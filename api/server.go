package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	types "github.com/oneeyedsunday/video_streaming_server/internal"
	"github.com/oneeyedsunday/video_streaming_server/pkg/video"
)

func ensureRequestIsRanged(w http.ResponseWriter, r *http.Request) (types.RangeValue, error) {
	rValue := r.Header["Range"]

	fmt.Printf("range header value is: %s\n", rValue)

	if len(strings.Trim(rValue[0], "")) == 0 {
		w.WriteHeader(http.StatusBadRequest)

		return types.RangeValue{0, 0}, errors.New("is not a range request")
	}

	rv, err := types.NewRangeValue(rValue[0])

	if err != nil {
		return rv, errors.New("is not a range request")
	}
	fmt.Printf("parse range is: %v, %v\n", rv[0], rv[1])

	return rv, nil
}

func Stream(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("got HTTP %s %s request\n", r.Method, r.URL.Path)

	rV, err := ensureRequestIsRanged(w, r)

	if err != nil {
		return
	}

	videoId := strings.TrimPrefix(r.URL.Path, "/api/video/")
	fmt.Printf("requesting video with id: %s", videoId)

	err, v := video.GetVideoById(videoId)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	err, bytes, end, videoSize, contentLength := video.SeekVideoFileByRange(v, rV)

	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	handleChunkedResponse(w, bytes, contentLength, fmt.Sprintf("bytes %v-%v/%v", rV[0], end, videoSize), "video/mp4")
}

func handleError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	msg := fmt.Sprintf("{ \"message\": \"%s\" }", err)
	io.WriteString(w, msg)
}

func handleChunkedResponse(w http.ResponseWriter, b []byte, l uint64, r string, t string) {
	w.Header().Set("X-Server-Name", "VideoStreamServer")
	w.Header().Set("Accept-Range", "bytes")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", l))
	w.Header().Set("Content-Range", r)
	w.Header().Set("Content-Type", t)
	w.WriteHeader(http.StatusPartialContent)
	w.Write(b)
}
