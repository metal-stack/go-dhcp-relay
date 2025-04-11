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


all: test build

.PHONY: test
test:
	go test ./... -race -coverprofile=coverage.out -covermode=atomic && go tool cover -func=coverage.out

.PHONY: build
build:
	go build \
		$(TAGS) \
		-ldflags \
		"$(LINKMODE)" \
		-o bin/$(BINARY) \
		github.com/metal-stack/go-dhcp-relay/cmd/
	strip bin/go-dhcp-relay

.PHONY: lint
lint:
	golangci-lint run --build-tags client -p bugs -p unused

.PHONY: docker-build
docker-build:
	docker build -t go-dhcp-relay:latest .

.PHONY: lab-up
lab-up: docker-build
	docker build -t dhcp-relay:local -f lab/dhcp-relay/Dockerfile .
	docker build -t dhcp-server:local lab/dhcp-server
	docker build -t dhcp-client:local lab/dhcp-client
	sudo containerlab deploy -t ./lab/go-dhcp-relay.clab.yaml

.PHONY: lab-down
lab-down:
	sudo containerlab destroy -t ./lab/go-dhcp-relay.clab.yaml

.PHONY: lab-restart
lab-restart: lab-down lab-up
