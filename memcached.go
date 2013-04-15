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
)

const (
    _CONFIG_SERVER_PREFIX = "--SERVER="
    _CONFIG_SOCKET_PREFIX = "--SOCKET="
)

type mcClient struct {
    mc *C.memcached_st
    encoding EncodingType
}

func NewClient(servers []string) (self *mcClient, err error) {
    for i, server := range servers {
        if strings.HasPrefix(server, "/") {
            servers[i] = _CONFIG_SOCKET_PREFIX + server
        } else {
            servers[i] = _CONFIG_SERVER_PREFIX + server
        }
    }
    config := strings.Join(servers, " ")
    self = new(mcClient)
    self.mc = C.memcached(C.CString(config), C.size_t(len(config)))
    if self.mc == nil {
        err = self.checkError(C.libmemcached_check_configuration(C.CString(config), C.size_t(len(config)), nil, 0))
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
        return nil
    }
    return errors.New(C.GoString(C.memcached_strerror(self.mc, returnCode)))
}

func (self *mcClient) LastErrorMessage() string {
    return C.GoString(C.memcached_last_error_message(self.mc))
}

func (self *mcClient) AddServer(host string, port int, weight uint32) error {
    return self.checkError(C.memcached_server_add_with_weight(self.mc, C.CString(host), C.in_port_t(port), C.uint32_t(weight)))
}

func (self *mcClient) SetBehavior(behavior C.memcached_behavior_t, value uint64) error {
    return self.checkError(C.memcached_behavior_set(self.mc, behavior, C.uint64_t(value)))
}

func (self *mcClient) GetBehavior(behavior C.memcached_behavior_t) uint64 {
    return uint64(C.memcached_behavior_get(self.mc, behavior))
}

func (self *mcClient) Increment(key string, offset uint32) (value uint64, err error) {
    err = self.checkError(C.memcached_increment(self.mc, C.CString(key), C.size_t(len(key)), C.uint32_t(offset), (*C.uint64_t)(&value)))
    return
}

func (self *mcClient) Decrement(key string, offset uint32) (value uint64, err error) {
    err = self.checkError(C.memcached_decrement(self.mc, C.CString(key), C.size_t(len(key)), C.uint32_t(offset), (*C.uint64_t)(&value)))
    return
}

func (self *mcClient) Delete(key string, expiration time.Duration) error {
    return self.checkError(C.memcached_delete(self.mc, C.CString(key), C.size_t(len(key)), C.time_t(expiration.Seconds())))
}

func (self *mcClient) Exist(key string) error {
    return self.checkError(C.memcached_exist(self.mc, C.CString(key), C.size_t(len(key))))
}

func (self *mcClient) FlushBuffers() error {
    return self.checkError(C.memcached_flush_buffers(self.mc))
}

func (self *mcClient) Flush(expiration time.Duration) error {
    return self.checkError(C.memcached_flush(self.mc, C.time_t(expiration.Seconds())))
}

func (self *mcClient) Get(key string, value interface{}) error {
    return nil
}

func (self *mcClient) Add(key string, value interface{}, expiration time.Duration) (err error) {
    buffer, flag, err := self.encode(value)
    if err != nil {
        return
    }
    return self.checkError(C.memcached_add(self.mc, C.CString(key), C.size_t(len(key)), C.CString(string(buffer)), C.size_t(len(buffer)), C.time_t(expiration.Seconds()), C.uint32_t(flag)))
}

func (self *mcClient) Replace(key string, value interface{}, expiration time.Duration) (err error) {
    buffer, flag, err := self.encode(value)
    if err != nil {
        return
    }
    return self.checkError(C.memcached_replace(self.mc, C.CString(key), C.size_t(len(key)), C.CString(string(buffer)), C.size_t(len(buffer)), C.time_t(expiration.Seconds()), C.uint32_t(flag)))
}

func (self *mcClient) Set(key string, value interface{}, expiration time.Duration) (err error) {
    buffer, flag, err := self.encode(value)
    if err != nil {
        return
    }
    return self.checkError(C.memcached_set(self.mc, C.CString(key), C.size_t(len(key)), C.CString(string(buffer)), C.size_t(len(buffer)), C.time_t(expiration.Seconds()), C.uint32_t(flag)))
}

