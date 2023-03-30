package http

import (
	"log"
	"net/http"
)

// middleware is a function that wraps a handler to augment it with some extra functionality.
type middleware func(http.Handler) http.Handler

func WithLoggingRequest(h func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handling %s request for %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		h(w, r)
	})
}
