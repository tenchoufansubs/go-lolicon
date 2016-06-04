package nyaa

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tenchoufansubs/go-lolicon"
	"github.com/tenchoufansubs/go-lolicon/storage"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
)

func init() {
	plugin := new(NyaaPlugin)

	var (
		_ lolicon.Plugin         = plugin
		_ lolicon.MessageHandler = plugin
	)

	lolicon.RegisterPlugin(plugin)
}

type NyaaPlugin struct {
	cache   storage.Driver     `json:"-"`
	session *discordgo.Session `json:"-"`

	NyaaId   string              `json:"nyaa_id"`
	Torrents map[string]*Torrent `json:"torrents"`

	CheckIntervalMinutes int `json:"check_interval_minutes"`

	NotificationsChannelId    string `json:"notifications_channel_id"`
	NotificationDownloadsStep int    `json:"notification_downloads_step"`

	chInterrupt chan struct{} `json:"-"`
}

func (p *NyaaPlugin) Id() lolicon.PluginId {
	return lolicon.PluginId("nyaa")
}

func (p *NyaaPlugin) Setup(cache storage.Driver) (err error) {
	p.cache = cache

	return
}

func (p *NyaaPlugin) Open(s *discordgo.Session) (err error) {
	err = storage.GetJSON(p.cache, string(p.Id()), p)
	if err != nil {
		return
	}

	p.session = s

	if p.Torrents == nil {
		p.Torrents = make(map[string]*Torrent)
	}

	if p.CheckIntervalMinutes < 5 {
		p.CheckIntervalMinutes = 5
	}

	if p.NotificationDownloadsStep < 0 {
		p.NotificationDownloadsStep = 0
	} else if p.NotificationDownloadsStep < 500 {
		p.NotificationDownloadsStep = 500
	}

	p.chInterrupt = make(chan struct{}, 1)

	go p.watch()

	return
}

func (p *NyaaPlugin) Close() (err error) {
	close(p.chInterrupt)

	err = storage.SetJSON(p.cache, string(p.Id()), p)
	return
}

func (p *NyaaPlugin) HandleMessage(msg *lolicon.Message) (done bool, err error) {
	if msg.Command != "nyaa" {
		return
	}
	if msg.Trailing == "" {
		return
	}

	done = true

	// TODO

	return
}

func (p *NyaaPlugin) watch() {
	for {
		_ = p.update()

		select {
		case <-p.chInterrupt:
			return
		case <-time.After(time.Duration(p.CheckIntervalMinutes) * time.Minute):
			continue
		}
	}
}

func (p *NyaaPlugin) update() (err error) {
	if p.NyaaId == "" {
		return
	}

	var (
		torrentsURL *url.URL

		doc *goquery.Document
	)

	torrentsURLStr := fmt.Sprintf("http://www.nyaa.se/?user=")

	torrentsURL, err = url.Parse(torrentsURLStr)
	if err != nil {
		return
	}

	q := torrentsURL.Query()
	q.Set("user", p.NyaaId)
	torrentsURL.RawQuery = q.Encode()

	torrentsURLStr = torrentsURL.String()

	doc, err = goquery.NewDocument(torrentsURLStr)
	if err != nil {
		return
	}

	torrents := make(map[string]*Torrent)

	doc.Find(".tlistrow").Each(func(i int, row *goquery.Selection) {
		title := row.Find(".tlistname").First().Text()
		title = strings.TrimSpace(title)

		if title == "" {
			return
		}

		downloadsStr := row.Find(".tlistdn").First().Text()
		downloadsStr = strings.TrimSpace(downloadsStr)

		if downloadsStr == "" {
			return
		}

		downloads, err := strconv.Atoi(downloadsStr)
		if err != nil {
			return
		}

		a := row.Find(".tlistname a").First()
		if a.Length() == 0 {
			return
		}

		href := a.AttrOr("href", "")
		if href == "" {
			return
		}

		torrentURL, err := torrentsURL.Parse(href)
		if err != nil {
			return
		}

		id := torrentURL.Query().Get("tid")
		if id == "" {
			return
		}

		torrents[id] = &Torrent{
			Id:        id,
			Title:     title,
			URL:       torrentURL.String(),
			Downloads: downloads,
		}
	})

	step := float64(p.NotificationDownloadsStep)

	for id, torrent := range torrents {
		// If the torrent is "new", we won't do anything with it
		// until next check.
		if _, ok := p.Torrents[id]; !ok {
			p.Torrents[id] = torrent
			continue
		}

		if p.NotificationsChannelId != "" && step > 0 {
			dlPrev := float64(p.Torrents[id].Downloads)
			dlCur := float64(torrent.Downloads)

			a := math.Floor(dlPrev / step)
			b := math.Floor(dlCur / step)

			if b > a {
				milestone := int(math.Floor(b * step))
				message := fmt.Sprintf("%s has reached %d downloads! %s", torrent.Title, milestone, torrent.URL)
				_, _ = p.session.ChannelMessageSend(p.NotificationsChannelId, message)
			}
		}

		// Save changes.
		p.Torrents[id] = torrent
	}

	return
}
