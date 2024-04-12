package bloomfilter

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"math/rand"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spaolacci/murmur3"
)

type DiskBloomFilterConfig struct {
	caseSensitive   bool
	hashFnCount     uint16
	hashFnSeeds     []uint32
	bloomFilterSize uint32
	headerOffset    uint64
}
type DiskBloomFilter struct {
	DiskBloomFilterConfig
	opbfFilePath       string
	dictionaryFilePath string
	hashFns            []hash.Hash32
	file               *os.File
	// inMemoryBf []byte
}

func CalculateHeaderOffset(config DiskBloomFilterConfig) uint64 {
	var IDENTIFIER_SIZE uint64 = 4
	var VERSION_SIZE uint64 = 2
	var CASE_SENSITIVITY uint64 = 2
	var HASH_FN_COUNT_SIZE uint64 = 2
	var BLOOM_FILTER_SIZE uint64 = 4

	return IDENTIFIER_SIZE + VERSION_SIZE + CASE_SENSITIVITY + HASH_FN_COUNT_SIZE + uint64(config.hashFnCount*4) + BLOOM_FILTER_SIZE
}

func BuildOpbfFile(opbfFilePath string, hashFnCount uint16, bfSize uint32, caseSensitivity bool) (DiskBloomFilterConfig, error) {
	config := DiskBloomFilterConfig{
		hashFnCount:     hashFnCount,
		hashFnSeeds:     make([]uint32, hashFnCount),
		bloomFilterSize: bfSize,
	}

	file, err := os.Create(opbfFilePath)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			color.Red(fmt.Sprintf("Permission denied to create file at %s\n", opbfFilePath))
		}
		return config, err
	}
	defer file.Close()

	opbfVersion := uint16(1)

	// OPBF identifer (uint32 - 4 bytes)
	file.Write([]byte("OPBF"))

	// Temporary reusable buffers during write
	buffer2Bytes := make([]byte, 2)
	buffer4Bytes := make([]byte, 4)

	// OPBF version (uint16 - 2 bytes)
	binary.BigEndian.PutUint16(buffer2Bytes, opbfVersion)
	file.Write(buffer2Bytes)

	// Case sensitivity (uint16 - 2 bytes)
	var caseSensitivityIdentifier uint16 = 0
	if caseSensitivity {
		caseSensitivityIdentifier = 1
	}
	binary.BigEndian.PutUint16(buffer2Bytes, caseSensitivityIdentifier)
	file.Write(buffer2Bytes)

	// No. of hash functions (uint16 - 2 bytes)
	binary.BigEndian.PutUint16(buffer2Bytes, hashFnCount)
	file.Write(buffer2Bytes)

	// Hash function seeds (uint32 - 4 bytes each)
	for i := range hashFnCount {
		config.hashFnSeeds[i] = rand.Uint32()
		binary.BigEndian.PutUint32(buffer4Bytes, config.hashFnSeeds[i])
		file.Write(buffer4Bytes)
	}

	// Bloomfilter size (in bits) (uint32 - 4 bytes)
	binary.BigEndian.PutUint32(buffer4Bytes, bfSize)
	file.Write(buffer4Bytes)

	// Bloomfilter bytes
	bytesRequired := (bfSize + 7) / 8 // Same as ceil(bfSize/8)
	bloomFilterBytes := make([]byte, bytesRequired)
	// fmt.Printf("For bloomfilter size %d, bytes needed: %d\n", bfSize, bytesRequired)
	file.Write(bloomFilterBytes)

	config.headerOffset = CalculateHeaderOffset(config)
	return config, nil
}

func ParseOpbfFile(opbfFilePath string) (bool, DiskBloomFilterConfig) {
	var config DiskBloomFilterConfig

	file, err := os.Open(opbfFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return false, config
	}

	// Temporary reusable buffers during read
	buffer2Bytes := make([]byte, 2)
	buffer4Bytes := make([]byte, 4)

	file.Read(buffer4Bytes)
	if !bytes.Equal(buffer4Bytes, []byte("OPBF")) {
		return false, config
	}

	file.Read(buffer2Bytes)
	version := binary.BigEndian.Uint16(buffer2Bytes)
	if uint16(1) != version {
		return false, config
	}

	file.Read(buffer2Bytes)
	caseSensitivityIdentifier := binary.BigEndian.Uint16(buffer2Bytes)
	if caseSensitivityIdentifier == 1 {
		config.caseSensitive = true
	} else {
		config.caseSensitive = false
	}

	file.Read(buffer2Bytes)
	hashFnCount := binary.BigEndian.Uint16(buffer2Bytes)

	hashFnSeeds := make([]uint32, hashFnCount)
	for i := range hashFnCount {
		file.Read(buffer4Bytes)
		hashFnSeeds[i] = binary.BigEndian.Uint32(buffer4Bytes)
	}

	file.Read(buffer4Bytes)
	bloomFilterSize := binary.BigEndian.Uint32(buffer4Bytes)

	config = DiskBloomFilterConfig{
		hashFnCount:     hashFnCount,
		hashFnSeeds:     hashFnSeeds,
		bloomFilterSize: bloomFilterSize,
	}
	config.headerOffset = CalculateHeaderOffset(config)
	return true, config
}

