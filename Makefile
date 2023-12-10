setup:
	go get ./...
	go install github.com/dmarkham/enumer@latest
	go generate ./...

build:
	go build ./...

test:
	go test ./... -count=1

update:
	go get -u ./...
	go mod tidy

lint:
	golangci-lint run