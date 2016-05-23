package storage

import (
	"encoding/json"
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

func GetJSON(s Driver, key string, v interface{}) (err error) {
	getter, ok := s.(JSONGetter)
	if ok {
		err = getter.GetJSON(key, v)
		return
	}

	var (
		rawValueStr string
	)

	rawValueStr, err = s.Get(key)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(rawValueStr), v)
	return
}

func SetJSON(s Driver, key string, value interface{}) (err error) {
	setter, ok := s.(JSONSetter)
	if ok {
		err = setter.SetJSON(key, value)
		return
	}

	var (
		rawValue []byte
	)

	rawValue, err = json.Marshal(value)
	if err != nil {
		return
	}

	err = s.Set(key, string(rawValue))
	return
}
