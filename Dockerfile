# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Required tools
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Install swag CLI and generate docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN $(go env GOPATH)/bin/swag init

RUN go build -o inventory-app .

# Stage 2: Lightweight container to run the app
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy app binary and docs
COPY --from=builder /app/inventory-app .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./inventory-app"]
