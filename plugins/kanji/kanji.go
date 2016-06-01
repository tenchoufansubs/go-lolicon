package images

import (
	"fmt"
	"strings"

	"github.com/tenchoufansubs/go-lolicon"
	"github.com/tenchoufansubs/go-lolicon/storage"

	"github.com/bwmarrin/discordgo"
	"github.com/noisypixy/go-jisho"
)

func init() {
	plugin := new(KanjiPlugin)

	var (
		_ lolicon.Plugin         = plugin
		_ lolicon.MessageHandler = plugin
	)

	lolicon.RegisterPlugin(plugin)
}

type KanjiPlugin struct {
	cache storage.Driver `json:"-"`
}

func (p *KanjiPlugin) Id() lolicon.PluginId {
	return lolicon.PluginId("kanji")
}

func (p *KanjiPlugin) Setup(cache storage.Driver) (err error) {
	return
}

func (p *KanjiPlugin) Open(s *discordgo.Session) (err error) {
	return
}

func (p *KanjiPlugin) Close() (err error) {
	return
}

func (p *KanjiPlugin) HandleMessage(msg *lolicon.Message) (done bool, err error) {
	if msg.Command != "kanji" {
		return
	}
	if msg.Trailing == "" {
		return
	}

	done = true

	var (
		results []*jisho.Result
	)

	results, err = jisho.Search("#kanji " + msg.Trailing)
	if err != nil {
		return
	}

	messageParts := make([]string, 0)

	for _, result := range results {
		k := result.Kanji

		if k == nil || k.Character == "" || k.RawStrokes == "" || k.Meanings == "" {
			continue
		}

		messageParts = append(messageParts, fmt.Sprintf("%s (%s) -- %s", k.Character, k.RawStrokes, k.Meanings))
	}

	if len(messageParts) == 0 {
		_, err = msg.Raw.Session.ChannelMessageSend(msg.ChannelId, fmt.Sprintf("No results for %#v", msg.Trailing))
		return
	}

	message := "```\n" + strings.Join(messageParts, "\n") + "\n```"

	_, err = msg.Raw.Session.ChannelMessageSend(msg.ChannelId, message)

	return
}
