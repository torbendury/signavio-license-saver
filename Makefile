-include .env
export

.PHONY: run test qa build install

run:
	@go run cmd/signavio-license-saver.go -url=$(URL) -tenant=$(TENANT) -user=$(USER) -password=$(PASSWORD) -allowlist=$(ALLOWLIST)

test:
	@go test -race -v ./...

coverage:
	@go test -v -coverprofile cover.out ./...
	@go tool cover -html cover.out -o cover.html
	@open cover.html

qa: test
	@go mod verify
	@go vet ./...
	@go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all ./...
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./...

build: test qa
	@go build -o bin/signavio-license-saver cmd/signavio-license-saver.go

install:
	@go install cmd/signavio-license-saver.go
