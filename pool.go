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

type mcPool struct {
	pool     *C.memcached_pool_st
	encoding EncodingType
}

func newPool(servers []string, initSize, maxSize int) (self *mcPool, err error) {
	config := poolConfig(servers, initSize, maxSize)
	cs_config := C.CString(config)
	defer C.free(unsafe.Pointer(cs_config))

	self = new(mcPool)
	self.pool = C.memcached_pool(cs_config, C.size_t(len(config)))
	if self.pool == nil {
		err = self.checkError(
			C.libmemcached_check_configuration(
				cs_config, C.size_t(len(config)), nil, 0))
		return
	}
	return
}

func (self *mcPool) checkError(returnCode C.memcached_return_t) error {
	if C.memcached_failed(returnCode) {
		return errors.New(C.GoString(C.memcached_strerror(nil, returnCode)))
	}
	return nil
}

func (self *mcPool) SetBehavior(behavior BehaviorType, value uint64) error {
	return self.checkError(
		C.memcached_pool_behavior_set(
			self.pool, C.memcached_behavior_t(behavior), C.uint64_t(value)))
}

func (self *mcPool) GetBehavior(behavior BehaviorType) (value uint64, err error) {
	err = self.checkError(
		C.memcached_pool_behavior_get(
            self.pool, C.memcached_behavior_t(behavior), (*C.uint64_t)(&value)))
	return
}

func (self *mcPool) fetchConnection() (conn *memcached, err error) {
	ret := new(C.memcached_return_t)
	conn = &memcached{
		mmc:      C.memcached_pool_fetch(self.pool, nil, ret),
		encoding: self.encoding,
	}
	err = self.checkError(*ret)
	return
}

func (self *mcPool) releaseConnection(conn *memcached) error {
	return self.checkError(C.memcached_pool_release(self.pool, conn.mmc))
}

func (self *mcPool) GenerateHash(key string) (hash uint32, err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	hash, err = conn.GenerateHash(key)
	return
}

func (self *mcPool) Increment(key string, offset uint32) (value uint64, err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Increment(key, offset)
}

func (self *mcPool) Decrement(key string, offset uint32) (value uint64, err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Decrement(key, offset)
}

func (self *mcPool) Delete(key string, expiration time.Duration) (err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Delete(key, expiration)
}

func (self *mcPool) Exist(key string) (err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Exist(key)
}

func (self *mcPool) Flush(expiration time.Duration) (err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Flush(expiration)
}

func (self *mcPool) Get(key string, value interface{}) (err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Get(key, value)
}

func (self *mcPool) Add(key string, value interface{}, expiration time.Duration) (err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Add(key, value, expiration)
}

func (self *mcPool) Replace(key string, value interface{}, expiration time.Duration) (err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Replace(key, value, expiration)
}

func (self *mcPool) Set(key string, value interface{}, expiration time.Duration) (err error) {
	conn, err := self.fetchConnection()
	defer self.releaseConnection(conn)
	if err != nil {
		return
	}

	return conn.Set(key, value, expiration)
}
