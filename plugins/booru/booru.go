package booru

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/tenchoufansubs/go-lolicon"
	"github.com/tenchoufansubs/go-lolicon/storage"

	"github.com/bwmarrin/discordgo"
	"github.com/noisypixy/go-booru"
)

func init() {
	plugin := new(BooruPlugin)

	var (
		_ lolicon.Plugin         = plugin
		_ lolicon.MessageHandler = plugin
	)

	lolicon.RegisterPlugin(plugin)
}

type BooruPlugin struct{}

func (p *BooruPlugin) Id() lolicon.PluginId {
	return lolicon.PluginId("booru")
}

func (p *BooruPlugin) Setup(cache storage.Driver) (err error) {
	return
}

func (p *BooruPlugin) Open(s *discordgo.Session) (err error) {
	return
}

func (p *BooruPlugin) Close() (err error) {
	return
}

func (p *BooruPlugin) HandleMessage(msg *lolicon.Message) (done bool, err error) {
	if msg.Command != "" {
		return
	}
	if msg.Trailing == "" {
		return
	}
	if strings.HasSuffix(msg.Trailing, ".png") || strings.HasSuffix(msg.Trailing, ".jpg") || strings.HasSuffix(msg.Trailing, ".jpeg") || strings.HasSuffix(msg.Trailing, ".gif") {
		return
	}

	var (
		parsedURL *url.URL
		image     *booru.Image
		resp      *http.Response
	)

	parts := strings.Split(msg.Trailing, " ")
	for _, p := range parts {
		var (
			err error
		)

		parsedURL = nil

		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		parsedURL, err = url.Parse(p)
		if err != nil {
			continue
		}

		if booru.Supports(parsedURL) {
			break
		}
	}

	if parsedURL == nil {
		return
	}

	done = true

	image, err = booru.Resolve(parsedURL.String())
	if err != nil {
		return
	}

	if parsedURL.String() == image.URL {
		return
	}

	resp, err = http.Get(image.URL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, err = msg.Raw.Session.ChannelFileSend(msg.ChannelId, image.Filename, resp.Body)

	return
}