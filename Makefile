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

test-redis-podman:
	podman run --rm --name hearchco-redis -d -p 6379:6379 redis
	go test ./src/cache/redis -count=1
	podman stop hearchco-redis

test-redis-docker:
	docker run --rm --name hearchco-redis -d -p 6379:6379 redis
	go test ./src/cache/redis -count=1
	docker stop hearchco-redis

update:
	go get -u ./...
	go mod tidy

lint:
	golangci-lint run
