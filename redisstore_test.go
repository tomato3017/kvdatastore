package datastore

import (
	"fmt"
	"github.com/alicebob/miniredis"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/suite"
	"github.com/tomato3017/kvdatastore/util"
	"math/rand"
	"strconv"
	"testing"
)

type RedisTestSuite struct {
	suite.Suite
	rds   *miniredis.Miniredis
	rpool *redis.Pool
}

func (r *RedisTestSuite) TearDownTest() {
	r.NoError(r.rpool.Close())
	r.rds.Close()
}

func (r *RedisTestSuite) SetupTest() {
	minirds, err := miniredis.Run()
	r.Require().NoError(err)

	r.rds = minirds

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", r.rds.Addr())
			if err != nil {
				return nil, fmt.Errorf("unable to make redis connection. Err: %w", err)
			}

			return conn, nil
		},
	}

	r.rpool = pool
}

func (r *RedisTestSuite) getRedisStore() (*RedisStore, error) {
	return NewRedisStore(r.rpool, "testprefix")
}

func (r *RedisTestSuite) TestNewRedisStore() {
	rStore, err := r.getRedisStore()
	r.Require().NoError(err)

	r.Require().Equal(rStore.prefix, "testprefix")
}

func (r *RedisTestSuite) TestRedisGet() {
	req := r.Require()
	rstore, err := r.getRedisStore()
	req.NoError(err)

	//Lets set some keys
	err = r.rds.Set(rstore.generateKey("test1"), "testkey")
	err = r.rds.Set(rstore.generateKey("test2"), "testkey2")
	req.NoError(err)

	keyI, ok := rstore.Get("test1")
	req.True(ok)

	req.Equal("testkey", keyI)

	_, ok = rstore.Get("nonexistentkey")
	req.False(ok)

}

func (r *RedisTestSuite) TestRedisSet() {
	req := r.Require()
	rstore, err := r.getRedisStore()
	req.NoError(err)

	req.NoError(rstore.Set("testkey1", "anothertestvalue"))
	print(r.rds.Dump())
	val, err := r.rds.Get(rstore.generateKey("testkey1"))
	req.NoError(err)

	req.Equal(val, "anothertestvalue")
}

func (r *RedisTestSuite) TestRedisDelete() {
	req := r.Require()

	rstore, err := r.getRedisStore()
	req.NoError(err)

	generatedKey := rstore.generateKey("testkey1")
	req.NoError(r.rds.Set(generatedKey, "someothertest"))
	req.NoError(r.rds.Set(generatedKey+"123", "someothertest123"))
	req.NoError(rstore.Delete("testkey1"))

	req.False(r.rds.Exists(generatedKey))
	req.True(r.rds.Exists(generatedKey + "123"))

}

func (r *RedisTestSuite) TestKeys() {
	req := r.Require()

	rstore, err := r.getRedisStore()
	req.NoError(err)

	keys, _ := r.generateKeysInRedis(rstore.generateKey(""), 10)

	keysDS, err := rstore.Keys()
	req.True(util.EqualStringSlices(keys, keysDS))
}

func (r *RedisTestSuite) TestValues() {
	req := r.Require()

	rstore, err := r.getRedisStore()
	req.NoError(err)

	_, vals := r.generateKeysInRedis(rstore.generateKey(""), 10)

	testVals, err := rstore.Values()
	req.NoError(err)
	req.True(util.EqualStringSlices(vals, testVals))
}

func (r *RedisTestSuite) generateKeysInRedis(prefix string, count int) ([]string, []string) {
	req := r.Require()

	keys := make([]string, count)
	for i := 0; i < count; i++ {
		keys[i] = fmt.Sprintf("%srandomkey-%d", prefix, i)
	}

	vals := make([]string, count)
	for i, v := range keys {
		vals[i] = strconv.Itoa(rand.Int())
		req.NoError(r.rds.Set(v, vals[i]))
	}

	return keys, vals
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
