package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	types "github.com/oneeyedsunday/video_streaming_server/internal"
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

const RangeCtxKey = "range"

func ensureRequestIsRanged(w http.ResponseWriter, r *http.Request) error {
	rValue := r.Header["Range"]

	fmt.Printf("range header value is: %s\n", rValue)

	if len(rValue) == 0 {
		return errors.New("is not a range request")
	}

	if len(strings.Trim(rValue[0], "")) == 0 {

		return errors.New("is not a range request")
	}

	rv, err := types.NewRangeValue(rValue[0])

	if err != nil {
		return errors.New("is not a range request")
	}
	fmt.Printf("parse range is: %v, %v\n", rv[0], rv[1])

	*r = *r.WithContext(context.WithValue(r.Context(), RangeCtxKey, rv))

	return nil
}

func EnsureRequestIsRanged(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := ensureRequestIsRanged(w, r)
		if err != nil {

			log.Printf("Non ranged request received %s request for %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			http.Error(w, "is not a range request", http.StatusBadRequest)
			return
		}
		h.ServeHTTP(w, r)
	})
}
