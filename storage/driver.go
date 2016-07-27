package storage

// Driver is the interface for a storage engine.
type Driver interface {
	Open(uri string) (err error)
	Close() (err error)

	Get(key string) (value string, err error)
	Set(key, value string) (err error)
	Delete(key string) (err error)
}
