TAG ?= dev
VERSION ?= $(shell git describe --tags --always --dirty='-dev')
CONTAINER_IMAGE ?= swoga/mikrotik-exporter

.PHONY: build
build: build-go build-container

.PHONY: build-go
build-go:
	env CGO_ENABLED=0 go build -v -ldflags="-X main.version=$(VERSION)" ./cmd/mikrotik-exporter

.PHONY: build-container
build-container:
	docker buildx build --load -t $(CONTAINER_IMAGE):$(TAG) .

.PHONY: push-container
push-container:
	docker tag $(CONTAINER_IMAGE):$(TAG) quay.io/$(CONTAINER_IMAGE):$(TAG)
	docker tag $(CONTAINER_IMAGE):$(TAG) ghcr.io/$(CONTAINER_IMAGE):$(TAG)
	docker push $(CONTAINER_IMAGE):$(TAG)
	docker push quay.io/$(CONTAINER_IMAGE):$(TAG)
	docker push ghcr.io/$(CONTAINER_IMAGE):$(TAG)

.PHONY: act
act:
	act -r --artifact-server-path /tmp/act-artifacts