func (bf *DiskBloomFilter) LoadDictionary(dictFilePath string) error {
	dictionaryFile, err := os.Open(dictFilePath)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(dictionaryFile)
	defer dictionaryFile.Close()

	for scanner.Scan() {
		word := strings.Trim(scanner.Text(), " ,")
		bf.Add(word)
	}

	return nil
}

func InitDiskBloomFilter(opbfFilePath string, dictionaryFilePath string, hashFnCount uint16, size uint32, caseSensitive bool, forceReseed bool) *DiskBloomFilter {
	validOpbfFile, newConfig := ParseOpbfFile(opbfFilePath)
	var seedDictionary bool

	if validOpbfFile && !forceReseed {
		fmt.Printf("Using pre-seeded bloomfilter at %s\n", opbfFilePath)
		seedDictionary = false
	} else {
		var err error
		newConfig, err = BuildOpbfFile(opbfFilePath, hashFnCount, size, caseSensitive)
		if err != nil {
			color.Red(fmt.Sprintf("couldn't create bloom filter file at %s\n", opbfFilePath))
			os.Exit(127)
		}
		fmt.Printf("Initialised new bloomfilter at %s\n", opbfFilePath)
		seedDictionary = true
	}

	file, _ := os.OpenFile(opbfFilePath, os.O_RDWR, 0644)

	bf := DiskBloomFilter{
		opbfFilePath:          opbfFilePath,
		dictionaryFilePath:    dictionaryFilePath,
		DiskBloomFilterConfig: newConfig,
		hashFns:               make([]hash.Hash32, newConfig.hashFnCount),
		file:                  file,
	}

	for i := range bf.hashFnCount {
		bf.hashFns[i] = murmur3.New32WithSeed(bf.hashFnSeeds[i])
	}

	if seedDictionary {
		fmt.Printf("Feeding words from %s into bloom filter ...\n", dictionaryFilePath)
		err := bf.LoadDictionary(dictionaryFilePath)
		if err != nil {
			color.Red(fmt.Sprintf("couldn't find dictionary file at %s\n", dictionaryFilePath))
			os.Exit(127)
		}
	}

	fmt.Printf("Initialised disk bloom filter of size %d with %d hash functions\n", bf.bloomFilterSize, bf.hashFnCount)
	return &bf
}

func (bf *DiskBloomFilter) resolveIdx(key string, hashFnIdx uint32) (uint32, uint64, uint32) {
	bf.hashFns[hashFnIdx].Reset()
	bf.hashFns[hashFnIdx].Write([]byte(key))
	rawIdx := bf.hashFns[hashFnIdx].Sum32() % bf.bloomFilterSize
	byteIdx := uint64(rawIdx/8) + bf.headerOffset
	bitIdx := rawIdx % 8
	bf.hashFns[hashFnIdx].Reset()
	return rawIdx, byteIdx, bitIdx
}

func (bf *DiskBloomFilter) Close() {
	// fmt.Println("Closed file")
	bf.file.Close()
}

func (bf *DiskBloomFilter) Add(key string) {
	if !bf.caseSensitive {
		key = strings.ToLower(key)
	}

	for hashFnIdx := range len(bf.hashFns) {
		_, byteIdx, bitIdx := bf.resolveIdx(key, uint32(hashFnIdx))

		b := make([]byte, 1)
		bf.file.ReadAt(b, int64(byteIdx))

		// fmt.Println("Read => ", uint8(b[0]))

		b[0] = uint8(b[0]) | (1 << bitIdx)

		// fmt.Println("Wrote => ", uint8(b[0]))
		bf.file.WriteAt(b, int64(byteIdx))
	}
}

func (bf *DiskBloomFilter) Exists(key string) bool {
	if !bf.caseSensitive {
		key = strings.ToLower(key)
	}

	for hashFnIdx := range len(bf.hashFns) {
		_, byteIdx, bitIdx := bf.resolveIdx(key, uint32(hashFnIdx))
		b := make([]byte, 1)
		bf.file.ReadAt(b, int64(byteIdx))

		exists := (uint8(b[0]) & (1 << bitIdx)) > 0
		if !exists {
			return false
		}
	}
	return true
}
