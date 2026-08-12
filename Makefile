TARGET := formdrain

OPENAPI_SPEC := api/openapi.yaml
OAPI_CONFIG := routes/apiv1/api/oapi-codegen.yaml
OAPI_GENERATE := routes/apiv1/api/generate.go
GENERATED := routes/apiv1/api/server.gen.go

SOURCES := $(shell find . -name '*.go' -not -name '*_test.go')

.PHONY: all clean build format test check

all: clean build

clean:
	rm -f '$(TARGET)' '$(GENERATED)'
	go clean -testcache

$(GENERATED): $(OPENAPI_SPEC) $(OAPI_CONFIG) $(OAPI_GENERATE)
	go generate ./routes/apiv1/api/...

generate: $(GENERATED)

$(TARGET):  $(GENERATED) $(SOURCES)
	go build -o $(TARGET) .

build: $(TARGET)

format:
	go fmt ./...
	go mod tidy

test: $(GENERATED)
	go test -timeout 30s ./... | sed '/^?/d'

check: format test
