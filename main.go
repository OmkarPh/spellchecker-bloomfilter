package main

import (
	"fmt"

	"github.com/OmkarPh/spellchecker-bloomfilter/bloomfilter"
)

func main() {
	// bloomfilter.TestInMemoryBloomfilter()

	bf := bloomfilter.InitDiskBloomFilter("dict.opbf", "dict.txt", 4, 5000000, false)
	defer bf.Close()

	testcases := []string{"good", "bad", "write", "rgiht", "lead", "newd", "feed", "leda", "mikl"}
	incorrectDetected := []string{}
	for _, word := range testcases {
		if !bf.Exists(word) {
			incorrectDetected = append(incorrectDetected, word)
		}
	}

	fmt.Println()
	fmt.Print("Incorrect words: ")
	fmt.Printf("%v\n", incorrectDetected)

	// added := []string{"a", "b", "c", "d", "e", "f"}
	// verify := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}

	// validKeys := make(map[string]bool)
	// for _, key := range added {
	// 	validKeys[key] = true
	// }

	// for _, key := range added {
	// 	bf.Add(key)
	// }
	// for _, key := range verify {
	// 	actuallyExists := validKeys[key]
	// 	exists := bf.Exists(key)
	// 	fmt.Printf("Check '%s' \t=> %t (Req: %t)\n", key, exists, actuallyExists)
	// }

}
