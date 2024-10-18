chores:
	gofmt -w -s .
	go vet ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest -show verbose ./...
	go run golang.org/x/tools/cmd/deadcode@latest ./...

out_barito_batch_k8s: clean
	CGO_ENABLED=1 go build \
		-ldflags "-s -w" \
		-buildmode=c-shared \
		-o out_barito_batch_k8s.so ./cmd/out_barito_batch_k8s

clean: chores
	rm -rf *.so *.h

dev:
	docker compose up --build