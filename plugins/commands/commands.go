package commands

import (
	"strings"

	"github.com/tenchoufansubs/go-lolicon"
	"github.com/tenchoufansubs/go-lolicon/storage"

	"github.com/bwmarrin/discordgo"
)

const DefaultPrefix = "!"

func init() {
	plugin := new(CommandsPlugin)

	var (
		_ lolicon.Plugin         = plugin
		_ lolicon.MessageHandler = plugin
	)

	lolicon.RegisterPlugin(plugin)
}

type CommandsPlugin struct {
	cache storage.Driver `json:"-"`

	UserId string `json:"-"`
	Prefix string `json:"prefix"`

	Admins []string `json:"admins"`
}

func (p *CommandsPlugin) Id() lolicon.PluginId {
	return lolicon.PluginId("commands")
}

func (p *CommandsPlugin) Setup(cache storage.Driver) (err error) {
	p.cache = cache

	return
}

func (p *CommandsPlugin) Open(s *discordgo.Session) (err error) {
	err = storage.GetJSON(p.cache, string(p.Id()), p)
	if err != nil {
		return
	}

	p.UserId, err = p.cache.Get("userId")
	if err != nil {
		if err != storage.ErrNotFound {
			return
		}
		err = nil
	}

	if p.Prefix == "" {
		p.Prefix = DefaultPrefix
	}

	if p.Admins == nil {
		p.Admins = make([]string, 0)
	}

	return
}

func (p *CommandsPlugin) Close() (err error) {
	err = storage.SetJSON(p.cache, string(p.Id()), p)
	return
}

func (p *CommandsPlugin) HandleMessage(msg *lolicon.Message) (done bool, err error) {
	if !strings.HasPrefix(msg.Content, p.Prefix) {
		return
	}

	content := strings.TrimPrefix(msg.Content, p.Prefix)
	parts := strings.SplitN(content, " ", 2)

	msg.Command = parts[0]
	if len(parts) >= 2 {
		msg.Trailing = parts[1]
	}

	msg.Command = strings.TrimSpace(msg.Command)
	msg.Trailing = strings.TrimSpace(msg.Trailing)

	msg.IsAdmin = false
	for _, adminId := range p.Admins {
		if msg.UserId == adminId {
			msg.IsAdmin = true
			break
		}
	}

	return
}
