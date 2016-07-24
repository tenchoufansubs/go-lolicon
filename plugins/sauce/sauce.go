package sauce

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/tenchoufansubs/go-lolicon"
	"github.com/tenchoufansubs/go-lolicon/storage"

	"github.com/bwmarrin/discordgo"

	"github.com/noisypixy/go-sauce"
)

func init() {
	plugin := new(SaucePlugin)

	var (
		_ lolicon.Plugin         = plugin
		_ lolicon.MessageHandler = plugin
	)

	lolicon.RegisterPlugin(plugin)
}

type SaucePlugin struct {
	cache storage.Driver
}

func (p *SaucePlugin) Id() lolicon.PluginId {
	return lolicon.PluginId("sauce")
}

func (p *SaucePlugin) Help() map[string]string {
	return map[string]string{
		"sauce": "Lookup the sauce of the last image posted",
	}
}

func (p *SaucePlugin) Setup(cache storage.Driver) (err error) {
	p.cache = cache
	return
}

func (p *SaucePlugin) Open(s *discordgo.Session) (err error) {
	return
}

func (p *SaucePlugin) Close() (err error) {
	return
}

func (p *SaucePlugin) HandleMessage(msg *lolicon.Message) (done bool, err error) {
	if msg.Command != "sauce" {
		return
	}

	done = true

	var (
		selfId string
	)

	selfId, err = p.cache.Get("userId")
	if err != nil {
		return
	}

	var (
		imgURLStr string = msg.Trailing

		imgURL *url.URL
	)

	if imgURLStr == "" {
		var (
			messages []*discordgo.Message
		)

		s := msg.Raw.Session

		messages, err = s.ChannelMessages(msg.Raw.Event.ChannelID, 50, msg.Raw.Event.ID, "")
		if err != nil {
			return
		}

		if len(messages) == 0 {
			err = errors.New("no messages")
			return
		}

		for _, imgMsg := range messages {
			if imgMsg.Author.ID == selfId {
				continue
			}
			if len(imgMsg.Attachments) > 0 {
				for _, at := range imgMsg.Attachments {
					imgURL, err = url.Parse(at.ProxyURL)
					if err != nil {
						continue
					}
					p := path.Base(imgURL.Path)
					if strings.HasSuffix(p, ".png") || strings.HasSuffix(p, ".jpg") || strings.HasSuffix(p, ".jpeg") {
						break
					}
					imgURL = nil
				}
			} else if len(imgMsg.Embeds) > 0 {
				for _, em := range imgMsg.Embeds {
					if em.Thumbnail == nil {
						continue
					}
					imgURL, err = url.Parse(em.Thumbnail.ProxyURL)
					if err != nil {
						continue
					}
					p := path.Base(imgURL.Path)
					if strings.HasSuffix(p, ".png") || strings.HasSuffix(p, ".jpg") || strings.HasSuffix(p, ".jpeg") {
						break
					}
					imgURL = nil
				}
			} else {
				continue
			}
			if imgURL != nil {
				break
			}
		}
	} else {
		if strings.HasPrefix(imgURLStr, "<") && strings.HasSuffix(imgURLStr, ">") {
			imgURLStr = imgURLStr[1 : len(imgURLStr)-1]
		}
		imgURL, err = url.Parse(imgURLStr)
		if err != nil {
			return
		}
	}

	var (
		img *sauce.Image
	)

	img, err = sauce.Sauce(imgURL.String())
	if err != nil {
		return
	}

	messageParts := make([]string, 0)

	if len(img.Characters) > 0 {
		messageParts = append(messageParts, "[Characters]")
		for _, c := range img.Characters {
			messageParts = append(messageParts, fmt.Sprintf("- %s", c.Name))
		}
		messageParts = append(messageParts, "")
	}

	if img.Series != nil {
		messageParts = append(messageParts, "[Series]")
		messageParts = append(messageParts, fmt.Sprintf("- %s", img.Series.Title))
		messageParts = append(messageParts, "")
	}

	if img.Artist != nil {
		messageParts = append(messageParts, "[Artist]")
		messageParts = append(messageParts, fmt.Sprintf("- %s", img.Artist.Name))
		messageParts = append(messageParts, "")
	}

	if len(img.Copyrights) > 0 {
		messageParts = append(messageParts, "[Copyrights]")
		for _, c := range img.Copyrights {
			messageParts = append(messageParts, fmt.Sprintf("- %s", c))
		}
		messageParts = append(messageParts, "")
	}

	message := fmt.Sprintf(
		"%s\n```\n%s\n```",
		img.SourceURL,
		strings.TrimSpace(
			strings.Join(messageParts, "\n"),
		),
	)

	_, err = msg.Raw.Session.ChannelMessageSend(msg.ChannelId, message)

	return
}
