FROM golang:1.24 AS builder
WORKDIR /work
COPY . .
RUN apk add make git
RUN make build

FROM gcr.io/distroless/static-debian12
WORKDIR /go-dhcp-relay
COPY --from=builder /work/bin/go-dhcp-relay /usr/bin

ENTRYPOINT ["/usr/bin/go-dhcp-relay"]
