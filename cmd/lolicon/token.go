package main

import (
	"os"

	"github.com/tenchoufansubs/go-lolicon/storage"
)

var Token string

func initToken() {
	var (
		err error
		ok  bool
	)

	Token, ok = os.LookupEnv("LOLICON_DISCORD_TOKEN")
	if !ok {
		Token, err = cache.Get("token")
		if err != nil {
			if err != storage.ErrNotFound {
				panic(err)
			}

			panic("token not found")
		}
	}
}
