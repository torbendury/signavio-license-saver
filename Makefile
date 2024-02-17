include .env
export

run:
	go run cmd/signavio-license-saver.go -url=$(URL) -tenant=$(TENANT) -user=$(USER) -password=$(PASSWORD) -allowlist=$(ALLOWLIST)

test:
	go test -race -cover -v ./...
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

build:
	go build -o bin/signavio-license-saver cmd/signavio-license-saver.go

install:
	go install cmd/signavio-license-saver.go
