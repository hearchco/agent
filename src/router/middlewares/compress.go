package middlewares

import (
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/go-chi/chi/v5/middleware"
)

func compress(lvl int, types ...string) [](func(http.Handler) http.Handler) {
	// Deflate & GZIP.
	dig := middleware.Compress(lvl, types...)

	// Brotli.
	br := middleware.NewCompressor(lvl, types...)
	br.SetEncoder("br", func(w io.Writer, lvl int) io.Writer {
		return brotli.NewWriterOptions(w, brotli.WriterOptions{
			Quality: lvl,
		})
	})

	return [](func(http.Handler) http.Handler){dig, br.Handler}
}
