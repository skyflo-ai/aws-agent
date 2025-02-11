# Dockerfile
FROM golang:1.23.6 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o crawler ./cmd/crawler/main.go

FROM amazon/aws-cli:latest
RUN yum install -y ca-certificates

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
COPY --from=builder /app/crawler /crawler
EXPOSE 8181
ENTRYPOINT ["/entrypoint.sh"]
