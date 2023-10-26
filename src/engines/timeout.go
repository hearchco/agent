package engines

import "net"

func IsTimeoutError(err error) bool {
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}
