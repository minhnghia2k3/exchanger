# Development
FROM golang:1.23.2-alpine AS dev
WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]

# Build
FROM golang:1.23.2-alpine AS build

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/main ./cmd/api

# Production
FROM debian:bookworm-slim AS prod
WORKDIR /

COPY --from=build /app /app

# Install TLS certificatres
RUN apt-get update && apt-get install -y ca-certificates

CMD "/app"

