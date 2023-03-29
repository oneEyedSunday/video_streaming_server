package video

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"strings"

	types "github.com/oneeyedsunday/video_streaming_server/internal"
)

var (
	videoMap = map[string]string{
		// mkv may not be playable
		"kimetsu":  "/Users/ispoa/Downloads/[Hakata Ramen] Kimetsu no Yaiba (Demon Slayer) {Season 1} [1080p][HEVC][10bit][Opus][Multi-Subs](Doc_Ramen)/[Hakata Ramen] Kimetsu no Yaiba (Demon Slayer) - 01 [1080p][HEVC].mkv",
		"rashomon": "/Users/ispoa/Downloads/Rashomon (1950) [BluRay] [1080p] [YTS.AM]/Rashomon.1950.1080p.BluRay.x264-[YTS.AM].mp4",
		"drive":    "/Users/ispoa/Downloads/Drive (2011) [1080p] [BluRay] [5.1] [YTS.MX]/Drive.2011.1080p.BluRay.x264.AAC5.1-[YTS.MX].mp4",
		"nayakan":  "/Users/ispoa/Downloads/Nayakan (1987) 720p 10bit AMZN WEBRip x265 HEVC Tamil DDP 2.0 ESub ~ Immortal.mkv",
		"video":    "/Users/ispoa/Downloads/Rashomon (1950) [BluRay] [1080p] [YTS.AM]/Rashomon.1950.1080p.BluRay.x264-[YTS.AM].mp4",
	}
)

// 1MB
const chunkSize = 1024 * 1024

func GetChunkSize() int {
	return chunkSize
}

func GetVideoById(id string) (error, *types.VideoSource) {
	realPath := resolveFilePath(id)

	if realPath == "" {
		return errors.New("video not found"), nil
	}

	return nil, types.NewVideoSourceFromFilePath(realPath)
}

func SeekVideoFileByRange(video *types.VideoSource, seekRange types.RangeValue) (err error, bytes []byte, end uint64, videoSize int64, contentLength uint64) {
	file, err := os.Open(video.GetUrl())

	if err != nil {
		return handleError(err), nil, 0, 0, 0
	}

	defer file.Close()

	fs, err := file.Stat()

	if err != nil {
		return handleError(err), nil, 0, 0, 0
	}
	// size in bytes
	videoSize = fs.Size()

	fmt.Printf("file info is: %v and extension is: %s\n", videoSize, path.Ext(fs.Name()))

	//
	// Calculate end Content-Range
	//
	// Safari/iOS first sends a request with bytes=0-1 range HTTP header
	// probably to find out if the server supports byte ranges
	//

	fmt.Println(seekRange)

	fmt.Printf("Difference of seekRange is: %v\n", seekRange[1]-seekRange[0])

	end = seekRange[1]

	if seekRange[1] != 1 {
		end = uint64(math.Min(float64(seekRange[0]+chunkSize), float64(videoSize)) - 1)
	}

	contentLength = (end - seekRange[0]) + 1 // We add 1 byte because start and end start from 0

	fmt.Printf("Raw end is: %v\n", end)

	err, bytes = returnChunk(file, seekRange[0])

	if err != nil {
		return handleError(err), nil, 0, 0, 0
	}
	// return nil, end, videoSize, contentLength
	return
}

func resolveFilePath(id string) string {

	for k, v := range videoMap {
		if strings.HasPrefix(strings.ToLower(id), strings.ToLower(k)) {
			return v
		}
	}
	return ""
}

func returnChunk(f *os.File, start uint64) (error, []byte) {
	b := make([]byte, chunkSize)

	f.ReadAt(b, int64(start))
	// for {
	// read content to buffer
	_, err := f.Read(b)
	if err != nil {
		if err != io.EOF {
			fmt.Println(err)
			return err, nil
		}
		// break
	}
	// }
	return nil, b
}

func handleError(err error) error {
	// Err return custom errors
	// Issue with file
	// File not found
	// Replace file with video
	fmt.Println(err)
	return err
}
