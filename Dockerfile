FROM golang:1.24 AS builder
WORKDIR /work
COPY . .
RUN make

FROM gcr.io/distroless/static-debian12
WORKDIR /go-dhcp-relay
COPY --from=builder /work/bin/go-dhcp-relay /usr/bin

ENTRYPOINT ["/usr/bin/go-dhcp-relay"]
