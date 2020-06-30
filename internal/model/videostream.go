package model

import (
	"time"

	"github.com/google/uuid"
)

type VideoStreamStore interface {
	GetVideoStream(VideoStreamID) (*VideoStream, error)
	ListVideoStream(offset, limit int) ([]VideoStream, error)

	CreateVideoStream(VideoStream) error
	UpdateVideoStream(VideoStreamID, VideoStream) error
	DeleteVideoStream(VideoStreamID) error
}

type VideoStream struct {
	ID        VideoStreamID
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type VideoStreamID uuid.UUID

func (v VideoStreamID) String() string {
	return uuid.UUID(v).String()
}

func ParseVideoStreamID(s string) (VideoStreamID, error) {
	id, err := uuid.Parse(s)
	return VideoStreamID(id), err
}
