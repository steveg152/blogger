all: build

build:
	@echo "Building..."
	
	@go build -o ./tmp/main  cmd/api/main.go