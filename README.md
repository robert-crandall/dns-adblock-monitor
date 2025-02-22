# DNS Adblock Monitor

DNS Adblock Monitor is a simple Go application that checks the DNS resolution for given hostnames. It responds with a 200 status code if the resolution is blocked, and a 500 status code if it is successful.

This is used in order to validate DNS based adblockers are working correctly.

The motivation for this project came from needing to monitor my guest wifi network for adblocking.

## Project Structure

```
dns-adblock-monitor
├── src
│   ├── main.go         # Entry point of the application
│   └── handlers
│       ├── types.go    # Type definitions
│       ├── resolver.go # DNS resolver configuration
│       ├── check.go    # HTTP request handling
│       └── check_test.go # Unit tests
├── Dockerfile          # Container build configuration
├── go.mod             # Module definition
├── go.sum             # Module checksums
└── README.md          # Project documentation
```

## Getting Started

To run the application, ensure you have Go installed on your machine. Clone the repository and navigate to the project directory:

```bash
git clone github.com/robert-crandall/dns-adblock-monitor
cd dns-adblock-monitor
```

Then, run the application:

```bash
go build -o dns-monitor src/main.go
DNS_HOSTS=ads.example.com,tracker.example.com DNS_RESOLVER=192.168.1.1:53 ./dns-monitor
```

### Docker

You can also run this in docker.

#### Using Docker CLI

```bash
docker run -p 8080:8080 \
  -e DNS_HOSTS=ads.google.com,adservice.google.com \
  -e DNS_RESOLVER=1.1.1.1:53 \
  ghcr.io/robert-crandall/dns-adblock-monitor:latest
```

#### Using Docker Compose


```yaml
services:
  dns-monitor:
    image: ghcr.io/robert-crandall/dns-adblock-monitor:latest
    ports:
      - "8080:8080"
    environment:
      - DNS_HOSTS=ads.google.com,adservice.google.com
      - DNS_RESOLVER=1.1.1.1:53
      # Optional environment variables with defaults
      - BLOCKING_IPV4=0.0.0.0/8,127.0.0.0/8
      - BLOCKING_IPV6=::/128,::1/128,fc00::/7
    restart: unless-stopped
```

## Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| DNS_HOSTS | Yes | none | Comma-separated list of hostnames to check for DNS blocking. Check https://paileactivist.github.io/toolz/adblock.html for some example hosts. |
| BLOCKING_IPV4 | No | `0.0.0.0,127.0.0.1` | Comma-separated list of IPv4 addresses or CIDR blocks (e.g., `0.0.0.0,127.0.0.0/8`) |
| BLOCKING_IPV6 | No | `::,::1` | Comma-separated list of IPv6 addresses or CIDR blocks (e.g., `::1/128,fc00::/7`) |
| DNS_RESOLVER | No | System Default | DNS resolver to use (e.g., `1.1.1.1:53`, `8.8.8.8:53`) |

## License

This project is licensed under the MIT License.
