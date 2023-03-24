package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Stream(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("got HTTP %s %s request\n", r.Method, r.URL.Path)
	videoId := strings.TrimPrefix(r.URL.Path, "/api/video/")
	fmt.Printf("requesting video with id: %s", videoId)
	w.Header().Set("X-Server-Name", "VideoStreamServer")
	io.WriteString(w, "This is my website!\n")
}

func writeResponseToBody() {}

func writeHeaders() {}
