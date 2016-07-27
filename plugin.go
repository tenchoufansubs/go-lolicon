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

// Plugin is a plugin.
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

// MessageHandler is a plugin that can handle messages.
type MessageHandler interface {
	Plugin
	HandleMessage(msg *Message) (done bool, err error)
}

// HelpProvider is a plugin that provides information for the "help" command.
type HelpProvider interface {
	Plugin

	// Help must return a map describing every supported command.
	//
	// A map key is the trigger for a command (without any prefix),
	// and its value is the description for that trigger.
	Help() map[string]string
}

// RegisterPlugin registers a plugin.
//
// This function panics if a plugin with the same id was already registered.
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

// Plugins returns a slice containing the registered plugins.
//
// The slice is sorted in the order the plugins were registered. This means
// the first element of the slice is the plugin that was registered first.
func Plugins() []Plugin {
	pluginSlice := make([]Plugin, len(pluginOrder))

	for i, id := range pluginOrder {
		pluginSlice[i] = plugins[id]
	}

	return pluginSlice
}
