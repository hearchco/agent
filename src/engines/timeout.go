package engines

import (
	"net"
)

func IsTimeoutError(err error) bool {
	if perr, ok := err.(net.Error); ok && perr.Timeout() {
		return true
	}
	return false
}
