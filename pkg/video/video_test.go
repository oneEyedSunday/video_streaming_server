package video

import (
	"testing"

	types "github.com/oneeyedsunday/video_streaming_server/internal"
	"github.com/stretchr/testify/assert"
)

func Test_Correctly_Seeks(t *testing.T) {
	c := uint64(GetChunkSize())
	tests := []struct {
		name      string
		seekRange types.RangeValue
	}{
		{
			name:      "seeks from beginning",
			seekRange: types.RangeValue{0, 0},
		},
		{
			name:      "seeks intermediate point",
			seekRange: types.RangeValue{c, 0},
		},
		{
			name:      "seeks correctly with custom intermediate point",
			seekRange: types.RangeValue{c * 3, c * 4},
		},
	}

	videoSource := types.NewVideoSourceFromFilePath(resolveFilePath("video"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, b, end, _, _ := SeekVideoFileByRange(videoSource, tt.seekRange)

			assert.Nil(t, err)
			assert.NotNil(t, b)
			assert.Equal(t, tt.seekRange[0]+c-1, end)
			assert.Equal(t, c, uint64(len(b)))
		})
	}
}
