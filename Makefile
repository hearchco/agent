install:
	go get ./...
	go install github.com/dmarkham/enumer@latest
	go generate ./...

build:
	go build -o brzaguza-bin ./src

test:
	go test ./... -count=1

update:
	go get -u ./...
	go mod tidy