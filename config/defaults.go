package config

import (
	"os"
	"path/filepath"
)

type DefaultParams struct {
	DictPath    string
	OpbfPath    string
	HashFnCount uint16
	BfSize      uint32
}

func GetDefaultParams(storageDir string) DefaultParams {
	return DefaultParams{
		DictPath:    filepath.Join(storageDir, "dict.txt"),
		OpbfPath:    filepath.Join(storageDir, "dict.opbf"),
		HashFnCount: 4,
		BfSize:      4000000,
	}
}

func GetDefaultStorageDir() string {
	cacheDir, _ := os.UserCacheDir()
	cacheDir = filepath.Join(cacheDir, "spellchecker-bloomfilter")
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.MkdirAll(cacheDir, os.ModePerm)
	}
	return cacheDir
}
