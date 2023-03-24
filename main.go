package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "{ \"message\": \"Server is up and running\" }")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/video/:id", getRoot)
	mux.HandleFunc("/api", getHealth)

	err := http.ListenAndServe(":3000", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
