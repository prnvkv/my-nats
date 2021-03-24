FROM golang:1.14.6-alpine3.12 AS builder

RUN apk add --update git
RUN git clone https://github.com/nats-io/nats.go.git
RUN cd nats.go && go build -o /nats-bench examples/nats-bench/main.go

FROM alpine:3.11

RUN apk add --no-cache --update ca-certificates
USER nobody
COPY --from=builder /nats-bench /bin/nats-bench

ENTRYPOINT ["/bin/nats-bench"]

