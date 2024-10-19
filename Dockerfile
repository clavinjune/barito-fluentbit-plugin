# can't use alpine issue with musl
FROM golang:1.23.2 AS builder

ARG VERSION
ARG BUILD_TIME
ENV CGO_ENABLED=1

WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download -x
COPY . /app/
RUN make out_barito_batch_k8s

FROM ghcr.io/fluent/fluent-operator/fluent-bit:3.1.8-debug AS dev
COPY --from=builder /app/out_barito_batch_k8s.so /fluent-bit/plugins/out_barito_batch_k8s.so
ENTRYPOINT ["/fluent-bit/bin/fluent-bit", "-e", "/fluent-bit/plugins/out_barito_batch_k8s.so"]

FROM ghcr.io/fluent/fluent-operator/fluent-bit:3.1.8 AS prod
LABEL org.opencontainers.image.authors="Clavin June <juneardoc@gmail.com>"
LABEL org.opencontainers.image.description="fluentbit with barito integration"
LABEL org.opencontainers.image.source="https://github.com/clavinjune/barito-fluentbit-plugin"
LABEL org.opencontainers.image.title="barito-fluentbit-plugin"
LABEL org.opencontainers.image.url="https://github.com/clavinjune/barito-fluentbit-plugin"

COPY --from=builder /app/out_barito_batch_k8s.so /fluent-bit/plugins/out_barito_batch_k8s.so
ENTRYPOINT ["/fluent-bit/bin/fluent-bit", "-e", "/fluent-bit/plugins/out_barito_batch_k8s.so"]