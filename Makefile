install:
	go get ./...
	go install github.com/dmarkham/enumer@latest
	go generate ./...
	go mod tidy