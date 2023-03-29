package video_streaming_server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Correctly_Initialiaze_RangeValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		expected RangeValue
	}{
		{
			name:    "error when init with no value",
			wantErr: true,
		},
		{
			name:     "initialize correctly on empty range",
			wantErr:  false,
			input:    "bytes=0-",
			expected: RangeValue{0, 0},
		},
		{

			name:     "initialize correctly on subsequent values",
			wantErr:  false,
			input:    "bytes=103456-",
			expected: RangeValue{103456, 0},
		},
		{
			name:     "subsequent range requests should work",
			input:    "bytes=1048576-",
			expected: RangeValue{1048576, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRangeValue(tt.input)
			if !tt.wantErr {
				assert.Nil(t, err)
				require.Equal(t, got, tt.expected)
			} else {
				require.NotNil(t, err)
			}

		})
	}

}
