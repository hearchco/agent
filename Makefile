run:
	air --pretty
run-cli:
	go run ./src --cli --pretty

debug:
	air -- -v --pretty
debug-cli:
	go run ./srv -v --cli --pretty

trace:
	air -- -vv --pretty
trace-cli:
	go run ./src -vv --cli --pretty

install:
	go get ./...
	go install github.com/dmarkham/enumer@latest
	go generate ./...
	go install github.com/cosmtrek/air@latest

compile:
	CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath ./src/...
compile-linux:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -trimpath -o bin/hearchco ./src
compile-macos:
	CGO_ENABLED=0 GOOS=darwin go build -ldflags "-s -w" -trimpath -o bin/hearchco ./src
compile-windows:
	CGO_ENABLED=0 GOOS=windows go build -ldflags "-s -w" -trimpath -o bin/hearchco.exe ./src

test:
	sh ./scripts/test.sh
test-engines:
	sh ./scripts/test-engines.sh
test-all: test test-engines

update:
	go get -u ./...
	go mod tidy

lint:
	golangci-lint run
