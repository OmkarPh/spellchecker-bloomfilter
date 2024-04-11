package bloomfilter

import (
	"fmt"
)

func TestInMemoryBloomfilter() {
	added := []string{"a", "b", "c", "d", "e", "f"}
	verify := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}

	validKeys := make(map[string]bool)
	for _, key := range added {
		validKeys[key] = true
	}

	bfSize := uint32(32)
	bfSeed := 226
	fmt.Printf("Single hash Bloom filter with size %d\n", bfSize)
	bf := InitBloomFilter(bfSize, uint32(bfSeed))
	for _, key := range added {
		bf.Add(key)
	}
	for _, key := range verify {
		actuallyExists := validKeys[key]
		idx, exists := bf.ExistsUtil(key)
		fmt.Printf("Check '%s' (%d)\t=> %t (Req: %t)\n", key, idx, exists, actuallyExists)
	}

	fmt.Println()

	mhBfSize := uint32(32)
	hashFnCount := uint32(3)
	fmt.Printf("Multi hash (%d hashes) with size %d\n", hashFnCount, mhBfSize)
	mhBf := InitMultiHashBloomFilter(mhBfSize, hashFnCount)
	for _, key := range added {
		mhBf.Add(key)
	}
	for _, key := range verify {
		actuallyExists := validKeys[key]
		exists := mhBf.Exists(key)
		fmt.Printf("Check '%s' => %t (Req: %t)\n", key, exists, actuallyExists)
	}

}
