PACKAGE=go-api-proxy
OUTPUT_DIR=$(shell pwd)/bin

.PHONY: all build test clean config deploy local

all: build

build: clean
	go build -v -o $(OUTPUT_DIR)/$(PACKAGE) .

test: 
	go test -v .

clean:
	rm -rf $(OUTPUT_DIR)

config:
	cat .env | xargs heroku config:set

deploy: config
	git push heroku master

local: build
	heroku local web
