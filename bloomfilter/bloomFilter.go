package bloomfilter

import (
	"hash"

	"github.com/spaolacci/murmur3"
)

type BloomFilter struct {
	size   uint32
	hasher hash.Hash32
	filter []byte
}

func InitBloomFilter(size uint32, seed uint32) *BloomFilter {
	return &BloomFilter{
		size:   size,
		hasher: murmur3.New32WithSeed(seed),
		filter: make([]byte, (size+7)/8), // Same as ceil(size/8)
	}
}

func (bf *BloomFilter) resolveIdx(key string) (uint32, uint32, uint32) {
	bf.hasher.Reset()
	bf.hasher.Write([]byte(key))
	rawIdx := bf.hasher.Sum32() % bf.size
	byteIdx := rawIdx / 8
	bitIdx := rawIdx % 8
	bf.hasher.Reset()
	return rawIdx, byteIdx, bitIdx
}

func (bf *BloomFilter) Add(key string) uint32 {
	rawIdx, byteIdx, bitIdx := bf.resolveIdx(key)
	bf.filter[byteIdx] = bf.filter[byteIdx] | (1 << bitIdx)
	return rawIdx
}

func (bf *BloomFilter) ExistsUtil(key string) (uint32, bool) {
	rawIdx, byteIdx, bitIdx := bf.resolveIdx(key)
	exists := (bf.filter[byteIdx] & (1 << bitIdx)) > 0
	return rawIdx, exists
}

func (bf *BloomFilter) Exists(key string) bool {
	_, exists := bf.ExistsUtil(key)
	return exists
}
