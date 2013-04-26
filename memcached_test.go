package gomc

import (
	"os/exec"
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

func start(servers []string, t *testing.T) (cmds []*exec.Cmd) {
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
		if err := cmd.Start(); err != nil {
			t.Error("Fail to start:", err)
		}
		cmds = append(cmds, cmd)
	}
	return
}

func stop(cmds []*exec.Cmd, t *testing.T) {
	for _, cmd := range cmds {
		if err := cmd.Process.Kill(); err != nil {
			t.Error("Fail to kill:", err)
		}
	}
}

func TestBehavior(t *testing.T) {
	cmds := start(testHosts, t)
	defer stop(cmds, t)

	cli, err := NewClient(testHosts)
	if err != nil {
		t.Error("Fail to new client:", err)
	}

	if err = cli.SetBehavior(BEHAVIOR_DISTRIBUTION, uint64(DISTRIBUTION_CONSISTENT_KETAMA)); err != nil {
		t.Error("Fail to new client:", err)
	}

	if behavior := cli.GetBehavior(BEHAVIOR_DISTRIBUTION); DistributionType(behavior) != DISTRIBUTION_CONSISTENT_KETAMA {
		t.Error("Error behavior:", behavior, ", expect:", DISTRIBUTION_CONSISTENT_KETAMA)
	}
}

func TestSetGet(t *testing.T) {
	cmds := start(testHosts, t)
	defer stop(cmds, t)

	var (
		testKey   = "test-key"
		testValue = "test-value"
		testExpr  = time.Second
	)

	cli, err := NewClient(testHosts)
	if err != nil {
		t.Error("Fail to new client:", err)
	}

	if err = cli.Set(testKey, testValue, testExpr); err != nil {
		t.Error("Fail to set:", err)
	}

	var val string
	if err = cli.Get(testKey, &val); err != nil {
		t.Error("Fail to get:", err)
	} else if val != testValue {
		t.Error("Error get:", val, ", expect:", testValue)
	}

	time.Sleep(testExpr)

	if err = cli.Get(testKey, &val); err == nil && val == testValue {
		t.Error("Fail to expire")
	}
}

func TestDelete(t *testing.T) {
	cmds := start(testHosts, t)
	defer stop(cmds, t)

	var (
		testKey   = "test-key"
		testValue = "test-value"
	)

	cli, err := NewClient(testHosts)
	if err != nil {
		t.Error("Fail to new client:", err)
	}

	if err = cli.Set(testKey, testValue, 0); err != nil {
		t.Error("Fail to set:", err)
	}

	if err = cli.Delete(testKey, 0); err != nil {
		t.Error("Fail to delete:", err)
	}

	var val string
	if err = cli.Get(testKey, &val); err == nil && val == testValue {
		t.Error("Fail to delete")
	}
}
