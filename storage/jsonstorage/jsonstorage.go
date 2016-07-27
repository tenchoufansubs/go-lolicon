// Package jsonstorage implements the storage.Driver interface for persisting
// information in JSON files.
package jsonstorage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/tenchoufansubs/go-lolicon/storage"
)

type JSONStorage struct {
	sync.RWMutex

	file string
	data map[string]interface{}
}

func (s *JSONStorage) Open(uri string) (err error) {
	var (
		f *os.File

		rawData []byte
	)

	// Convert file uri to path.
	if strings.HasPrefix(uri, "file://") {
		uri = strings.TrimPrefix(uri, "file://")
		uri, err = url.QueryUnescape(uri)
		if err != nil {
			return
		}
	}

	s.file = uri
	s.data = make(map[string]interface{})

	s.Lock()
	defer s.Unlock()

	f, err = os.Open(s.file)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}

		f, err = os.OpenFile(s.file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
		if err != nil {
			return
		}
		f.Write([]byte("{}"))
		f.Close()
		return
	}
	defer f.Close()

	rawData, err = ioutil.ReadAll(f)
	if err != nil {
		return
	}

	err = json.Unmarshal(rawData, &s.data)
	if err != nil {
		return
	}

	return
}

func (s *JSONStorage) Close() (err error) {
	err = s.Flush()
	return
}

func (s *JSONStorage) Flush() (err error) {
	var (
		rawData []byte

		f *os.File
	)

	s.Lock()
	defer s.Unlock()

	rawData, err = json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return
	}

	f, err = os.OpenFile(s.file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.Write(rawData)
	if err != nil {
		return
	}

	_, err = f.Write([]byte{'\n'})

	return
}

func (s *JSONStorage) Get(key string) (value string, err error) {
	var (
		ok bool

		rawValue interface{}
	)

	s.RLock()
	defer s.RUnlock()

	rawValue, ok = s.data[key]
	if !ok {
		err = storage.ErrNotFound
		return
	}

	value, ok = rawValue.(string)
	if !ok {
		err = errors.New("value is not a string")
		return
	}

	return
}

func (s *JSONStorage) Set(key, value string) (err error) {
	s.Lock()

	s.data[key] = value

	s.Unlock()

	err = s.Flush()
	return
}

func (s *JSONStorage) GetJSON(key string, v interface{}) (err error) {
	s.RLock()
	defer s.RUnlock()

	var (
		rawData []byte
	)

	rawData, err = json.Marshal(s.data[key])
	if err != nil {
		return
	}

	err = json.Unmarshal(rawData, v)
	return
}

func (s *JSONStorage) SetJSON(key string, value interface{}) (err error) {
	s.Lock()

	s.data[key] = value

	s.Unlock()

	err = s.Flush()
	return
}

func (s *JSONStorage) Delete(key string) (err error) {
	s.Lock()

	delete(s.data, key)

	s.Unlock()

	err = s.Flush()
	return
}

func New() (storage storage.Driver, err error) {
	storage = new(JSONStorage)
	return
}
