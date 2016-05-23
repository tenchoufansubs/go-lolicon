package images

import (
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/tenchoufansubs/go-lolicon"
	"github.com/tenchoufansubs/go-lolicon/storage"
)

func init() {
	plugin := new(ImagesPlugin)

	var (
		_ lolicon.Plugin         = plugin
		_ lolicon.MessageHandler = plugin
	)

	lolicon.RegisterPlugin(plugin)
}

type ImagesPlugin struct {
	cache storage.Driver `json:"-"`

	Directories []string `json:"directories"`
}

func (p *ImagesPlugin) Id() lolicon.PluginId {
	return lolicon.PluginId("images")
}

func (p *ImagesPlugin) Setup(cache storage.Driver) (err error) {
	p.cache = cache
	return
}

func (p *ImagesPlugin) Open() (err error) {
	err = storage.GetJSON(p.cache, string(p.Id()), &p)
	if err != nil {
		return
	}

	if p.Directories == nil {
		p.Directories = make([]string, 0)
	}

	return
}

func (p *ImagesPlugin) Close() (err error) {
	err = storage.SetJSON(p.cache, string(p.Id()), p)
	return
}

func (p *ImagesPlugin) HandleMessage(msg *lolicon.Message) (done bool, err error) {
	if msg.Command == "" || msg.Trailing != "" {
		return
	}

	var (
		file string

		f *os.File
	)

	file, err = p.pickImageFile(msg.Command)
	if err != nil {
		return
	}
	if file == "" {
		return
	}

	done = true

	filename := path.Base(file)

	f, err = os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = msg.Raw.Session.ChannelFileSend(msg.ChannelId, filename, f)
	return
}

// pickImageFile returns the path to an image file.
//
// It looks in every directory in p.Directories for a subdirectory
// named dirname, and then returns a random image file from there.
//
// If no image can be found, file will be an empty string.
func (p *ImagesPlugin) pickImageFile(dirname string) (file string, err error) {
	if len(p.Directories) == 0 {
		return
	}

	for _, root := range p.Directories {
		file, err = randomImageFromDirectory(root, dirname)
		if err != nil {
			return
		}
		if file != "" {
			return
		}
	}

	return
}

// randomImageFromDirectory returns the path to a random image file
// in the directory path.Join(root, name).
//
// If the directory does not exist, or it exists and an image could
// not be found, imageFile will be an empty string and err will be nil.
func randomImageFromDirectory(root, name string) (imageFile string, err error) {
	var (
		fi os.FileInfo

		d *os.File

		entries []string
		files   []string
	)

	targetDir := path.Join(root, name)

	fi, err = os.Stat(targetDir)
	if err != nil {
		return
	}
	if !fi.IsDir() {
		return
	}

	d, err = os.Open(targetDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	defer d.Close()

	entries, err = d.Readdirnames(-1)
	if err != nil {
		return
	}

	for _, entry := range entries {
		file := path.Join(targetDir, entry)

		fi, err = os.Stat(file)
		if err != nil {
			return
		}

		// Ignore if it's not a normal file.
		if fi.Mode()&os.ModeType != 0 {
			continue
		}

		files = append(files, file)
	}

	if len(files) == 0 {
		return
	}

	rand.Seed(time.Now().UTC().UnixNano())

	idx := rand.Intn(len(files))
	imageFile = files[idx]

	return
}
