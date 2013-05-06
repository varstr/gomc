package gomc

/*
#cgo LDFLAGS: -L/usr/lib -lmemcached -lmemcachedutil
#include <libmemcached/memcached.h>
#include <libmemcached/util.h>
#include <stdlib.h>
#include <stdint.h>
*/
import "C"

import (
	"errors"
	"time"
	"unsafe"
)

type memcachedPool struct {
	pool     *C.memcached_pool_st
	encoding EncodingType
}

func newPool(servers []string, initSize, maxSize int, encoding EncodingType) (self *memcachedPool, err error) {
	config := poolConfig(servers, initSize, maxSize)
	cs_config := C.CString(config)
	defer C.free(unsafe.Pointer(cs_config))

	self = new(memcachedPool)
	self.pool = C.memcached_pool(cs_config, C.size_t(len(config)))
	if self.pool == nil {
		err = self.checkError(
			C.libmemcached_check_configuration(
				cs_config, C.size_t(len(config)), nil, 0))
		return
	}
	self.encoding = encoding
	return
}

func (self *memcachedPool) checkError(returnCode C.memcached_return_t) error {
	if C.memcached_failed(returnCode) {
		return errors.New(C.GoString(C.memcached_strerror(nil, returnCode)))
	}
	return nil
}

func (self *memcachedPool) SetBehavior(behavior BehaviorType, value uint64) error {
	return self.checkError(
		C.memcached_pool_behavior_set(
			self.pool, C.memcached_behavior_t(behavior), C.uint64_t(value)))
}

func (self *memcachedPool) GetBehavior(behavior BehaviorType) (value uint64, err error) {
	err = self.checkError(
		C.memcached_pool_behavior_get(
			self.pool, C.memcached_behavior_t(behavior), (*C.uint64_t)(&value)))
	return
}

func (self *memcachedPool) fetchConnection() (conn *memcached, err error) {
	ret := new(C.memcached_return_t)
	conn = &memcached{
		mmc:      C.memcached_pool_fetch(self.pool, nil, ret),
		encoding: self.encoding,
	}
	err = self.checkError(*ret)
	return
}

func (self *memcachedPool) releaseConnection(conn *memcached) error {
	return self.checkError(C.memcached_pool_release(self.pool, conn.mmc))
}

func (self *memcachedPool) GenerateHash(key string) (hash uint32, err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	hash, err = conn.GenerateHash(key)
	return
}

func (self *memcachedPool) Increment(key string, offset uint32) (value uint64, err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Increment(key, offset)
}

func (self *memcachedPool) Decrement(key string, offset uint32) (value uint64, err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Decrement(key, offset)
}

func (self *memcachedPool) Delete(key string, expiration time.Duration) (err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Delete(key, expiration)
}

func (self *memcachedPool) Exist(key string) (err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Exist(key)
}

func (self *memcachedPool) Flush(expiration time.Duration) (err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Flush(expiration)
}

func (self *memcachedPool) Get(key string, value interface{}) (err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Get(key, value)
}

func (self *memcachedPool) GetMulti(keys []string) (res Result, err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.GetMulti(keys)
}

func (self *memcachedPool) Add(key string, value interface{}, expiration time.Duration) (err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Add(key, value, expiration)
}

func (self *memcachedPool) Replace(key string, value interface{}, expiration time.Duration) (err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Replace(key, value, expiration)
}

func (self *memcachedPool) Set(key string, value interface{}, expiration time.Duration) (err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Set(key, value, expiration)
}
