all: fmt build

fmt:
	goimports -w .

test:
	go test -v ./...

build:
	go build -o output/