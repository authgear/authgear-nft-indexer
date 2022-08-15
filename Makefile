GIT_NAME ?= $(shell git describe --exact-match 2>/dev/null)
GIT_HASH ?= git-$(shell git rev-parse --short=12 HEAD)

LDFLAGS ?= "-X github.com/authgear/authgear-nft-indexer/pkg/version.Version=${GIT_HASH}"

.PHONY: start-worker
start-worker:
	go run ./cmd/indexer start

.PHONY: setup
setup: vendor
	cp authgear-nft-indexer.yaml.example authgear-nft-indexer.yaml

.PHONY: vendor
vendor:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.46.2
	go mod download
	go install github.com/google/wire/cmd/wire@v0.5.0

.PHONY: generate
generate:
	go generate ./pkg/... ./cmd/...

.PHONY: test
test:
	go test ./pkg/... -timeout 1m30s

.PHONY: lint
lint:
	golangci-lint run ./cmd/... ./pkg/...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build:
	go build -o $(BIN_NAME) -tags "$(GO_BUILD_TAGS)" -ldflags ${LDFLAGS} ./cmd/$(TARGET)

.PHONY: check-tidy
check-tidy:
	$(MAKE) fmt
	$(MAKE) generate
	go mod tidy
	git status --porcelain | grep '.*'; test $$? -eq 1

.PHONY: binary
binary:
	rm -rf ./dist
	mkdir ./dist
	$(MAKE) build TARGET=indexer BIN_NAME=./dist/authgear-nft-indexer-"$(shell go env GOOS)"-"$(shell go env GOARCH)"-${GIT_HASH}
	$(MAKE) build TARGET=server BIN_NAME=./dist/authgear-nft-server-"$(shell go env GOOS)"-"$(shell go env GOARCH)"-${GIT_HASH}


.PHONY: build-image
build-image:
	# Add --pull so that we are using the latest base image.
	docker build --pull --file ./cmd/$(TARGET)/Dockerfile --tag $(IMAGE_NAME) --build-arg GIT_HASH=$(GIT_HASH) .

.PHONY: tag-image
tag-image: DOCKER_IMAGE = quay.io/theauthgear/$(IMAGE_NAME)
tag-image:
	docker tag $(IMAGE_NAME) $(DOCKER_IMAGE):latest
	docker tag $(IMAGE_NAME) $(DOCKER_IMAGE):$(GIT_HASH)
	if [ ! -z $(GIT_NAME) ]; then docker tag $(IMAGE_NAME) $(DOCKER_IMAGE):$(GIT_NAME); fi

.PHONY: push-image
push-image: DOCKER_IMAGE = quay.io/theauthgear/$(IMAGE_NAME)
push-image:
	docker push $(DOCKER_IMAGE):latest
	docker push $(DOCKER_IMAGE):$(GIT_HASH)
	if [ ! -z $(GIT_NAME) ]; then docker push $(DOCKER_IMAGE):$(GIT_NAME); fi
