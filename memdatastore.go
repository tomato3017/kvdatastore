package datastore

import (
	"github.com/tomato3017/kvdatastore/util/errcodes"
	"sync"
)

//MemoryDataStore Memory based datastore
type MemoryDataStore struct {
	storage map[string]string
	lock    sync.RWMutex
}

func NewMemoryDataStore() MemoryDataStore {
	return MemoryDataStore{storage: make(map[string]string)}
}

func (m *MemoryDataStore) Get(key string) (string, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if val, ok := m.storage[key]; ok {
		return val, ok
	}

	return "", false
}

func (m *MemoryDataStore) Set(key string, value string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.storage[key] = value
	return nil
}

func (m *MemoryDataStore) Delete(key string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.storage[key]; ok {
		delete(m.storage, key)
		return nil
	}

	return errcodes.ErrMissingKey{
		Key: key,
	}
}

func (m *MemoryDataStore) Keys() ([]string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	keys := make([]string, 0)

	for k := range m.storage {
		keys = append(keys, k)
	}

	return keys, nil
}

func (m *MemoryDataStore) Values() ([]string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	values := make([]string, 0)

	for _, v := range m.storage {
		values = append(values, v)
	}

	return values, nil
}
