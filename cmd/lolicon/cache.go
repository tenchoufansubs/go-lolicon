package main

import (
	"os"
	"path"

	"github.com/tenchoufansubs/go-lolicon/storage"
	"github.com/tenchoufansubs/go-lolicon/storage/jsonstorage"
)

var (
	cache storage.Driver
)

func initCache() {
	var (
		err  error
		file string
	)

	cache, err = jsonstorage.New()
	if err != nil {
		panic(err)
	}

	file, err = CacheFile()
	if err != nil {
		panic(err)
	}

	err = cache.Open(file)
	if err != nil {
		panic(err)
	}
}

// CacheFile returns the path to the cache file.
//
// It may exist or not.
func CacheFile() (file string, err error) {
	var (
		wd string
	)

	wd, err = os.Getwd()
	if err != nil {
		return
	}

	file = path.Join(wd, "lolicon.json")
	return
}
