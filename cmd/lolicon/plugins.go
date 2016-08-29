package main

import (
	"log"

	"github.com/tenchoufansubs/go-lolicon"

	_ "github.com/tenchoufansubs/go-lolicon/plugins/commands"

	_ "github.com/tenchoufansubs/go-lolicon/plugins/booru"
	_ "github.com/tenchoufansubs/go-lolicon/plugins/id"
	_ "github.com/tenchoufansubs/go-lolicon/plugins/kanji"
	_ "github.com/tenchoufansubs/go-lolicon/plugins/kudos"
	_ "github.com/tenchoufansubs/go-lolicon/plugins/nyaa"
	_ "github.com/tenchoufansubs/go-lolicon/plugins/pixiv"
	_ "github.com/tenchoufansubs/go-lolicon/plugins/sauce"

	_ "github.com/tenchoufansubs/go-lolicon/plugins/images"

	_ "github.com/tenchoufansubs/go-lolicon/plugins/upload"
)

var (
	plugins []lolicon.Plugin
)

func initPlugins() {
	plugins = lolicon.Plugins()

	for _, p := range plugins {
		log.Printf("setup: %s", p.Id())
		p.Setup(cache)
	}
}
