# brzaguza
BrzaGuza is a private distributed metasearch engine with crowdsourced cache and speed as the primary goal.

## How to run
- `make install` to install required dependencies
- `go test` to run automated unit tests
- `go run ./src` to start the webserver (add `-v` for debug logs or `-vv` for trace logs)

### CLI mode
If you want to do a single search from the cli (add `-v` for debug logs or `-vv` for trace logs):
- `go run ./src --cli --query "my first search"`