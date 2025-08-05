TARGET_DIRECTORY := $(shell pwd)

all: build/thape

build/thape:
	@mkdir -p $(TARGET_DIRECTORY)/build
	cd cmd/thape && go build -o $(TARGET_DIRECTORY)/build

.PHONY: clean dev

clean:
	rm -rf $(TARGET_DIRECTORY)/build

dev:
	@go install github.com/air-verse/air@latest
	cd $(TARGET_DIRECTORY) && air
