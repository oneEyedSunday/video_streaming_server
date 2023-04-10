package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"bytes"
	"encoding/gob"
	"encoding/json"

	types "github.com/oneeyedsunday/video_streaming_server/internal"
	my_http "github.com/oneeyedsunday/video_streaming_server/pkg/http"
	"github.com/oneeyedsunday/video_streaming_server/pkg/video"
)

func Health(w http.ResponseWriter, r *http.Request) {
	handleJsonResponse(w, map[string]interface{}{
		"Message": "Server is up and running",
	})
}

func Stream(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("got HTTP %s %s request\n", r.Method, r.URL.Path)

	rV := r.Context().Value(my_http.RangeCtxKey).(types.RangeValue)

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

func handleResponse(w http.ResponseWriter, data interface{}) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return
	}
	w.Write(buf.Bytes())
}

func handleJsonResponse(w http.ResponseWriter, data map[string]interface{}) {
	jsonB, err := json.Marshal(data)

	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "text/json")
	w.Write(jsonB)
}
