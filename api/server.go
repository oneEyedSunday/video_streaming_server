package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ensureRequestIsRanged(w http.ResponseWriter, r *http.Request) error {
	rangeValue := r.Header["Range"]

	fmt.Printf("range header value is: %s\n", rangeValue)

	if len(strings.Trim(rangeValue[0], "")) == 0 {
		w.WriteHeader(http.StatusBadRequest)

		return errors.New("is not a range request")
	}

	return nil
}

func Stream(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("got HTTP %s %s request\n", r.Method, r.URL.Path)

	err := ensureRequestIsRanged(w, r)

	if err != nil {
		return
	} else {
		videoId := strings.TrimPrefix(r.URL.Path, "/api/video/")
		fmt.Printf("requesting video with id: %s", videoId)
		w.Header().Set("X-Server-Name", "VideoStreamServer")
		io.WriteString(w, "This is my website!\n")
	}

}

func writeResponseToBody() {}

func writeHeaders() {}
