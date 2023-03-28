package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/oneeyedsunday/video_streaming_server/api"
)

func getHealth(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "{ \"message\": \"Server is up and running\" }")
}

const keyServerAddr = "serverAddr"

func main() {
	mux := http.NewServeMux()
	// Wrap handlers in request and response logger middlewares
	mux.HandleFunc("/api/video/", api.Stream)
	mux.HandleFunc("/api", getHealth)

	ctx, cancelCtx := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:    ":3000",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	go func() {
		fmt.Printf("Booting up server on: %s \n", server.Addr)
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else if err != nil {
			fmt.Printf("error starting server: %s\n", err)
		}
		cancelCtx()
	}()

	<-ctx.Done()
}
