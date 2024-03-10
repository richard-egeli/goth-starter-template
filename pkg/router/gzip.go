package router

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodings := r.Header.Get("Accept-Encoding")

		if strings.Contains(encodings, "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			writer := gzipResponseWriter{Writer: gz, ResponseWriter: w}
			next.ServeHTTP(writer, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
