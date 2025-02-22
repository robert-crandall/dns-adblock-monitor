FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /dns-checker ./src

# Default environment variables
ENV DNS_HOSTS=ads.google.com,adservice.google.com
ENV BLOCKING_IPV4=0.0.0.0,127.0.0.1
ENV BLOCKING_IPV6=::,::1
ENV DNS_RESOLVER=1.1.1.1:53

EXPOSE 8080

CMD ["/dns-checker"]
