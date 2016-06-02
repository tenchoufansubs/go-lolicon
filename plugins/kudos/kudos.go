package kudos

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tenchoufansubs/go-lolicon"
	"github.com/tenchoufansubs/go-lolicon/storage"

	"github.com/bwmarrin/discordgo"
	"github.com/noisypixy/go-jisho"
)

func init() {
	plugin := new(KudosPlugin)

	var (
		_ lolicon.Plugin         = plugin
		_ lolicon.MessageHandler = plugin
	)

	lolicon.RegisterPlugin(plugin)
}

type KudosPlugin struct {
	cache  storage.Driver `json:"-"`
	selfId string         `json:"-"`

	Kudos map[string]int `json:"kudos"`
}

func (p *KudosPlugin) Id() lolicon.PluginId {
	return lolicon.PluginId("kudos")
}

func (p *KudosPlugin) Setup(cache storage.Driver) (err error) {
	p.cache = cache

	p.selfId, err = cache.Get("userId")
	if err != nil {
		return
	}

	if p.selfId == "" {
		err = errors.New("couldn't get own id")
		return
	}

	return
}

func (p *KudosPlugin) Open(s *discordgo.Session) (err error) {
	err = storage.GetJSON(p.cache, string(p.Id()), p)
	if err != nil {
		return
	}

	if p.Kudos == nil {
		p = make(map[string]int)
	}

	return
}

func (p *KudosPlugin) Close() (err error) {
	err = storage.SetJSON(p.cache, string(p.Id()), p)
	return
}

func (p *KudosPlugin) HandleMessage(msg *lolicon.Message) (done bool, err error) {
	if msg.Command != "kudos" && msg.Command != "kudos?" && msg.Command != "damedesu" {
		return
	}
	if msg.Trailing == "" {
		return
	}
	if !msg.IsAdmin {
		return
	}
	if len(msg.Raw.Event.Mentions) == 0 {
		return
	}

	mentions := make([]*discordgo.User, 0)
	for _, user := range msg.Raw.Event.Mentions {
		if user.ID == p.selfId {
			continue
		}

		mentions = append(mentions, user)
	}
	if len(mentions) != 1 {
		return
	}

	done = true

	user := mentions[0]
	message := ""

	switch msg.Command {
	case "kudos":
		p.Kudos[user.ID] += 1

		kudosStr := "kudo"
		if p.Kudos[user.ID] != 1 {
			kudosStr = "kudos"
		}

		message = fmt.Sprintf("<@%s> now has %d %s!", user.ID, p.Kudos[user.ID], kudosStr)

		if p.Kudos[user.ID] == 0 {
			message += " "
			message += "(Welcome back. Just don't do it again, 'kay?)"
		} else if p.Kudos[user.ID] < 0 {
			message += " "
			message += "(Come on, you can do it)"
		}

		break

	case "damedesu":
		p.Kudos[user.ID] -= 1

		kudosStr := "kudo"
		if p.Kudos[user.ID] != 1 {
			kudosStr = "kudos"
		}

		message = fmt.Sprintf("Dame desu, <@%s>. You now have %d %s.", user.ID, p.Kudos[user.ID], kudosStr)
		if p.Kudos[user.ID] == 0 {
			message += " "
			message += "(All that effort went to waste)"
		} else if p.Kudos[user.ID] < 0 {
			message += " "
			message += "(You don't learn, do you?)"
		}

		break

	case "kudos?":
		kudosStr := "kudo"
		if p.Kudos[user.ID] != 1 {
			kudosStr = "kudos"
		}

		message = fmt.Sprintf("<@%s> has **%d** %s.", user.ID, p.Kudos[user.ID], kudosStr)
		if p.Kudos[user.ID] < 0 {
			message += " "
			message += "(The hell is wrong with you?)"
		}

		break
	}

	_, err = msg.Raw.Session.ChannelMessageSend(msg.ChannelId, message)

	return
}
