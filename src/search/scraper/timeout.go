package scraper

import (
	"context"
	"net"
	"strings"
)

func IsTimeoutError(err error) bool {
	// Check if the error is a cancelled context error.
	if strings.HasSuffix(err.Error(), context.Canceled.Error()) {
		return true
	}

	// Check if the error is a timeout error.
	if perr, ok := err.(net.Error); ok && perr.Timeout() {
		return true
	}

	return false
}
