package datastore

import (
	"io"
)

//DataStore interface, main interface for all datastores
type DataStore interface {
	Get(key string) (string, bool)
	Set(key string, value string) error
	Delete(key string) error
	Keys() ([]string, error)
	Values() ([]string, error)
}

//PersistentDataStore flushes to some storage medium
type PersistentDataStore interface {
	DataStore
	io.Closer
	Sync() error
}
