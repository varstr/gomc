package gomc

import (
	"os/exec"
	"reflect"
    "strconv"
	"strings"
	"testing"
	"time"
)

const (
	_MC_CMD           = "memcached"
	_MC_FLAG_SOCKET   = "-s"
	_MC_FLAG_TCP_PORT = "-p"
)

var (
	testHosts = []string{
		"localhost:11211",
		"localhost:11212",
		"localhost:11213",
	}

	testSockets = []string{
		"/tmp/test-gomc-1.sock",
		"/tmp/test-gomc-2.sock",
		"/tmp/test-gomc-3.sock",
	}
)

func start(servers []string) (cmds []*exec.Cmd) {
	cmds = make([]*exec.Cmd, 0, len(servers))
	for _, server := range servers {
		var args []string
		if strings.HasPrefix(server, "/") {
			args = []string{_MC_FLAG_SOCKET, server}
		} else {
			port := strings.Split(server, ":")[1]
			args = []string{_MC_FLAG_TCP_PORT, port}
		}
		cmd := exec.Command(_MC_CMD, args...)
		cmd.Start()
		cmds = append(cmds, cmd)
	}
	return
}

func stop(cmds []*exec.Cmd) {
	for _, cmd := range cmds {
		cmd.Process.Kill()
	}
}

func TestBehavior(t *testing.T) {
	cmds := start(testHosts)
	defer stop(cmds)

	mc, err := newMemcached(testHosts, ENCODING_DEFAULT)
	if err != nil {
		t.Error("Fail to new client:", err)
	}

	if err = mc.SetBehavior(BEHAVIOR_DISTRIBUTION, uint64(DISTRIBUTION_CONSISTENT_KETAMA)); err != nil {
		t.Error("Fail to new client:", err)
	}

	if behavior, _ := mc.GetBehavior(BEHAVIOR_DISTRIBUTION); DistributionType(behavior) != DISTRIBUTION_CONSISTENT_KETAMA {
		t.Error("Error behavior:", behavior, ", expect:", DISTRIBUTION_CONSISTENT_KETAMA)
	}
}

func TestSetGet(t *testing.T) {
	cmds := start(testHosts)
	defer stop(cmds)

	testKey := "test-key"
	testValue := "test-value"
	testExpr := time.Second
	mc, err := newMemcached(testHosts, ENCODING_DEFAULT)
	if err != nil {
		t.Error("Fail to new client:", err)
	}

	if err = mc.Set(testKey, testValue, testExpr); err != nil {
		t.Error("Fail to set:", err)
	}

	var val string
	if err = mc.Get(testKey, &val); err != nil {
		t.Error("Fail to get:", err)
	} else if val != testValue {
		t.Error("Error get:", val, ", expect:", testValue)
	}

	time.Sleep(testExpr)

	if err = mc.Get(testKey, &val); err == nil && val == testValue {
		t.Error("Fail to expire")
	}
}

func testSetGetWithEncoding(t *testing.T, encoding EncodingType) {
	cmds := start(testHosts)
	defer stop(cmds)

	testKey := "test-key"
	testValue := randomStruct()
	restoreValue := new(TestStruct)
	mc, err := newMemcached(testHosts, encoding)
	if err != nil {
		t.Error("Fail to new client:", err)
	}
	if err = mc.Set(testKey, testValue, 0); err != nil {
		t.Error("Fail to set:", err)
	}

	if err = mc.Get(testKey, restoreValue); err != nil {
		t.Error("Fail to get:", err)
	} else if !reflect.DeepEqual(testValue, restoreValue) {
		t.Error("Error get:", restoreValue, ", expect:", testValue)
	}
}

func TestGobSetGet(t *testing.T) {
	testSetGetWithEncoding(t, ENCODING_GOB)
}

func TestJsonSetGet(t *testing.T) {
	testSetGetWithEncoding(t, ENCODING_JSON)
}

func TestGetMulti(t *testing.T) {
	cmds := start(testHosts)
	defer stop(cmds)

    num := 10
    testKeyPrefix := "test-key:"
    testKeys := make([]string, num)
    testStructs := make(map[string]*TestStruct, num)
	restoreValue := new(TestStruct)
	mc, err := newMemcached(testHosts, ENCODING_JSON)
	if err != nil {
		t.Error("Fail to new client:", err)
	}

    for i:=0; i<num; i++ {
        testKey := testKeyPrefix + strconv.Itoa(i)
        testValue := randomStruct()
        testKeys = append(testKeys, testKey)
        testStructs[testKey] = testValue
        if err = mc.Set(testKey, testValue, 0); err != nil {
            t.Error("Fail to set:", err)
        }
    }

    res, err := mc.GetMulti(testKeys)
    if err != nil {
		t.Error("Fail to get-multi:", err)
	}

    t.Log(res)

    for testKey, testValue := range testStructs {
        if err := res.Get(testKey, restoreValue); err != nil {
            t.Error("Fail to get:", err)
        } else if !reflect.DeepEqual(testValue, restoreValue) {
            t.Error("Error get:", restoreValue, ", expect:", testValue)
        }
    }
}

func TestDelete(t *testing.T) {
	cmds := start(testHosts)
	defer stop(cmds)

    testKey   := "test-key"
    testValue := "test-value"
	mc, err := newMemcached(testHosts, ENCODING_DEFAULT)
	if err != nil {
		t.Error("Fail to new client:", err)
	}

	if err = mc.Set(testKey, testValue, 0); err != nil {
		t.Error("Fail to set:", err)
	}

	if err = mc.Delete(testKey, 0); err != nil {
		t.Error("Fail to delete:", err)
	}

	var val string
	if err = mc.Get(testKey, &val); err == nil && val == testValue {
		t.Error("Fail to delete")
	}
}

func BenchmarkSet(b *testing.B) {
	b.StopTimer()

	cmds := start(testHosts)
	defer stop(cmds)

	mc, _ := newMemcached(testHosts, ENCODING_DEFAULT)
	testKey := "test-key"
	testValue := "test-value"

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		mc.Set(testKey, testValue, 0)
	}
}

func BenchmarkGet(b *testing.B) {
	b.StopTimer()

	cmds := start(testHosts)
	defer stop(cmds)

	mc, _ := newMemcached(testHosts, ENCODING_DEFAULT)
	testKey := "test-key"
	testValue := "test-value"
	restoreValue := new(string)

	mc.Set(testKey, testValue, 0)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		mc.Get(testKey, restoreValue)
	}
}
