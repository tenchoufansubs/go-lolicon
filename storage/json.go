package storage

type JSONGetter interface {
	GetJSON(key string, v interface{}) (err error)
}

type JSONSetter interface {
	SetJSON(key string, value interface{}) (err error)
}
