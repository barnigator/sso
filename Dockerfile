FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/bin/sso \
    ./cmd/sso


FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/sso /app/sso
COPY config /app/config

ENV CONFIG_PATH=/app/config/local.yaml

EXPOSE 44044

ENTRYPOINT ["/app/sso"]