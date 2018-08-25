IMAGE_NAME := tehwey/feedbridge
VERSION_DOCKER := $(shell git describe --abbrev=0 --tags  | sed 's/^v\(.*\)/\1/')

all: install

install:
	go install -v

test:
	go test ./... -v

image:
	docker build -t $(IMAGE_NAME) .
	docker tag $(IMAGE_NAME):latest $(IMAGE_NAME):$(VERSION_DOCKER)

image-push:
	docker push $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME):$(VERSION_DOCKER)

release:
	git tag -a $(VERSION) -m "Release $(VERSION)" || true
	git push origin $(VERSION)
	goreleaser --rm-dist

.PHONY: install test