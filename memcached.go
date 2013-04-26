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
	"strings"
	"time"
	"unsafe"
)

const (
	_CONFIG_SERVER_PREFIX = "--SERVER="
	_CONFIG_SOCKET_PREFIX = "--SOCKET="
)

type ReturnType int
type BehaviorType int
type DistributionType int
type HashType int
type ConnectionType int

type mcClient struct {
	mc       *C.memcached_st
	encoding EncodingType
}

func NewClient(servers []string) (self *mcClient, err error) {
	cfg := make([]string, len(servers))
	for i, server := range servers {
		if strings.HasPrefix(server, "/") {
			cfg[i] = _CONFIG_SOCKET_PREFIX + server
		} else {
			cfg[i] = _CONFIG_SERVER_PREFIX + server
		}
	}
	config := strings.Join(cfg, " ")
	cs_config := C.CString(config)
	defer C.free(unsafe.Pointer(cs_config))

	self = new(mcClient)
	self.mc = C.memcached(cs_config, C.size_t(len(config)))
	if self.mc == nil {
		err = self.checkError(
			C.libmemcached_check_configuration(
				cs_config, C.size_t(len(config)), nil, 0))
		return
	}
	return
}

func (self *mcClient) encode(object interface{}) ([]byte, uint32, error) {
	return encode(object, self.encoding)
}

func (self *mcClient) decode(buffer []byte, flag uint32, object interface{}) error {
	return decode(buffer, flag, object)
}

func (self *mcClient) checkError(returnCode C.memcached_return_t) error {
	if C.memcached_failed(returnCode) {
		return errors.New(C.GoString(C.memcached_strerror(self.mc, returnCode)))
	}
	return nil
}

func (self *mcClient) LastErrorMessage() string {
	return C.GoString(C.memcached_last_error_message(self.mc))
}

func (self *mcClient) AddServer(host string, port int, weight uint32) error {
	cs_host := C.CString(host)
	defer C.free(unsafe.Pointer(cs_host))
	return self.checkError(
		C.memcached_server_add_with_weight(
			self.mc, cs_host, C.in_port_t(port), C.uint32_t(weight)))
}

func (self *mcClient) SetBehavior(behavior BehaviorType, value uint64) error {
	return self.checkError(
		C.memcached_behavior_set(
			self.mc, C.memcached_behavior_t(behavior), C.uint64_t(value)))
}

func (self *mcClient) GetBehavior(behavior BehaviorType) uint64 {
	return uint64(C.memcached_behavior_get(self.mc, C.memcached_behavior_t(behavior)))
}

func (self *mcClient) GenerateHash(key string) uint32 {
	cs_key := C.CString(key)
	defer C.free(unsafe.Pointer(cs_key))

	return uint32(C.memcached_generate_hash(self.mc, cs_key, C.size_t(len(key))))
}

func (self *mcClient) Increment(key string, offset uint32) (value uint64, err error) {
	cs_key := C.CString(key)
	defer C.free(unsafe.Pointer(cs_key))
	err = self.checkError(
		C.memcached_increment(
			self.mc, cs_key, C.size_t(len(key)), C.uint32_t(offset), (*C.uint64_t)(&value)))
	return
}

func (self *mcClient) Decrement(key string, offset uint32) (value uint64, err error) {
	cs_key := C.CString(key)
	defer C.free(unsafe.Pointer(cs_key))
	err = self.checkError(
		C.memcached_decrement(
			self.mc, cs_key, C.size_t(len(key)), C.uint32_t(offset), (*C.uint64_t)(&value)))
	return
}

func (self *mcClient) Delete(key string, expiration time.Duration) error {
	cs_key := C.CString(key)
	defer C.free(unsafe.Pointer(cs_key))
	return self.checkError(
		C.memcached_delete(
			self.mc, cs_key, C.size_t(len(key)), C.time_t(expiration.Seconds())))
}

func (self *mcClient) Exist(key string) error {
	cs_key := C.CString(key)
	defer C.free(unsafe.Pointer(cs_key))
	return self.checkError(C.memcached_exist(self.mc, cs_key, C.size_t(len(key))))
}

func (self *mcClient) FlushBuffers() error {
	return self.checkError(C.memcached_flush_buffers(self.mc))
}

func (self *mcClient) Flush(expiration time.Duration) error {
	return self.checkError(C.memcached_flush(self.mc, C.time_t(expiration.Seconds())))
}

func (self *mcClient) Get(key string, value interface{}) (err error) {
	flags := new(C.uint32_t)
	cErr := new(C.memcached_return_t)
	valueLen := new(C.size_t)
	cs_key := C.CString(key)
	defer C.free(unsafe.Pointer(cs_key))

	raw := C.memcached_get(self.mc, cs_key, C.size_t(len(key)), valueLen, flags, cErr)
	buffer := C.GoBytes(unsafe.Pointer(raw), C.int(*valueLen))
	if err = self.checkError(*cErr); err != nil {
		return
	}
	return decode(buffer, uint32(*flags), value)
}

func (self *mcClient) Add(key string, value interface{}, expiration time.Duration) (err error) {
	buffer, flag, err := self.encode(value)
	if err != nil {
		return
	}
	cs_key := C.CString(key)
	defer C.free(unsafe.Pointer(cs_key))
	cs_value := C.CString(string(buffer))
	defer C.free(unsafe.Pointer(cs_value))

	return self.checkError(
		C.memcached_add(
			self.mc, cs_key, C.size_t(len(key)), cs_value, C.size_t(len(buffer)),
			C.time_t(expiration.Seconds()), C.uint32_t(flag)))
}

func (self *mcClient) Replace(key string, value interface{}, expiration time.Duration) (err error) {
	buffer, flag, err := self.encode(value)
	if err != nil {
		return
	}
	cs_key := C.CString(key)
	defer C.free(unsafe.Pointer(cs_key))
	cs_value := C.CString(string(buffer))
	defer C.free(unsafe.Pointer(cs_value))

	return self.checkError(
		C.memcached_replace(
			self.mc, cs_key, C.size_t(len(key)), cs_value, C.size_t(len(buffer)),
			C.time_t(expiration.Seconds()), C.uint32_t(flag)))
}

func (self *mcClient) Set(key string, value interface{}, expiration time.Duration) (err error) {
	buffer, flag, err := self.encode(value)
	if err != nil {
		return
	}
	cs_key := C.CString(key)
	defer C.free(unsafe.Pointer(cs_key))
	cs_value := C.CString(string(buffer))
	defer C.free(unsafe.Pointer(cs_value))

	return self.checkError(
		C.memcached_set(
			self.mc, cs_key, C.size_t(len(key)), cs_value, C.size_t(len(buffer)),
			C.time_t(expiration.Seconds()), C.uint32_t(flag)))
}
