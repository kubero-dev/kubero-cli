.PHONY: all build test install clean docs

all: build

build:
	go build -o bin/kubero cmd/kubero/main.go

test:
	go test ./...

install:
	./scripts/install.sh

clean:
	rm -rf bin/

docs:
	godoc -http=:6060