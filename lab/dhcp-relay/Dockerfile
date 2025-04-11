FROM golang:1.24 AS builder
WORKDIR /work
COPY . .
RUN make

FROM alpine:3.21
WORKDIR /go-dhcp-relay
COPY --from=builder /work/bin/go-dhcp-relay /usr/bin
RUN apk add iproute2 tcpdump

ENTRYPOINT ["/usr/bin/go-dhcp-relay"]
