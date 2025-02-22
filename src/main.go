package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/robert-crandall/dns-adblock-monitor/src/handlers"
)

type Config struct {
	Hosts       []string
	ExpectedIPs []string
}

func loadConfig() Config {
	hostsEnv := os.Getenv("DNS_HOSTS")
	if hostsEnv == "" {
		log.Fatal("DNS_HOSTS environment variable is required")
	}

	expectedIPsEnv := os.Getenv("EXPECTED_IP_RESOLUTIONS")
	var expectedIPs []string
	if expectedIPsEnv != "" {
		expectedIPs = strings.Split(expectedIPsEnv, ",")
	}

	return Config{
		Hosts:       strings.Split(hostsEnv, ","),
		ExpectedIPs: expectedIPs,
	}
}

func main() {
	config := loadConfig()

	// Make the config available to handlers
	handlers.Initialize(config.Hosts, config.ExpectedIPs...)

	http.HandleFunc("/", handlers.CheckHandler)
	log.Printf("Starting server on :8080 with hosts: %v", config.Hosts)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
