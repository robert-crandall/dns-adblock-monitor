# README.md

# DNS Checker

DNS Checker is a simple Go application that checks the DNS resolution for a given hostname. It responds with a 200 status code if the resolution is successful and a 500 status code if it fails.

## Project Structure

```
dns-checker
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
git clone <repository-url>
cd dns-checker
```

Then, run the application:

```bash
go run src/main.go
```

## Usage

Send a request to the server with a hostname to check its DNS resolution. The server will respond with either a 200 or a 500 status code based on the result of the DNS check.

## License

This project is licensed under the MIT License.