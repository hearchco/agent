run:
	go run ./src

debug:
	go run ./src -v

trace:
	go run ./src -vv

setup:
	go get ./...
	go install github.com/dmarkham/enumer@latest
	go generate ./...

build:
	go build ./...

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

update:
	go get -u ./...
	go mod tidy

lint:
	golangci-lint run
