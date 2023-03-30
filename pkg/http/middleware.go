package http

import (
	"log"
	"net/http"
)

// middleware is a function that wraps a handler to augment it with some extra functionality.
type middleware func(h func(http.ResponseWriter, *http.Request)) http.HandlerFunc

type loggingResponseWriter struct {
	// from https://gist.github.com/Boerworz/b683e46ae0761056a636
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func WithLoggingRequest(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handling %s request for %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		lrw := NewLoggingResponseWriter(w)

		h.ServeHTTP(lrw, r)

		statusCode := lrw.statusCode
		log.Printf("Handled request %s %s with %d %s", r.Method, r.URL.Path, statusCode, http.StatusText(statusCode))
	})
}
