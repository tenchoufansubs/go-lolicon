package pixiv

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/tenchoufansubs/go-lolicon"
	"github.com/tenchoufansubs/go-lolicon/storage"

	"github.com/bwmarrin/discordgo"

	"github.com/PuerkitoBio/goquery"
)

func init() {
	plugin := new(PixivPlugin)

	var (
		_ lolicon.Plugin         = plugin
		_ lolicon.MessageHandler = plugin
	)

	lolicon.RegisterPlugin(plugin)
}

type PixivPlugin struct{}

func (p *PixivPlugin) Id() lolicon.PluginId {
	return lolicon.PluginId("pixiv")
}

func (p *PixivPlugin) Setup(cache storage.Driver) (err error) {
	return
}

func (p *PixivPlugin) Open(s *discordgo.Session) (err error) {
	return
}

func (p *PixivPlugin) Close() (err error) {
	return
}

func (p *PixivPlugin) HandleMessage(msg *lolicon.Message) (done bool, err error) {
	if msg.Command != "" {
		return
	}

	var (
		parsedURL *url.URL
		req       *http.Request
		resp      *http.Response
		doc       *goquery.Document
	)

	parts := strings.Split(msg.Content, " ")
	for _, p := range parts {
		var (
			err error
		)

		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		if !strings.Contains(p, "pixiv.net") {
			continue
		}

		parsedURL, err = url.Parse(p)
		if err != nil {
			continue
		}

		break
	}

	if parsedURL == nil {
		return
	}

	done = true

	doc, err = goquery.NewDocument(parsedURL.String())
	if err != nil {
		return
	}

	imageURLStr := doc.Find(`meta[property="og:image"]`).First().AttrOr("content", "")
	imageURLStr = strings.TrimSpace(imageURLStr)
	if imageURLStr == "" {
		return
	}

	req, err = http.NewRequest(http.MethodGet, imageURLStr, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "lolicon")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	filename := parsedURL.Query().Get("illust_id")
	if filename == "" {
		return
	}

	switch resp.Header.Get("Content-Type") {
	case "image/png":
		filename += ".png"
		break
	case "image/jpeg":
		fallthrough
	case "image/jpg":
		filename += ".jpg"
		break
	default:
		return
	}

	_, err = msg.Raw.Session.ChannelFileSend(msg.ChannelId, filename, resp.Body)

	return
}
