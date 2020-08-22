BUILD_FILES = $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}}\
{{end}}' ./...)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOLINT=golint

# Build parameters
BINARY_NAME=glassfactory
BINARY_OUT=bin/$(BINARY_NAME)

.PHONY: all
all: build

.PHONY: build
build: $(BINARY_OUT)

$(BINARY_OUT): $(BUILD_FILES)
	$(GOBUILD) -trimpath -o "$@" ./cmd/glassfactory

.PHONY: install
install: $(BINARY_OUT)
	cp $^ /usr/local/bin/$(BINARY_NAME)

.PHONY: uninstall
uninstall:
	rm /usr/local/bin/$(BINARY_NAME)

.PHONY:
lint:
	$(GOVET) ./...
	$(GOLINT) ./...

.PHONY: format
format:
	$(GOFMT)

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_OUT)


