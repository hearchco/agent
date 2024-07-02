package middlewares

import (
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/go-chi/chi/v5/middleware"
)

func compress(lvl int, types ...string) func(next http.Handler) http.Handler {
	// Already has deflate and gzip.
	comp := middleware.NewCompressor(lvl, types...)

	// Add brotli.
	comp.SetEncoder("br", func(w io.Writer, lvl int) io.Writer {
		return brotli.NewWriterOptions(w, brotli.WriterOptions{
			Quality: lvl,
		})
	})

	return comp.Handler
}
