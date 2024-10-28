PLATFORMS  ?= linux/arm64
VERSION    ?= dev
BUILD_TIME ?= N/A

chores:
	gofmt -w -s .
	go vet ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest -show verbose ./...
	go run golang.org/x/tools/cmd/deadcode@latest ./...

out_barito:
	CGO_ENABLED=1 go build \
		-ldflags "-s -w -X main.PluginVersion=$(VERSION) -X main.PluginBuildTime=$(BUILD_TIME)" \
		-buildmode=c-shared \
		-o out_barito.so ./cmd/out_barito

clean: chores
	rm -rf *.so *.h

dev:
	docker compose up --build

build:
	docker buildx build \
		--build-arg=VERSION=$(shell git describe --tags HEAD)-$(shell git rev-parse HEAD) \
		--build-arg=BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ') \
		--target prod \
		--tag barito-fluentbit-plugin:latest .

ci-build:
	docker buildx build \
		--platform=$(PLATFORMS) \
		--push \
		--build-arg=VERSION=$(shell git describe --tags HEAD)-$(shell git rev-parse HEAD) \
		--build-arg=BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ') \
		--target prod \
		--label "org.opencontainers.image.created=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')" \
		--label "org.opencontainers.image.revision=$(shell git rev-parse HEAD)" \
		--label "org.opencontainers.image.version=$(shell git describe --tags HEAD)" \
		--annotation "index:org.opencontainers.image.created=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')" \
		--annotation "index:org.opencontainers.image.revision=$(shell git rev-parse HEAD)" \
		--annotation "index:org.opencontainers.image.version=$(shell git describe --tags HEAD)" \
		--annotation "index:org.opencontainers.image.authors=Clavin June <juneardoc@gmail.com>" \
		--annotation "index:org.opencontainers.image.description=fluentbit with barito integration" \
		--annotation "index:org.opencontainers.image.source=https://github.com/clavinjune/barito-fluentbit-plugin" \
		--annotation "index:org.opencontainers.image.title=barito-fluentbit-plugin" \
		--annotation "index:org.opencontainers.image.url=https://github.com/clavinjune/barito-fluentbit-plugin" \
		--tag docker.io/juneardoc/barito-fluentbit-plugin:latest \
		--tag docker.io/juneardoc/barito-fluentbit-plugin:$(shell git describe --tags HEAD)-fluentbit-3.1.8 \
		--tag docker.io/juneardoc/barito-fluentbit-plugin:$(shell git rev-parse --short HEAD)-fluentbit-3.1.8 \
		--tag ghcr.io/clavinjune/barito-fluentbit-plugin:latest \
		--tag ghcr.io/clavinjune/barito-fluentbit-plugin:$(shell git describe --tags HEAD)-fluentbit-3.1.8 \
		--tag ghcr.io/clavinjune/barito-fluentbit-plugin:$(shell git rev-parse --short HEAD)-fluentbit-3.1.8 .

ci-build-debug:
	docker buildx build \
		--platform=$(PLATFORMS) \
		--push \
		--build-arg=VERSION=$(shell git describe --tags HEAD)-$(shell git rev-parse HEAD) \
		--build-arg=BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ') \
		--target dev \
		--label "org.opencontainers.image.created=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')" \
		--label "org.opencontainers.image.revision=$(shell git rev-parse HEAD)" \
		--label "org.opencontainers.image.version=$(shell git describe --tags HEAD)" \
		--annotation "index:org.opencontainers.image.created=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')" \
		--annotation "index:org.opencontainers.image.revision=$(shell git rev-parse HEAD)" \
		--annotation "index:org.opencontainers.image.version=$(shell git describe --tags HEAD)" \
		--annotation "index:org.opencontainers.image.authors=Clavin June <juneardoc@gmail.com>" \
		--annotation "index:org.opencontainers.image.description=fluentbit with barito integration" \
		--annotation "index:org.opencontainers.image.source=https://github.com/clavinjune/barito-fluentbit-plugin" \
		--annotation "index:org.opencontainers.image.title=barito-fluentbit-plugin" \
		--annotation "index:org.opencontainers.image.url=https://github.com/clavinjune/barito-fluentbit-plugin" \
		--tag docker.io/juneardoc/barito-fluentbit-plugin:latest-debug \
		--tag docker.io/juneardoc/barito-fluentbit-plugin:$(shell git describe --tags HEAD)-fluentbit-3.1.8-debug \
		--tag docker.io/juneardoc/barito-fluentbit-plugin:$(shell git rev-parse --short HEAD)-fluentbit-3.1.8-debug \
		--tag ghcr.io/clavinjune/barito-fluentbit-plugin:latest-debug \
		--tag ghcr.io/clavinjune/barito-fluentbit-plugin:$(shell git describe --tags HEAD)-fluentbit-3.1.8-debug \
		--tag ghcr.io/clavinjune/barito-fluentbit-plugin:$(shell git rev-parse --short HEAD)-fluentbit-3.1.8-debug .