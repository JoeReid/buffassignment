package types

import (
	"time"

	"github.com/JoeReid/buffassignment/internal/model"
)

type VideoStream struct {
	UUID      string    `json:"stream_id" yaml:"stream_id"`
	Title     string    `json:"stream_title" yaml:"stream_title"`
	CreatedAt time.Time `json:"stream_created_at" yaml:"stream_created_at"`
	UpdatedAt time.Time `json:"stream_updated_at" yaml:"stream_updated_at"`
}

func NewVideoStream(mvs model.VideoStream) VideoStream {
	return VideoStream{
		UUID:      mvs.ID.String(),
		Title:     mvs.Title,
		CreatedAt: mvs.CreatedAt,
		UpdatedAt: mvs.UpdatedAt,
	}
}

func NewVideoStreams(mvs []model.VideoStream) []VideoStream {
	vs := make([]VideoStream, 0, len(mvs))

	for _, mv := range mvs {
		vs = append(vs, NewVideoStream(mv))
	}
	return vs
}
