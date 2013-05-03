package gomc

import (
    "fmt"
)

/*
#cgo LDFLAGS: -L/usr/lib -lmemcached -lmemcachedutil
#include <libmemcached/memcached.h>
#include <libmemcached/util.h>
#include <stdlib.h>
#include <stdint.h>
*/
import "C"

type row struct {
    buffer []byte
    flags uint32
}

type result struct {
    raw map[string]*row
}

func newResult(size int) *result {
    return &result{raw: make(map[string]*row, size)}
}

func (self *result) set(key string, buffer []byte, flags uint32) {
    self.raw[key] = &row{
        buffer: buffer,
        flags: flags,
    }
}

func (self *result) Get(key string, value interface{}) (err error) {
    if row, ok := self.raw[key]; ok {
        return decode(row.buffer, row.flags, value)
    }
    err = fmt.Errorf("No result for key `%s`", key)
    return
}
