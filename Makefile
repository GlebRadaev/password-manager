.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	golangci-lint run ./...

.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

.PHONY: coverage
coverage:
	@echo "Running coverage..." 
	go test ./... -v -parallel=32 -coverprofile=coverage.txt -covermode=atomic && go tool cover -html=coverage.txt && rm -rf coverage.txt

.PHONY: build
build:
	@echo "Building binary..."
	cd cmd/gophermart && go build .

.PHONY: run
run:
	@echo "Running binary..."
	cd cmd/gophermart && go run .

.PHONY: swagger
swagger: 
	@echo "Generating swagger docs..."
	swag fmt
	swag init -g cmd/gophermart/main.go 
