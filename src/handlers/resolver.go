package handlers

import (
	"context"
	"net"
	"time"
)

type Resolver interface {
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
}

var resolver Resolver = net.DefaultResolver

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
