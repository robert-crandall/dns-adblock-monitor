package dns

import (
    "net"
)

// ResolveDNS checks the DNS resolution for a given hostname.
func ResolveDNS(hostname string) error {
    _, err := net.LookupHost(hostname)
    return err
}