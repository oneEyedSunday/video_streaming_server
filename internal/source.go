package video_streaming_server

type VideoSource struct {
	url string
}

func NewVideoSourceFromFilePath(f string) *VideoSource {
	return &VideoSource{
		url: f,
	}
}

func (v *VideoSource) GetUrl() string {
	return v.url
}
