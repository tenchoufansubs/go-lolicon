// Package storage defines interfaces for persisting data in storage engines.
package storage

import (
	"encoding/json"
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

// GetJSON is a helper for retrieving and unmarshalling a json-encoded string
// from a storage engine.
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

// SetJSON is a helper for storing arbitrary information in a storage engine as
// a json-encoded string.
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
