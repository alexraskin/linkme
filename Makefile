.PHONY: build run clean docker-build docker-run

BUILD_TIME := $(shell date '+%b %d, %Y')
DIST := ./dist
BINARY := $(DIST)/linkme
IMAGE := linkme

build:
	@mkdir -p $(DIST)
	go build -ldflags "-X 'main.buildTime=$(BUILD_TIME)'" -o $(BINARY) .

run: build
	$(BINARY)

clean:
	rm -rf $(DIST)

docker-build:
	docker build --build-arg BUILD_TIME="$(BUILD_TIME)" -t $(IMAGE) .

docker-run: docker-build
	docker run -p 8080:8080 $(IMAGE)
