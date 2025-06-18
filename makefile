.PHONY: test
BINDIR = $(shell pwd)/bin
HOST=irm@49.13.141.169
PORT=22

dev-install-tooling:
	GOBIN=$(BINDIR) go install github.com/golang/mock/mockgen@latest
	GOBIN=$(BINDIR) go install github.com/swaggo/swag/cmd/swag@latest
	GOBIN=$(BINDIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	GOBIN=$(BINDIR) go install github.com/divan/expvarmon@latest
	GOBIN=$(BINDIR) go install github.com/rakyll/hey@latest
	GOBIN=$(BINDIR) go install github.com/pressly/goose/v3/cmd/goose@latest

dev-up:
	docker compose -f docker/docker-compose.yml up -d

dev-down:
	docker compose -f docker/docker-compose.yml down

gen-api-swag:
	swag init -g cmd/api/main.go -o docs

gen-proto:
	@PATH=$(BINDIR):$(PATH) protoc --proto_path=shared/pkg/proto --go_out=shared/pkg/proto --go-grpc_out=shared/pkg/proto shared/pkg/proto/*.proto
	@git add shared/pkg/proto/pb/*.pb.go

lint:
	@$(BINDIR)/golangci-lint run ./...

gen-mocks:
	GOBIN=$(BINDIR) go generate ./...
	git add *.go

bench:
	go test -bench=. -benchmem -cpu=4 -run=^# ./...

build-api:
	make _build SERVICE=api

_build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o bin/$(SERVICE) cmd/$(SERVICE)/main.go