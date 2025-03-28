TAGS := -tags 'netgo'
BINARY := go-dhcp-relay

SHA := $(shell git rev-parse --short=8 HEAD)
GITVERSION := $(shell git describe --long --all)
# gnu date format iso-8601 is parsable with Go RFC3339
BUILDDATE := $(shell date --iso-8601=seconds)
VERSION := $(or ${VERSION},$(shell git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD))

LINKMODE := $(LINKMODE) \
		 -X 'github.com/metal-stack/v.Version=$(VERSION)' \
		 -X 'github.com/metal-stack/v.Revision=$(GITVERSION)' \
		 -X 'github.com/metal-stack/v.GitSHA1=$(SHA)' \
		 -X 'github.com/metal-stack/v.BuildDate=$(BUILDDATE)'

.PHONY: build
build:
	go build \
		$(TAGS) \
		-ldflags \
		"$(LINKMODE)" \
		-o bin/$(BINARY) \
		github.com/metal-stack/go-dhcp-relay

.PHONY: start-test-server
start-test-server: docker-build
	docker run --rm -it --network host -v $(shell pwd)/test:/etc/go-dhcp-relay go-dhcp-relay:local

.PHONY: run-test-client
run-test-client:
	docker run --rm -it --network host go-dhcp-relay:local test-client -i lo

.PHONY: lint
lint:
	golangci-lint run --build-tags client -p bugs -p unused

.PHONY: docker-build
docker-build:
	docker build -t go-dhcp-relay:local .
