.PHONY: compile build build_all fmt test vet bootstrap
	
SOURCE_FOLDER := .

BINARY_PATH ?= ./bin/forward2docker

GOARCH ?= amd64

ifdef GOOS
BINARY_PATH :=$(BINARY_PATH).$(GOOS)-$(GOARCH)
endif

export GO15VENDOREXPERIMENT=1

default: build
	
build_all: vet fmt
	for GOOS in darwin linux windows; do \
		$(MAKE) compile GOOS=$$GOOS GOARCH=amd64 ; \
	done; \
	$(MAKE) compile GOOS=windows GOARCH=386

compile:
	CGO_ENABLED=0 go build -i -v -ldflags '-s' -o $(BINARY_PATH) $(SOURCE_FOLDER)/

build: vet fmt compile
	
fmt:
	go fmt $(glide novendor)

vet:
	go vet $(glide novendor)

lint:
	go list $(SOURCE_FOLDER)/... | grep -v /vendor/ | xargs -L1 golint
	
bootstrap:
	go get github.com/Masterminds/glide
	glide install --use-gopath --cache-gopath