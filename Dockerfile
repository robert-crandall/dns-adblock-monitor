FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /dns-checker ./src

ENV DNS_HOSTS=example.com,google.com
ENV EXPECTED_IP_RESOLUTIONS=0.0.0.0,127.0.0.1


EXPOSE 8080

CMD ["/dns-checker"]
