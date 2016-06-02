package id

import (
	"fmt"

	"github.com/tenchoufansubs/go-lolicon"
	"github.com/tenchoufansubs/go-lolicon/storage"

	"github.com/bwmarrin/discordgo"
)

func init() {
	plugin := new(IdPlugin)

	var (
		_ lolicon.Plugin         = plugin
		_ lolicon.MessageHandler = plugin
	)

	lolicon.RegisterPlugin(plugin)
}

type IdPlugin struct{}

func (p *IdPlugin) Id() lolicon.PluginId {
	return lolicon.PluginId("id")
}

func (p *IdPlugin) Setup(cache storage.Driver) (err error) {
	return
}

func (p *IdPlugin) Open(s *discordgo.Session) (err error) {
	return
}

func (p *IdPlugin) Close() (err error) {
	return
}

func (p *IdPlugin) HandleMessage(msg *lolicon.Message) (done bool, err error) {
	if msg.Command != "id" {
		return
	}
	if msg.Trailing != "" {
		return
	}

	done = true

	message := fmt.Sprintf("<@%s> Your ID is `%s`. You're welcome.", msg.UserId, msg.UserId)

	_, err = msg.Raw.Session.ChannelMessageSend(msg.ChannelId, message)

	return
}
