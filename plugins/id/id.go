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

type IdPlugin struct {
	selfId string `json:"-"`
}

func (p *IdPlugin) Id() lolicon.PluginId {
	return lolicon.PluginId("id")
}

func (p *IdPlugin) Setup(cache storage.Driver) (err error) {
	p.selfId, err = cache.Get("userId")
	if err != nil {
		return
	}

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

	var (
		message string
	)

	mentions := make([]*discordgo.User, 0)
	for _, user := range msg.Raw.Event.Mentions {
		if user.ID == p.selfId {
			continue
		}

		mentions = append(mentions, user)
	}

	if len(mentions) > 1 {
		return
	}

	if len(mentions) == 0 {
		message = fmt.Sprintf("<@%s> Your ID is `%s`. You're welcome.", msg.UserId, msg.UserId)
	} else {
		message = fmt.Sprintf("The ID of <@%s> is `%s`. You're welcome.", mentions[0], msg.UserId)
	}

	done = true

	_, err = msg.Raw.Session.ChannelMessageSend(msg.ChannelId, message)

	return
}
