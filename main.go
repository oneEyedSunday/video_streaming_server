package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/oneeyedsunday/video_streaming_server/api"
	my_http "github.com/oneeyedsunday/video_streaming_server/pkg/http"
)

const keyServerAddr = "serverAddr"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/video/", my_http.WithLoggingRequest(http.HandlerFunc(api.Stream)))
	mux.HandleFunc("/api", my_http.WithLoggingRequest(http.HandlerFunc(api.Health)))

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
