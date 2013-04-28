package gomc

import (
	"testing"
	"time"
)

func TestPoolBehavior(t *testing.T) {
	cmds := start(testHosts)
	defer stop(cmds)

	pool, err := newPool(testHosts, 1, 2, ENCODING_DEFAULT)
	if err != nil {
		t.Error("Fail to new client:", err)
	}

	if err = pool.SetBehavior(BEHAVIOR_DISTRIBUTION, uint64(DISTRIBUTION_CONSISTENT_KETAMA)); err != nil {
		t.Error("Fail to new client:", err)
	}

	if behavior, _ := pool.GetBehavior(BEHAVIOR_DISTRIBUTION); DistributionType(behavior) != DISTRIBUTION_CONSISTENT_KETAMA {
		t.Error("Error behavior:", behavior, ", expect:", DISTRIBUTION_CONSISTENT_KETAMA)
	}
}

func TestPoolSetGet(t *testing.T) {
	cmds := start(testHosts)
	defer stop(cmds)

	var (
		testKey   = "test-key"
		testValue = "test-value"
		testExpr  = time.Second
	)

	pool, err := newPool(testHosts, 1, 2, ENCODING_DEFAULT)
	if err != nil {
		t.Error("Fail to new client:", err)
	}

	if err = pool.Set(testKey, testValue, testExpr); err != nil {
		t.Error("Fail to set:", err)
	}

	var val string
	if err = pool.Get(testKey, &val); err != nil {
		t.Error("Fail to get:", err)
	} else if val != testValue {
		t.Error("Error get:", val, ", expect:", testValue)
	}

	time.Sleep(testExpr)

	if err = pool.Get(testKey, &val); err == nil && val == testValue {
		t.Error("Fail to expire")
	}
}

func BenchmarkPoolGet(b *testing.B) {
	b.StopTimer()

	cmds := start(testHosts)
	defer stop(cmds)

	pool, _ := newPool(testHosts, 1, 2, ENCODING_DEFAULT)
	testKey := "test-key"
	testValue := "test-value"
	restoreValue := new(string)

	pool.Set(testKey, testValue, 0)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		pool.Get(testKey, restoreValue)
	}
}
