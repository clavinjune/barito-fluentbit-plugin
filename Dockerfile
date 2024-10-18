FROM golang:1.23.2 AS builder
ENV CGO_ENABLED=1
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download -x
COPY . /app/
RUN make out_barito_batch_k8s

FROM ghcr.io/fluent/fluent-operator/fluent-bit:3.1.8-debug
COPY --from=builder /app/out_barito_batch_k8s.so /fluent-bit/plugins/out_barito_batch_k8s.so
ENTRYPOINT ["/fluent-bit/bin/fluent-bit", "-e", "/fluent-bit/plugins/out_barito_batch_k8s.so"]