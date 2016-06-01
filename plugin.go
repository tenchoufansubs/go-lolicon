package lolicon

import (
	"fmt"

	"github.com/tenchoufansubs/go-lolicon/storage"

	"github.com/bwmarrin/discordgo"
)

type PluginId string

var (
	plugins     map[PluginId]Plugin
	pluginOrder []PluginId
)

func initPlugin() {
	plugins = make(map[PluginId]Plugin)
	pluginOrder = make([]PluginId, 0)
}

type Plugin interface {
	Id() PluginId

	// Setup is called when the plugin is loaded.
	Setup(storage storage.Driver) (err error)

	// Open is called when the connection to Discord has been
	// successfully established.
	Open(s *discordgo.Session) (err error)

	// Close is called after the connection to Discord has been closed.
	Close() (err error)
}

type MessageHandler interface {
	Plugin
	HandleMessage(msg *Message) (done bool, err error)
}

func RegisterPlugin(plugin Plugin) {
	if plugin == nil {
		panic("nil plugin")
	}

	id := plugin.Id()

	if _, ok := plugins[id]; ok {
		panic(fmt.Sprintf("plugin %#v already registered", id))
	}

	plugins[id] = plugin
	pluginOrder = append(pluginOrder, id)
}

func Plugins() []Plugin {
	pluginSlice := make([]Plugin, len(pluginOrder))

	for i, id := range pluginOrder {
		pluginSlice[i] = plugins[id]
	}

	return pluginSlice
}
