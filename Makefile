ARGS ?= ''
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X main.version=$(VERSION)

all: clean build run

run:
	./bin/watcher $(ARGS)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/watcher

clean:
	rm -rf bin/*

