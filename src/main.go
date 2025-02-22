package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/robert-crandall/dns-adblock-monitor/src/handlers"
)

type Config struct {
	Hosts        []string
	BlockingIPv4 []string
	BlockingIPv6 []string
	Resolver     string
}

func loadConfig() Config {
	hostsEnv := os.Getenv("DNS_HOSTS")
	if hostsEnv == "" {
		log.Fatal("DNS_HOSTS environment variable is required")
	}

	blockingIPv4Env := os.Getenv("BLOCKING_IPV4")
	var blockingIPv4 []string
	if blockingIPv4Env != "" {
		blockingIPv4 = strings.Split(blockingIPv4Env, ",")
	} else {
		blockingIPv4 = []string{"0.0.0.0", "127.0.0.1"}
	}

	blockingIPv6Env := os.Getenv("BLOCKING_IPV6")
	var blockingIPv6 []string
	if blockingIPv6Env != "" {
		blockingIPv6 = strings.Split(blockingIPv6Env, ",")
	} else {
		blockingIPv6 = []string{"::", "::1"}
	}

	resolverEnv := os.Getenv("DNS_RESOLVER")

	return Config{
		Hosts:        strings.Split(hostsEnv, ","),
		BlockingIPv4: blockingIPv4,
		BlockingIPv6: blockingIPv6,
		Resolver:     resolverEnv,
	}
}

func main() {
	config := loadConfig()

	// Make the config available to handlers
	handlers.Initialize(config.Hosts, config.BlockingIPv4, config.BlockingIPv6, config.Resolver)

	http.HandleFunc("/", handlers.CheckHandler)
	log.Printf("Starting server on :8080 with hosts: %v", config.Hosts)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
