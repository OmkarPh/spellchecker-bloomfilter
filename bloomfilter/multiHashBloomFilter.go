package bloomfilter

import (
	"hash"
	"math/rand"

	"github.com/spaolacci/murmur3"
)

type MultiHashBloomFilter struct {
	size    uint32
	hashFns []hash.Hash32
	filter  []byte
}

func InitMultiHashBloomFilter(size uint32, hashFnCount uint32) *MultiHashBloomFilter {
	bf := MultiHashBloomFilter{
		size:    size,
		hashFns: make([]hash.Hash32, hashFnCount),
		filter:  make([]byte, (size+7)/8), // Same as ceil(size/8)
	}
	for i := range hashFnCount {
		bf.hashFns[i] = murmur3.New32WithSeed(rand.Uint32())
	}
	return &bf
}

func (bf *MultiHashBloomFilter) resolveIdx(key string, hashFnIdx uint32) (uint32, uint32, uint32) {
	bf.hashFns[hashFnIdx].Reset()
	bf.hashFns[hashFnIdx].Write([]byte(key))
	rawIdx := bf.hashFns[hashFnIdx].Sum32() % bf.size
	byteIdx := rawIdx / 8
	bitIdx := rawIdx % 8
	bf.hashFns[hashFnIdx].Reset()
	return rawIdx, byteIdx, bitIdx
}

func (bf *MultiHashBloomFilter) Add(key string) {
	for hashFnIdx := range len(bf.hashFns) {
		_, byteIdx, bitIdx := bf.resolveIdx(key, uint32(hashFnIdx))
		bf.filter[byteIdx] = bf.filter[byteIdx] | (1 << bitIdx)
	}
}

func (bf *MultiHashBloomFilter) Exists(key string) bool {
	for hashFnIdx := range len(bf.hashFns) {
		_, byteIdx, bitIdx := bf.resolveIdx(key, uint32(hashFnIdx))
		exists := (bf.filter[byteIdx] & (1 << bitIdx)) > 0
		if !exists {
			return false
		}
	}
	return true
}
