.PHONY: setup
setup: vendor
	cp authgear-nft-indexer.yaml.example authgear-nft-indexer.yaml.yaml

.PHONY: vendor
vendor:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.46.2
	go mod download

.PHONY: test
test:
	go test ./pkg/... -timeout 1m30s