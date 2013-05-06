package gomc

/*
#cgo LDFLAGS: -L/usr/lib -lmemcached
#include <libmemcached/memcached.h>
#include <stdlib.h>
#include <stdint.h>
*/
import "C"

import (
	"errors"
	"time"
	"unsafe"
)

type ReturnType int
type BehaviorType int
type DistributionType int
type HashType int
type ConnectionType int

func cString(str string) (*C.char, C.size_t) {
	return C.CString(str), C.size_t(len(str))
}

type memcached struct {
	mc       *C.memcached_st
	encoding EncodingType
}

func newMemcached(servers []string, encoding EncodingType) (self *memcached, err error) {
	config := clientConfig(servers)
	cs_config, config_len := cString(config)
	defer C.free(unsafe.Pointer(cs_config))

	self = new(memcached)
	self.mc = C.memcached(cs_config, config_len)
	if self.mc == nil {
		err = self.checkError(
			C.libmemcached_check_configuration(
				cs_config, config_len, nil, 0))
		return
	}
	self.encoding = encoding
	return
}

func (self *memcached) encode(object interface{}) ([]byte, uint32, error) {
	return encode(object, self.encoding)
}

func (self *memcached) checkError(returnCode C.memcached_return_t) error {
	if C.memcached_failed(returnCode) {
		return errors.New(C.GoString(C.memcached_strerror(self.mc, returnCode)))
	}
	return nil
}

func (self *memcached) LastErrorMessage() string {
	return C.GoString(C.memcached_last_error_message(self.mc))
}

func (self *memcached) AddServer(host string, port int, weight uint32) error {
	cs_host := C.CString(host)
	defer C.free(unsafe.Pointer(cs_host))
	return self.checkError(
		C.memcached_server_add_with_weight(
			self.mc, cs_host, C.in_port_t(port), C.uint32_t(weight)))
}

func (self *memcached) SetBehavior(behavior BehaviorType, value uint64) error {
	return self.checkError(
		C.memcached_behavior_set(
			self.mc, C.memcached_behavior_t(behavior), C.uint64_t(value)))
}

func (self *memcached) GetBehavior(behavior BehaviorType) (uint64, error) {
	return uint64(C.memcached_behavior_get(self.mc, C.memcached_behavior_t(behavior))), nil
}

func (self *memcached) GenerateHash(key string) (uint32, error) {
	cs_key, key_len := cString(key)
	defer C.free(unsafe.Pointer(cs_key))

	return uint32(C.memcached_generate_hash(self.mc, cs_key, key_len)), nil
}

func (self *memcached) Increment(key string, offset uint32) (value uint64, err error) {
	cs_key, key_len := cString(key)
	defer C.free(unsafe.Pointer(cs_key))
	err = self.checkError(
		C.memcached_increment(
			self.mc, cs_key, key_len, C.uint32_t(offset), (*C.uint64_t)(&value)))
	return
}

func (self *memcached) Decrement(key string, offset uint32) (value uint64, err error) {
	cs_key, key_len := cString(key)
	defer C.free(unsafe.Pointer(cs_key))
	err = self.checkError(
		C.memcached_decrement(
			self.mc, cs_key, key_len, C.uint32_t(offset), (*C.uint64_t)(&value)))
	return
}

func (self *memcached) Delete(key string, expiration time.Duration) error {
	cs_key, key_len := cString(key)
	defer C.free(unsafe.Pointer(cs_key))
	return self.checkError(
		C.memcached_delete(
			self.mc, cs_key, key_len, C.time_t(expiration.Seconds())))
}

func (self *memcached) Exist(key string) error {
	cs_key, key_len := cString(key)
	defer C.free(unsafe.Pointer(cs_key))
	return self.checkError(C.memcached_exist(self.mc, cs_key, key_len))
}

func (self *memcached) FlushBuffers() error {
	return self.checkError(C.memcached_flush_buffers(self.mc))
}

func (self *memcached) Flush(expiration time.Duration) error {
	return self.checkError(C.memcached_flush(self.mc, C.time_t(expiration.Seconds())))
}

