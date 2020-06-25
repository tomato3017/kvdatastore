package datastore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/spf13/afero"
	"github.com/tomato3017/kvdatastore/util/errcodes"
	utilio "github.com/tomato3017/kvdatastore/util/io"
)

type FileStore struct {
	storage      afero.Fs
	cache        map[string]string
	filename     string
	dirty        bool
	SaveOnChange bool
	lock         sync.RWMutex
}

func NewFileStore(storage afero.Fs, filename string, saveOnChange bool) (*FileStore, error) {
	//Create the file if it doesn't exist
	if err := utilio.TouchFile(storage, filename); err != nil {
		return nil, errcodes.ErrDataSourceDoesntExist{Name: filename}
	}

	//Open the file
	f, err := storage.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	byteArr, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	//Load the existing file into the cache if it exists
	cache := make(map[string]string)
	if len(byteArr) > 0 {
		if err := json.Unmarshal(byteArr, &cache); err != nil {
			return nil, fmt.Errorf("unable to unmarshal json for filestore. Err: %w", err)
		}
	}

	return &FileStore{storage: storage,
		filename:     filename,
		SaveOnChange: saveOnChange,
		cache:        cache}, nil
}

func (fs *FileStore) Get(key string) (string, bool) {
	fs.lock.RLock()
	defer fs.lock.RUnlock()

	val, ok := fs.cache[key]

	return val, ok
}

func (fs *FileStore) Set(key string, value string) error {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	fs.cache[key] = value

	//Set dirty
	fs.dirty = true

	if fs.SaveOnChange {
		return fs.syncToDiskRaw()
	}

	return nil
}

func (fs *FileStore) Delete(key string) error {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	delete(fs.cache, key)

	return nil
}

func (fs *FileStore) Keys() ([]string, error) {
	fs.lock.RLock()
	defer fs.lock.RUnlock()

	keys := make([]string, 0)

	for k := range fs.cache {
		keys = append(keys, k)
	}

	return keys, nil
}

func (fs *FileStore) Values() ([]string, error) {
	fs.lock.RLock()
	defer fs.lock.Unlock()

	values := make([]string, 0)

	for _, v := range fs.cache {
		values = append(values, v)
	}

	return values, nil
}

func (fs *FileStore) Close() error {
	return fs.syncToDisk()
}

func (fs *FileStore) Sync() error {
	return fs.syncToDisk()
}

//sync to disk without locking
func (fs *FileStore) syncToDiskRaw() error {
	//The cache isn't dirty, no need to save.
	if !fs.dirty {
		return nil
	}

	//Open the file for writing
	fData, err := fs.storage.OpenFile(fs.filename, os.O_WRONLY, 644)
	if err != nil {
		return err
	}
	defer fData.Close()

	//Convert map into json
	jsonBytes, err := json.Marshal(fs.cache)
	if err != nil {
		//TODO: Convert to error type
		return fmt.Errorf("unable to marshal json. Err: %w", err)
	}

	if _, err := fData.Write(jsonBytes); err != nil {
		return fmt.Errorf("unable to write data to file %s. Err: %w", fs.filename, err)
	}
	fs.dirty = false

	return nil
}

func (fs *FileStore) syncToDisk() error {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	return fs.syncToDiskRaw()
}
