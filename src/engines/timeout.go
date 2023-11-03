package engines

import (
	"errors"
	"net"
)

func IsTimeoutError(err error) bool {
	var netError *net.Error
	is := errors.As(err, netError)
	return is && (*netError).Timeout()
}
