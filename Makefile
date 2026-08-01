TARGET := formdrain
SOURCES := $(shell find . -name '*.go' -not -name '*_test.go')

.PHONY: all clean build format test

all: clean build

clean:
	rm -f $(TARGET)
	go clean -testcache

build: $(TARGET)

$(TARGET): $(SOURCES)
	go build -o $(TARGET) .

format:
	go fmt ./...
	go mod tidy

test: format
	go test -timeout 30s ./... | sed '/^?/d'
