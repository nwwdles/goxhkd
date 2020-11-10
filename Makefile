VERSION=$(shell printf "r%s.%s" "`git rev-list --count HEAD`" "`git rev-parse --short HEAD`")
LDFLAGS=-w -s -extldflags -static -X 'main.version=$(VERSION)' -X 'main.build=$(shell date +%FT%T%z)'

.PHONY:all
all: build

.PHONY:install
install:
	go install -ldflags "$(LDFLAGS)" ./...

.PHONY:build
build:
	go build -ldflags "$(LDFLAGS)" -o . ./...
