package datastore

import (
	"encoding/json"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"github.com/tomato3017/kvdatastore/util"
	"io/ioutil"
	"testing"
)

type TestFileStoreSuite struct {
	suite.Suite
	datastore *FileStore
	fs        afero.Fs
}

func (fss *TestFileStoreSuite) SetupTest() {
	fss.fs = afero.NewMemMapFs()

	var err error
	fss.datastore, err = NewFileStore(
		fss.fs,
		"/testfilename",
		false)
	fss.Require().NoError(err)
}

func (fss *TestFileStoreSuite) TestFileStoreGet() {
	req := fss.Require()
	keys := []string{
		"test1",
		"test2",
		"test3",
		"helloworld",
	}

	for _, key := range keys {
		req.NoError(fss.datastore.Set(key, "testvaluefor_"+key))
	}

	val, ok := fss.datastore.Get(keys[1])
	req.True(ok)
	req.Equal("testvaluefor_"+keys[1], val)
}

func (fss *TestFileStoreSuite) TestFileStoreSet() {
	req := fss.Require()
	keys := []string{
		"test1",
		"test2",
		"test3",
		"helloworld",
	}

	for _, key := range keys {
		req.NoError(fss.datastore.Set(key, "testvaluefor_"+key))
	}

	//Checks if the dirty bit was set
	req.True(fss.datastore.dirty)

	//check if a map value exists
	req.Equal("testvaluefor_"+keys[1], fss.datastore.cache[keys[1]])

	req.NoError(fss.datastore.Sync())

	//Check if the dirty bit was unset
	req.False(fss.datastore.dirty)

	req.NoError(fss.datastore.Close())

	//Next lets load the file
	f, err := fss.fs.Open("/testfilename")
	req.NoError(err)
	defer f.Close()

	byteArr, err := ioutil.ReadAll(f)
	req.NoError(err)
	jsonMap := make(map[string]interface{})
	req.NoError(json.Unmarshal(byteArr, &jsonMap))

	value, ok := jsonMap[keys[1]]
	req.True(ok)

	valStr, ok := value.(string)
	req.True(ok)
	req.Equal("testvaluefor_"+keys[1], valStr)
}

func (fss *TestFileStoreSuite) TestKeys() {
	req := fss.Require()
	keys := []string{
		"testworld",
		"testone",
		"testgroup",
	}

	for _, v := range keys {
		err := fss.datastore.Set(v, "helloworld")
		req.NoError(err)
	}

	keysDS, err := fss.datastore.Keys()
	req.NoError(err)
	req.True(util.EqualStringSlices(keys, keysDS))
}

func TestFileStore(t *testing.T) {
	suite.Run(t, new(TestFileStoreSuite))
}
