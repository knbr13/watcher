

ARGS ?= ''

all: clean build run

run:
	./bin/watcher $(ARGS)

build:
	go build -o bin/watcher

clean:
	rm -rf bin/*

