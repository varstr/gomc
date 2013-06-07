package gomc

import (
	"time"
)

type Result interface {
	Get(string, interface{}) error
}

type Client interface {
	SetBehavior(BehaviorType, uint64) error
	GetBehavior(BehaviorType) (uint64, error)
	GenerateHash(string) (uint32, error)
	Increment(string, uint32) (uint64, error)
	Decrement(string, uint32) (uint64, error)
	Delete(string, time.Duration) error
	Exist(string) error
	Flush(time.Duration) error
	Get(string, interface{}) error
	GetMulti([]string) (Result, error)
	Add(string, interface{}, time.Duration) error
	Replace(string, interface{}, time.Duration) error
	Set(string, interface{}, time.Duration) error
	Close()
}

func NewClient(servers []string, poolSize int, encoding EncodingType) (self Client, err error) {
	if poolSize <= 1 {
		return newMemcached(servers, encoding)
	}
	return newPool(servers, 1, poolSize, encoding)
}
