# ---------- Build Stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker ./cmd/worker
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o adminctl ./cmd/adminctl

# ---------- Runtime Stage ----------
FROM alpine:latest

WORKDIR /root/
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/api .
COPY --from=builder /app/worker .
COPY --from=builder /app/adminctl .
