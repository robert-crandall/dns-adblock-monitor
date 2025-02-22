# README.md

# DNS Adblock Monitor

DNS Adblock Monitor is a simple Go application that checks the DNS resolution for given hostnames. It responds with a 200 status code if the resolution is blocked, and a 500 status code if it is successful.

This is used in order to validate DNS based adblockers are working correctly.

## Project Structure

```
dns-adblock-monitor
├── src
│   ├── main.go          # Entry point of the application
│   ├── dns
│   │   └── resolver.go  # DNS resolution logic
│   └── handlers
│       └── check.go     # HTTP request handling
├── go.mod               # Module definition
├── go.sum               # Module checksums
└── README.md            # Project documentation
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
DNS_HOSTS=ads.example.com,tracker.example.com ./dns-monitor
```

## Variables

## Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| DNS_HOSTS | Yes | none | Comma-separated list of hostnames to check for DNS blocking. Check https://paileactivist.github.io/toolz/adblock.html for some example hosts. |
| BLOCKING_IPS | No | `0.0.0.0,127.0.0.1` | Comma-separated list of IP addresses that indicate successful blocking (e.g., `0.0.0.0,127.0.0.1`) |

## License

This project is licensed under the MIT License.
