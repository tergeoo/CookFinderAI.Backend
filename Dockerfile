FROM golang:1.24.1-alpine AS builder

WORKDIR /build

COPY go.mod ./
COPY go.sum ./   # если есть go.sum
RUN go mod download

COPY . .

# Собираем из правильной директории
RUN go build -o /build/app ./cmd/api/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /build/app .

EXPOSE 8080
ENTRYPOINT ["/app/app"]
