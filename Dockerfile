ARG GO_VERSION=1.23.1
ARG GO_BASE_IMAGE=alpine

FROM golang:${GO_VERSION}-${GO_BASE_IMAGE} AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o bin/drop-auth ./cmd/auth
RUN go build -o bin/migrator ./cmd/migrator

FROM alpine:latest
WORKDIR /app
COPY --from=builder ./app/bin ./bin
COPY --from=builder ./app/internal/data ./data
COPY --from=builder ./app/tls ./tls