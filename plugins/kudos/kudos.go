package kudos

import (
	"errors"
	"fmt"
	"os"

	"github.com/tenchoufansubs/go-lolicon"
	"github.com/tenchoufansubs/go-lolicon/storage"

	"github.com/bwmarrin/discordgo"
)

func init() {
	plugin := new(KudosPlugin)

	var (
		_ lolicon.Plugin         = plugin
		_ lolicon.MessageHandler = plugin
		_ lolicon.HelpProvider   = plugin
	)

	lolicon.RegisterPlugin(plugin)
}

type KudosPlugin struct {
	cache  storage.Driver `json:"-"`
	selfId string         `json:"-"`

	LogFile string         `json:"log_file"`
	Kudos   map[string]int `json:"kudos"`

	logf *os.File `json:"-"`
}

func (p *KudosPlugin) Id() lolicon.PluginId {
	return lolicon.PluginId("kudos")
}

func (p *KudosPlugin) Help() map[string]string {
	return map[string]string{
		"kudos <user>":    "Give kudos to <user>",
		"damedesu <user>": "Remove 1 kudo from <user>",
		"kudos? <user>":   "Display amount of kudos received by <user>",
	}
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

	if p.LogFile == "" {
		p.LogFile = "kudos.log"
	}

	return
}

func (p *KudosPlugin) Open(s *discordgo.Session) (err error) {
	err = storage.GetJSON(p.cache, string(p.Id()), p)
	if err != nil {
		return
	}

	if p.Kudos == nil {
		p.Kudos = make(map[string]int)
	}

	p.logf, err = os.OpenFile(p.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return
	}

	return
}

func (p *KudosPlugin) Close() (err error) {
	_ = p.logf.Close()
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

	var (
		channel   *discordgo.Channel
		guild     *discordgo.Guild
		whoSaidIt *discordgo.User
	)

	user := mentions[0]

	channel, err = msg.Raw.Session.Channel(msg.ChannelId)
	if err != nil {
		return
	}

	guild, err = msg.Raw.Session.Guild(channel.GuildID)
	if err != nil {
		return
	}

	whoSaidIt, err = msg.Raw.Session.User(msg.UserId)
	if err != nil {
		return
	}

	message := ""

	switch msg.Command {
	case "kudos":
		p.Kudos[user.ID] += 1

		kudosStr := "kudo"
		if p.Kudos[user.ID] != 1 {
			kudosStr = "kudos"
		}

		message = fmt.Sprintf("<@%s> now has **%d** %s!", user.ID, p.Kudos[user.ID], kudosStr)

		if p.Kudos[user.ID] == 0 {
			message += " "
			message += "(Welcome back. Just don't do it again, 'kay?)"
		} else if p.Kudos[user.ID] < 0 {
			message += " "
			message += "(Come on, you can do it)"
		}

		_, err = fmt.Fprintf(p.logf, "[%s] [%s|%s] [%s|#%s] -- %s <%s> sent +1 kudo to %s <%s>\n", msg.Date, guild.ID, guild.Name, channel.ID, channel.Name, whoSaidIt.Username, whoSaidIt.ID, user.Username, user.ID)

		break

	case "damedesu":
		p.Kudos[user.ID] -= 1

		kudosStr := "kudo"
		if p.Kudos[user.ID] != 1 {
			kudosStr = "kudos"
		}

		message = fmt.Sprintf("Dame desu, <@%s>. You now have **%d** %s.", user.ID, p.Kudos[user.ID], kudosStr)
		if p.Kudos[user.ID] == 0 {
			message += " "
			message += "(All that effort went to waste)"
		} else if p.Kudos[user.ID] < 0 {
			message += " "
			message += "(You don't learn, do you?)"
		}

		_, err = fmt.Fprintf(p.logf, "[%s] [%s|%s] [%s|#%s] -- %s <%s> took 1 kudo from %s <%s>\n", msg.Date, guild.ID, guild.Name, channel.ID, channel.Name, whoSaidIt.Username, whoSaidIt.ID, user.Username, user.ID)

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
