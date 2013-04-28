package gomc

import (
	"strconv"
	"strings"
)

const (
	_CONFIG_SEPARATOR = " "

	_CONFIG_SERVER_PREFIX = "--SERVER="
	_CONFIG_SOCKET_PREFIX = "--SOCKET="

	_CONFIG_POOL_MIN = "--POOL-MIN="
	_CONFIG_POOL_MAX = "--POOL-MAX="
)

func join(options []string) string {
	return strings.Join(options, _CONFIG_SEPARATOR)
}

func clientOptions(servers []string) (options []string) {
	options = make([]string, len(servers))
	for i, server := range servers {
		if strings.HasPrefix(server, "/") {
			options[i] = _CONFIG_SOCKET_PREFIX + server
		} else {
			options[i] = _CONFIG_SERVER_PREFIX + server
		}
	}
	return
}

func clientConfig(servers []string) string {
	return join(clientOptions(servers))
}

func poolOptions(servers []string, initSize, maxSize int) (options []string) {
	options = clientOptions(servers)
	options = append(options, _CONFIG_POOL_MIN+strconv.Itoa(initSize))
	options = append(options, _CONFIG_POOL_MAX+strconv.Itoa(maxSize))
	return
}

func poolConfig(servers []string, initSize, maxSize int) string {
	return join(poolOptions(servers, initSize, maxSize))
}
