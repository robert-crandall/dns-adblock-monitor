package main

import (
	"net/http"

	"github.com/robert-crandall/dns-adblock-monitor/src/handlers"
)

func main() {
	http.HandleFunc("/check", handlers.CheckHandler)
	http.ListenAndServe(":8080", nil)
}
