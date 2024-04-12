package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/OmkarPh/spellchecker-bloomfilter/bloomfilter"
	"github.com/OmkarPh/spellchecker-bloomfilter/config"
)

func main() {
	// bloomfilter.TestInMemoryBloomfilter()

	storageDir := config.GetDefaultStorageDir()
	EnsureEmbeddedFilesInStorageDir(storageDir)

	defaultParams := config.GetDefaultParams(storageDir)

	// Flags
	dictPath := (flag.String("dict", defaultParams.DictPath, "Dictionary file path"))
	build := flag.Bool("build", false, "Rebuild the bloom filter")
	hashFnCount := flag.Uint("hashes", uint(defaultParams.HashFnCount), "Number of hash functions")
	bfSize := flag.Uint("size", uint(defaultParams.BfSize), "Bloom filter size")
	flag.Parse()

	forceReseed := (*dictPath != defaultParams.DictPath) || *build
	bf := bloomfilter.InitDiskBloomFilter(defaultParams.OpbfPath, *dictPath, uint16(*hashFnCount), uint32(*bfSize), false, forceReseed)
	defer bf.Close()

	// Words to check for
	words := flag.Args()

	if len(words) == 0 {
		if !forceReseed {
			fmt.Println("\nPlease provide words to check !")
		}
		return
	}

	// words := []string{"good", "bad", "write", "rgiht", "lead", "newd", "feed", "leda", "mikl"}

	incorrectDetected := []string{}
	for _, word := range words {
		if !bf.Exists(word) {
			incorrectDetected = append(incorrectDetected, word)
		}
	}

	fmt.Println()
	if len(incorrectDetected) > 0 {
		fmt.Println("Incorrect words:", strings.Join(incorrectDetected, ", "))
	} else {
		fmt.Println("All words are spelt correctly !")
	}
}