func (self *memcached) Get(key string, value interface{}) (err error) {
	flags := new(C.uint32_t)
	ret := new(C.memcached_return_t)
	value_len := new(C.size_t)
	cs_key, key_len := cString(key)
	defer C.free(unsafe.Pointer(cs_key))

	raw := C.memcached_get(self.mc, cs_key, key_len, value_len, flags, ret)
	buffer := C.GoBytes(unsafe.Pointer(raw), C.int(*value_len))
	if err = self.checkError(*ret); err != nil {
		return
	}
	return decode(buffer, uint32(*flags), value)
}

func (self *memcached) getMulti(keys []string) (res *result, err error) {
	char_size := unsafe.Sizeof(new(C.char))
	cs_keys := C.malloc(C.size_t(len(keys)) * C.size_t(char_size))
	defer C.free(cs_keys)

	len_size := unsafe.Sizeof(C.size_t(0))
	key_lens := C.malloc(C.size_t(len(keys)) * C.size_t(len_size))
	defer C.free(key_lens)

	for i, key := range keys {
		cs_key, key_len := cString(key)
		defer C.free(unsafe.Pointer(cs_key))

		key_pos := (**C.char)(unsafe.Pointer(uintptr(cs_keys) + uintptr(i)*char_size))
		*key_pos = cs_key

		size_pos := (*C.size_t)(unsafe.Pointer(uintptr(key_lens) + uintptr(i)*len_size))
		*size_pos = key_len
	}

	ret := C.memcached_mget(self.mc, (**C.char)(cs_keys), (*C.size_t)(key_lens), C.size_t(len(keys)))
	if err = self.checkError(ret); err != nil {
		return
	}

	rc := new(C.memcached_return_t)
	res = newResult(len(keys))
	for {
		if raw := C.memcached_fetch_result(self.mc, nil, rc); raw != nil && ReturnType(*rc) != END {
			key := C.memcached_result_key_value(raw)
			buffer := C.memcached_result_value(raw)
			buffer_len := C.memcached_result_length(raw)
			flags := C.memcached_result_flags(raw)
			res.set(C.GoString(key), C.GoBytes(unsafe.Pointer(buffer), C.int(buffer_len)), uint32(flags))
		} else {
			break
		}
	}
	return
}

func (self *memcached) GetMulti(keys []string) (Result, error) {
	return self.getMulti(keys)
}

func (self *memcached) Add(key string, value interface{}, expiration time.Duration) (err error) {
	buffer, flag, err := self.encode(value)
	if err != nil {
		return
	}
	cs_key, key_len := cString(key)
	defer C.free(unsafe.Pointer(cs_key))
	cs_value, value_len := cString(string(buffer))
	defer C.free(unsafe.Pointer(cs_value))

	return self.checkError(
		C.memcached_add(
			self.mc, cs_key, key_len, cs_value, value_len,
			C.time_t(expiration.Seconds()), C.uint32_t(flag)))
}

func (self *memcached) Replace(key string, value interface{}, expiration time.Duration) (err error) {
	buffer, flag, err := self.encode(value)
	if err != nil {
		return
	}
	cs_key, key_len := cString(key)
	defer C.free(unsafe.Pointer(cs_key))
	cs_value, value_len := cString(string(buffer))
	defer C.free(unsafe.Pointer(cs_value))

	return self.checkError(
		C.memcached_replace(
			self.mc, cs_key, key_len, cs_value, value_len,
			C.time_t(expiration.Seconds()), C.uint32_t(flag)))
}

func (self *memcached) Set(key string, value interface{}, expiration time.Duration) (err error) {
	buffer, flag, err := self.encode(value)
	if err != nil {
		return
	}
	cs_key, key_len := cString(key)
	defer C.free(unsafe.Pointer(cs_key))
	cs_value, value_len := cString(string(buffer))
	defer C.free(unsafe.Pointer(cs_value))

	return self.checkError(
		C.memcached_set(
			self.mc, cs_key, key_len, cs_value, value_len,
			C.time_t(expiration.Seconds()), C.uint32_t(flag)))
}
