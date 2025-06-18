FROM golang:alpine3.21 AS builder

ARG APP

WORKDIR /build

COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o ${APP} ${APP}/cmd/main.go

FROM alpine:3.21

ARG APP

WORKDIR /app
COPY ./api/font /font
COPY --from=builder /build/${APP}/main /app/main

ENTRYPOINT ["/app/main"]
