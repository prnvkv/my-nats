FROM golang:1.13.1 AS builder
LABEL stage=server-intermediate

WORKDIR /go/src/my-nats
COPY . .
RUN go build -o bin/nats-req ./cmd/request-reply/pub/

FROM alpine:latest AS runner
RUN mkdir -p /lib64 && \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /usr/local/bin
COPY --from=builder /go/src/my-nats/bin/nats-req .
# CMD ["nats-pub"]
ENTRYPOINT ["nats-req"]
