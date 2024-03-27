package router

import (
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/go-chi/chi/v5/middleware"
)

func compress(level int, types ...string) [](func(http.Handler) http.Handler) {
	// deflate & gzip
	dig := middleware.Compress(level, types...)

	// brotli
	br := middleware.NewCompressor(level, types...)
	br.SetEncoder("br", func(w io.Writer, level int) io.Writer {
		return brotli.NewWriterOptions(w, brotli.WriterOptions{
			Quality: level,
		})
	})

	return [](func(http.Handler) http.Handler){dig, br.Handler}
}
