BINARY_NAME=ghex
VERSION=$(shell cat VERSION 2>/dev/null || echo "1.0.0")
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

.PHONY: all build build-all clean test run

all: build

build:
	go build $(LDFLAGS) -o build/$(BINARY_NAME) ./cmd/ghex

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o build/$(BINARY_NAME)-linux-amd64 ./cmd/ghex

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o build/$(BINARY_NAME)-linux-arm64 ./cmd/ghex

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o build/$(BINARY_NAME)-darwin-amd64 ./cmd/ghex

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o build/$(BINARY_NAME)-darwin-arm64 ./cmd/ghex

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o build/$(BINARY_NAME)-windows-amd64.exe ./cmd/ghex

build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64

clean:
	rm -rf build/

test:
	go test -v ./...

test-prop:
	go test -v -run "Prop" ./...

run:
	go run ./cmd/ghex

deps:
	go mod download
	go mod tidy

fmt:
	go fmt ./...

lint:
	golangci-lint run

checksums:
	cd build && sha256sum $(BINARY_NAME)-* > checksums.txt

install: build
	cp build/$(BINARY_NAME) /usr/local/bin/

uninstall:
	rm -f /usr/local/bin/$(BINARY_NAME)

# Release targets
release-dry:
	goreleaser release --snapshot --clean

release:
	goreleaser release --clean

# Create archives for release
archives: build-all
	cd build && tar -czvf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 ../README.md ../LICENSE
	cd build && tar -czvf $(BINARY_NAME)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64 ../README.md ../LICENSE
	cd build && tar -czvf $(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 ../README.md ../LICENSE
	cd build && tar -czvf $(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64 ../README.md ../LICENSE
	cd build && zip $(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe ../README.md ../LICENSE

# Version management
version:
	@echo $(VERSION)

bump-patch:
	@echo $$(echo $(VERSION) | awk -F. '{print $$1"."$$2"."$$3+1}') > VERSION
	@echo "Version bumped to $$(cat VERSION)"

bump-minor:
	@echo $$(echo $(VERSION) | awk -F. '{print $$1"."$$2+1".0"}') > VERSION
	@echo "Version bumped to $$(cat VERSION)"

bump-major:
	@echo $$(echo $(VERSION) | awk -F. '{print $$1+1".0.0"}') > VERSION
	@echo "Version bumped to $$(cat VERSION)"

tag:
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)

# Help
help:
	@echo "GHEX Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build          Build for current platform"
	@echo "  make build-all      Build for all platforms"
	@echo "  make test           Run tests"
	@echo "  make install        Install to /usr/local/bin"
	@echo "  make clean          Clean build artifacts"
	@echo ""
	@echo "Release:"
	@echo "  make archives       Create release archives"
	@echo "  make checksums      Generate checksums"
	@echo "  make release-dry    Test release with goreleaser"
	@echo "  make release        Create release with goreleaser"
	@echo ""
	@echo "Version:"
	@echo "  make version        Show current version"
	@echo "  make bump-patch     Bump patch version (1.0.0 -> 1.0.1)"
	@echo "  make bump-minor     Bump minor version (1.0.0 -> 1.1.0)"
	@echo "  make bump-major     Bump major version (1.0.0 -> 2.0.0)"
	@echo "  make tag            Create and push git tag"
