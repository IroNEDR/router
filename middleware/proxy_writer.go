package middleware

import (
	"net/http"
)

type responseWriter interface {
	http.ResponseWriter
	Status() int
	Written() bool
	Size() int
}

type proxyWriter struct {
	http.ResponseWriter
	status      int
	size        int
	wroteHeader bool
}

func newProxyWriter(w http.ResponseWriter) responseWriter {
	return &proxyWriter{ResponseWriter: w}
}

func (pw *proxyWriter) Size() int {
	return pw.size
}

func (pw *proxyWriter) Status() int {
	return pw.status
}

func (pw *proxyWriter) Written() bool {
	return pw.wroteHeader
}

func (pw *proxyWriter) WriteHeader(code int) {
	if pw.Written() {
		return
	}

	pw.status = code
	pw.ResponseWriter.WriteHeader(code)
	pw.wroteHeader = true
}

func (pw *proxyWriter) Write(b []byte) (int, error) {
	if !pw.Written() {
		pw.WriteHeader(http.StatusOK)
	}
	size, err := pw.ResponseWriter.Write(b)
	pw.size += size
	return size, err
}
