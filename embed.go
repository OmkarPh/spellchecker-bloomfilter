package main

import (
	"embed"
	"os"
	"path/filepath"
)

//go:embed data/*txt
var content embed.FS

func EnsureEmbeddedFilesInStorageDir(storageDir string) {

	// Check if the target directory exists, if not, create it
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	// List files in the embedded filesystem
	files, err := content.ReadDir("data")
	if err != nil {
		panic(err)
	}

	// Copy files
	for _, file := range files {
		targetFilePath := filepath.Join(storageDir, file.Name())

		// If the file doesn't exist in the cache directory, copy it
		if _, err := os.Stat(targetFilePath); os.IsNotExist(err) {
			data, err := content.ReadFile(filepath.Join("data", file.Name()))
			if err != nil {
				panic(err)
			}
			if err := os.WriteFile(targetFilePath, data, 0644); err != nil {
				panic(err)
			}
		}
	}
}
