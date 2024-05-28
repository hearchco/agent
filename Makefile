run:
	air -- --pretty
run-cli:
	go run ./src --cli --pretty

debug:
	air -- --pretty -v
debug-cli:
	go run ./srv --cli --pretty -v

trace:
	air -- --pretty -vv
trace-cli:
	go run ./src --cli --pretty -vv

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
test-redis:
	sh ./scripts/test-redis.sh
test-redis-podman:
	sh ./scripts/test-redis-podman.sh
test-redis-docker:
	sh ./scripts/test-redis-docker.sh
test-all: test test-redis test-engines
test-all-podman: test test-redis-podman test-engines
test-all-docker: test test-redis-docker test-engines

update:
	go get -u ./...
	go mod tidy

lint:
	golangci-lint run
