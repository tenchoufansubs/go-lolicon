package upload

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/tenchoufansubs/go-lolicon"
	"github.com/tenchoufansubs/go-lolicon/storage"

	"github.com/bwmarrin/discordgo"
)

func init() {
	plugin := new(UploadPlugin)

	var (
		_ lolicon.Plugin         = plugin
		_ lolicon.MessageHandler = plugin
	)

	lolicon.RegisterPlugin(plugin)
}

type UploadPlugin struct {
	cache storage.Driver `json:"-"`

	Directories []string `json:"directories"`
}

func (p *UploadPlugin) Id() lolicon.PluginId {
	return lolicon.PluginId("upload")
}

func (p *UploadPlugin) Help() map[string]string {
	return map[string]string{
		"upload <directory> [URL...]": "Upload the given images to a directory.",
	}
}

func (p *UploadPlugin) Setup(cache storage.Driver) (err error) {
	p.cache = cache
	return
}

func (p *UploadPlugin) Open(s *discordgo.Session) (err error) {
	err = storage.GetJSON(p.cache, "images", p)
	if err != nil {
		return
	}

	if p.Directories == nil {
		p.Directories = make([]string, 0)
	}

	return
}

func (p *UploadPlugin) Close() (err error) {
	return
}

func (p *UploadPlugin) HandleMessage(msg *lolicon.Message) (done bool, err error) {
	if msg.Command != "upload" {
		return
	}

	done = true

	if len(p.Directories) == 0 {
		err = errors.New("no configured directories")
		return
	}

	parts := strings.Split(msg.Trailing, " ")
	if len(parts) == 0 {
		err = errors.New("dirname is required")
		return
	}

	dirname := strings.TrimSpace(parts[0])
	if dirname == "" {
		err = errors.New("empty dirname")
		return
	}
	if strings.Contains(dirname, "/") {
		err = errors.New("invalid dirname")
		return
	}

	uploadDirectory := path.Join(p.Directories[0], dirname)

	urlList := make([]string, 0)
	for i := 1; i < len(parts); i++ {
		urlStr := strings.TrimSpace(parts[i])
		if urlStr == "" {
			continue
		}

		urlList = append(urlList, urlStr)
	}

	attachments := msg.Raw.Event.Attachments
	for i := 0; i < len(attachments); i++ {
		attachment := attachments[i]
		urlStr := attachment.URL

		ext := path.Ext(attachment.Filename)
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".gif" {
			continue
		}

		urlList = append(urlList, urlStr)
	}

	if len(urlList) == 0 {
		err = errors.New("no urls")
		return
	}

	_, err = os.Open(uploadDirectory)
	if err != nil {
		err = os.Mkdir(uploadDirectory, 0775)
	}
	if err != nil {
		return
	}

	for _, urlStr := range urlList {
		var (
			imageURL *url.URL
		)

		imageURL, err = url.Parse(urlStr)
		if err != nil {
			return
		}

		filename := fmt.Sprintf("%d-%s", time.Now().Unix(), path.Base(imageURL.Path))

		err = func() (err error) {
			var (
				resp *http.Response
				f    *os.File
			)

			resp, err = http.Get(urlStr)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			file := path.Join(uploadDirectory, filename)

			f, err = os.Create(file)
			if err != nil {
				return
			}

			_, err = io.Copy(f, resp.Body)
			return
		}()

		if err != nil {
			return
		}
	}

	message := ""
	switch len(urlList) {
	case 0:
		message = "No images uploaded."
	case 1:
		message = "Uploaded 1 image."
	default:
		message = fmt.Sprintf("Uploaded %d images.", len(urlList))
	}

	_, err = msg.Raw.Session.ChannelMessageSend(msg.ChannelId, message)

	return
}
