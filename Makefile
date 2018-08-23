IMAGE_NAME := tehwey/feedbridge
VERSION := 0.1.3

all: install

install:
	go install -v

test:
	go test ./... -v

image:
	docker build -t $(IMAGE_NAME) .
	docker tag $(IMAGE_NAME):latest $(IMAGE_NAME):$(VERSION)

image-push:
	docker push $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME):$(VERSION)

release:
	git tag -a $(VERSION) -m "Release" || true
	git push origin $(VERSION)
	goreleaser --rm-dist

.PHONY: install test