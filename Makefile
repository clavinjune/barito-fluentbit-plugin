chores:
	gofmt -w -s .
	go vet ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest -show verbose ./...
	go run golang.org/x/tools/cmd/deadcode@latest ./...

out_barito_batch_k8s:
	CGO_ENABLED=1 go build \
		-ldflags "-s -w -X main.PluginVersion=$(VERSION) -X main.PluginBuildTime=$(BUILD_TIME)" \
		-buildmode=c-shared \
		-o out_barito_batch_k8s.so ./cmd/out_barito_batch_k8s

clean: chores
	rm -rf *.so *.h

dev:
	docker compose up --build

build:
	docker build \
		--build-arg=VERSION=$(shell git describe --tags HEAD)-$(shell git rev-parse HEAD) \
		--build-arg=BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ') \
		--target prod \
		--tag ghcr.io/clavinjune/barito-fluentbit-plugin:latest \
		--tag ghcr.io/clavinjune/barito-fluentbit-plugin:$(shell git describe --tags HEAD) \
		--tag ghcr.io/clavinjune/barito-fluentbit-plugin:$(shell git rev-parse --short HEAD) .

build-and-push: build
	docker push ghcr.io/clavinjune/barito-fluentbit-plugin:latest
	docker push ghcr.io/clavinjune/barito-fluentbit-plugin:$(shell git describe --tags HEAD)
	docker push ghcr.io/clavinjune/barito-fluentbit-plugin:$(shell git rev-parse --short HEAD)