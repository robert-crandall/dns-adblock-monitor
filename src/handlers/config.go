package handlers

import (
	"net"
	"sync"
)

var (
	config Config
	mu     sync.RWMutex
)

func Initialize(hosts []string, blockingIPv4 []string, blockingIPv6 []string, dnsResolver string) {
	mu.Lock()
	defer mu.Unlock()

	config.Hosts = hosts
	parseIPv4Blocks(blockingIPv4)
	parseIPv6Blocks(blockingIPv6)
	// Only initialize resolver if a custom address is provided
	if dnsResolver != "" {
		initResolver(dnsResolver)
	}
}

func parseIPv4Blocks(blockingIPv4 []string) {
	config.BlockingIPv4 = make([]net.IPNet, 0)
	for _, ipv4 := range blockingIPv4 {
		addIPv4Block(ipv4)
	}
}

func parseIPv6Blocks(blockingIPv6 []string) {
	config.BlockingIPv6 = make([]net.IPNet, 0)
	for _, ipv6 := range blockingIPv6 {
		addIPv6Block(ipv6)
	}
}
