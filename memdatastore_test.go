package datastore

import (
	"github.com/stretchr/testify/suite"
	"github.com/tomato3017/kvdatastore/util"
	"testing"
)

type MemDataStoreTestSuite struct {
	suite.Suite
	datastore MemoryDataStore
}

func (m *MemDataStoreTestSuite) SetupTest() {
	m.datastore = NewMemoryDataStore()
}

func (m *MemDataStoreTestSuite) TestGet() {
	req := m.Require()
	m.datastore.storage["testkey"] = "teststring"

	val, ok := m.datastore.Get("testkey")
	req.True(ok)
	req.IsType("", val)
	req.Equal("teststring", val)

	_, ok = m.datastore.Get("adsasadas")
	req.False(ok)
}

func (m *MemDataStoreTestSuite) TestSet() {
	req := m.Require()

	err := m.datastore.Set("testkey", "teststring")
	req.NoError(err)
	req.Equal("teststring", m.datastore.storage["testkey"])
}

func (m *MemDataStoreTestSuite) TestKeys() {
	req := m.Require()
	keys := []string{
		"testworld",
		"testone",
		"testgroup",
	}

	for _, v := range keys {
		err := m.datastore.Set(v, "helloworld")
		req.NoError(err)
	}

	keysDS, err := m.datastore.Keys()
	req.NoError(err)
	req.True(util.EqualStringSlices(keys, keysDS))
}

func TestRunMemDataStoreSuite(t *testing.T) {
	suite.Run(t, new(MemDataStoreTestSuite))
}
