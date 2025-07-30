FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN $(go env GOPATH)/bin/swag init --parseDependency --parseInternal

RUN go build -o inventory-app .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/inventory-app .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./inventory-app"]
