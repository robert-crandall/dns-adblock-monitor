package handlers

import (
	"context"
	"net"
	"time"
)

var (
	resolver *net.Resolver
)

func initResolver(dnsResolver string) {
	if dnsResolver != "" {
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Second * 10,
				}
				return d.DialContext(ctx, "udp", dnsResolver)
			},
		}
	} else {
		resolver = net.DefaultResolver
	}
}
