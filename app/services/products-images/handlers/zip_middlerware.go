package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type GzipMiddleware struct {
}

// WrappedResponseWriter implements the http.ResponseWriter interface.
type WrappedResponseWriter struct {
	w  http.ResponseWriter
	gw *gzip.Writer
}

func NewWrappedResponseWriter(w http.ResponseWriter) *WrappedResponseWriter {
	gw := gzip.NewWriter(w)
	return &WrappedResponseWriter{
		w:  w,
		gw: gw,
	}
}

func (wr *WrappedResponseWriter) Header() http.Header {
	return wr.w.Header()
}

func (wr *WrappedResponseWriter) Write(d []byte) (int, error) {
	return wr.w.Write(d)
}

func (wr *WrappedResponseWriter) WriteHeader(statusCode int) {
	wr.w.WriteHeader(statusCode)
}

// Flush flushes any pending compressed data to the underlying writer.
func (wr *WrappedResponseWriter) Flush() {
	wr.gw.Flush()
	wr.gw.Close()
}

func (g *GzipMiddleware) GzipMiddlerware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// create a gziped response
			wrw := NewWrappedResponseWriter(w)
			wrw.Header().Set("COntent-Encoding", "gzip")
			defer wrw.Flush()

			next.ServeHTTP(wrw, r)
			return
		}

		// handler normal
		next.ServeHTTP(w, r)
	})
}
