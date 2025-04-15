Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")
SERVICES=$(shell ls -1 api | grep \.proto | sed s/\.proto//)
PROTO_FILE := ./
PROTOC_VER = 3.12.4
OS = linux
ifeq ($(shell uname -s), Darwin)
    OS = osx
endif

BINARY_NAME = pm
VERSION ?= $(shell git describe --tags 2>/dev/null || echo "v0.1.0")
CLIENT_DIR = client
BUILD_DIR = bin
LDFLAGS = -ldflags="-s -w -X main.version=$(VERSION)"

PLATFORMS = \
	darwin/amd64 \
	darwin/arm64 \
	linux/amd64 \
	windows/amd64

build-all:
	mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		OUTPUT=$(BUILD_DIR)/$(BINARY_NAME)-$$GOOS-$$GOARCH; \
		if [ $$GOOS = "windows" ]; then OUTPUT=$$OUTPUT.exe; fi; \
		echo "Building $$GOOS/$$GOARCH -> $$OUTPUT"; \
		GOOS=$$GOOS GOARCH=$$GOARCH go build $(LDFLAGS) -o $$OUTPUT ./$(CLIENT_DIR); \
	done
    

.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	golangci-lint run --config ./.golangci.yml --timeout=5m ./...

.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

.PHONY: coverage
coverage:
	@echo "Running coverage..." 
	go test ./... -v -parallel=32 -coverprofile=coverage.txt -covermode=atomic && go tool cover -html=coverage.txt && rm -rf coverage.txt

run: 
	$(info $(M) running)
	docker build -t manager:latest -f Dockerfile.local .
	docker compose up -d

.PHONY: bin
bin: $(info $(M) install bin)
	@GOBIN=$(CURDIR)/bin go install -mod=mod \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.25.1 \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.25.1; \
    GOBIN=$(CURDIR)/bin go install -mod=mod \
        google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.2; \
    GOBIN=$(CURDIR)/bin go install -mod=mod \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1; \
    GOBIN=$(CURDIR)/bin go install -mod=mod \
        github.com/envoyproxy/protoc-gen-validate@v1.1.0; \
    curl -Ls -o $(CURDIR)/bin/protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VER}/protoc-${PROTOC_VER}-${OS}-x86_64.zip; \
    mkdir $(CURDIR)/bin/.protoc; \
    unzip -q $(CURDIR)/bin/protoc.zip -d $(CURDIR)/bin/.protoc; \
    mv $(CURDIR)/bin/.protoc/bin/protoc $(CURDIR)/bin; \
    rm -rf $(CURDIR)/bin/.protoc; rm -rf $(CURDIR)/bin/protoc.zip;

.PHONY: gen
gen: $(info $(M) protoc gen)
	$(Q) for srv in $(SERVICES); do \
		echo "Generate $$srv" && \
		mkdir -p $(CURDIR)/cmd/$$srv && \
		mkdir -p $(CURDIR)/internal/$$srv && \
		mkdir -p $(CURDIR)/pkg/$$srv && \
		mkdir -p $(CURDIR)/docs/swagger && \
        $(CURDIR)/bin/protoc \
            --plugin=protoc-gen-grpc-gateway=$(CURDIR)/bin/protoc-gen-grpc-gateway \
            --plugin=protoc-gen-openapiv2=$(CURDIR)/bin/protoc-gen-openapiv2 \
            --plugin=protoc-gen-go-grpc=$(CURDIR)/bin/protoc-gen-go-grpc \
            --plugin=protoc-gen-validate=$(CURDIR)/bin/protoc-gen-validate \
            --plugin=protoc-gen-go=$(CURDIR)/bin/protoc-gen-go \
            -I$(CURDIR)/api:$(CURDIR)/vendor.pb \
            --go_out=$(CURDIR)/pkg \
            --validate_out=lang=go:$(CURDIR)/pkg \
            --go-grpc_out=$(CURDIR)/pkg \
			--experimental_allow_proto3_optional \
            --grpc-gateway_out=$(CURDIR)/pkg \
            --grpc-gateway_opt=logtostderr=true \
            --openapiv2_out=$(CURDIR)/docs/swagger \
            --openapiv2_opt=logtostderr=true \
            --openapiv2_opt=use_go_templates=true \
            $(CURDIR)/api/$$srv.proto ; \
	done


.PHONY: doc-check
doc-check:
	@echo "▶ Running documentation checks..."
	@if ! command -v revive >/dev/null 2>&1; then \
		echo "Error: revive not found. Install with: go install github.com/mgechev/revive@latest"; \
		exit 1; \
	fi
	@revive -config .revive.toml -formatter friendly \
		-exclude vendor/... \
		-exclude vendor.pb/... \
		./...

.PHONY: buffmt
buffmt:
	docker run --rm \
		-v $(shell pwd):/app \
		-w /app \
		bufbuild/buf \
		format -w $(PROTO_FILE)