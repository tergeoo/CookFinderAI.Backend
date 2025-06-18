FROM golang:1.24.1-alpine AS builder

ARG APP=app

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Укажи путь к main.go здесь
RUN go build -o /build/${APP} ./cmd/${APP}/main.go

FROM alpine:3.21

ARG APP=app

WORKDIR /app

COPY --from=builder /build/${APP} /app/main

EXPOSE 8080
ENTRYPOINT ["/app/main"]
