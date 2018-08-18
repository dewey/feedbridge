IMAGE_NAME := tehwey/feedbridge

all: install

install:
	go install -v

test:
	go test ./... -v

image:
	docker build -t $(IMAGE_NAME) .

# release:
# 	git tag -a $(VERSION) -m "Release" || true
# 	git push origin $(VERSION)
# 	goreleaser --rm-dist

.PHONY: install test