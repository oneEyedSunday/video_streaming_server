package video_streaming_server

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type RangeValue [2]uint64

func NewRangeValue(v string) (RangeValue, error) {
	fmt.Printf("parsing RangeValue for: %s\n", v)
	parts := strings.Split(strings.Replace(v, "bytes=", "", -1), "-")

	if parts[1] == "" {
		parts[1] = "0"
	}

	if len(parts) < 2 {
		return RangeValue{0, 0}, errors.New("invalid value")
	}

	f, err := strconv.ParseUint(parts[0], 10, 0)

	if err != nil {
		return RangeValue{0, 0}, err
	}
	s, err := strconv.ParseUint(parts[1], 10, 0)

	if err != nil {
		return RangeValue{0, 0}, err
	}

	return RangeValue{f, s}, nil
}
