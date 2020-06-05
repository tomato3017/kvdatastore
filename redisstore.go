package datastore

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	utilio "github.com/tomato3017/kvdatastore/util/io"
	"sync"
)

//RedisStore Redis Backed datastore
type RedisStore struct {
	pool   *redis.Pool
	prefix string
	lock   sync.RWMutex
}

func NewRedisStore(pool *redis.Pool, prefix string) (*RedisStore, error) {
	rds := RedisStore{
		pool:   pool,
		prefix: prefix,
	}

	if ok, err := rds.IsValidConnection(); !ok {
		return nil, err
	}

	return &rds, nil
}

func (r *RedisStore) IsValidConnection() (bool, error) {
	//Attempt a connection
	conn := r.pool.Get()
	defer utilio.SafeClose(conn)

	reply, err := redis.String(conn.Do("PING"))
	if err != nil {
		return false, fmt.Errorf("unable to make run command PING on redis. Err: %w", err)
	}

	if reply != "PONG" {
		return false, fmt.Errorf("unknown response to PING command. Response is %s", reply)
	}

	return true, nil
}

func (r *RedisStore) generateKey(key string) string {
	if r.prefix != "" {
		return fmt.Sprintf("%s--%s", r.prefix, key)
	}

	return key
}

func (r *RedisStore) Get(key string) (string, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	//Get a connection
	conn := r.pool.Get()
	defer utilio.SafeClose(conn)

	lookupKey := r.generateKey(key)

	rtnVal, err := redis.String(conn.Do("GET", lookupKey))
	if err != nil {
		return "", false
	}

	return rtnVal, true
}

func (r *RedisStore) Set(key string, value string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	//Get a connection
	conn := r.pool.Get()
	defer utilio.SafeClose(conn)
	derivedKey := r.generateKey(key)

	ok, err := conn.Do("SET", derivedKey, value)
	if ok != "OK" || err != nil {
		if err != nil {
			return fmt.Errorf("unable to set key %s. Err: %w", derivedKey, err)
		}
		return fmt.Errorf("unable to set key %s. REDIS Msg: %s", derivedKey, ok)
	}
	return nil
}

func (r *RedisStore) Delete(key string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	//Get a connection
	conn := r.pool.Get()
	defer utilio.SafeClose(conn)
	derivedKey := r.generateKey(key)

	_, err := conn.Do("DEL", derivedKey)
	if err != nil {
		return fmt.Errorf("unable to delete key %s. Err: %w", derivedKey, err)
	}
	return nil
}

func (r *RedisStore) Keys() ([]string, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	//Get a connection
	conn := r.pool.Get()
	defer utilio.SafeClose(conn)

	keys, err := redis.Strings(conn.Do("KEYS", r.generateKey("*")))
	if err != nil {
		return nil, fmt.Errorf("unable to get keys for prefix %s. Err: %w", r.prefix, err)
	}

	return keys, nil
}

func (r *RedisStore) Values() ([]string, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	keys, err := r.Keys()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve keys from redis. Err: %w", err)
	}

	conn := r.pool.Get()
	defer utilio.SafeClose(conn)

	vals := make([]string, 0)
	for _, key := range keys {
		val, err := redis.String(conn.Do("GET", key))
		if err == redis.ErrNil {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("unable to get key %s. Err: %w", key, err)
		}
		vals = append(vals, val)
	}

	return vals, nil
}
