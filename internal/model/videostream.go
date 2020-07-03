package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// VideoStreamStore defines all the actions needed to implement a videostream storage layer
// This could be implemented by:
//   - A relational database (for production)
//   - A mock implementation (for testing)
//   - An RPC backend (for unforeseen future developments)
//
// Genericising the storage actions in this way makes the code considerably
// easier to re-factor with respect to storage sub-systems, should they need to change
type VideoStreamStore interface {
	GetVideoStream(context.Context, VideoStreamID) (*VideoStream, error)
	ListVideoStream(ctx context.Context, offset, limit int) ([]VideoStream, error)

	CreateVideoStream(context.Context, VideoStream) error
	UpdateVideoStream(context.Context, VideoStreamID, VideoStream) error
	DeleteVideoStream(context.Context, VideoStreamID) error
}

// VideoStream defines the abstract representation of the VideoStream type in the data model
//
// This is how we can think about a VideoStream in the application
// (separate from the database or API encoding representations)
type VideoStream struct {
	ID        VideoStreamID
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// VideoStreamID is a uuid.UUID type
// It is defined as it's own type to make the use of IDs in the model type-safe
// E.g. you can't accidentally use a BuffID as a VideoStreamID
type VideoStreamID uuid.UUID

// String provides access to the underlying uuid.String function
func (v VideoStreamID) String() string {
	return uuid.UUID(v).String()
}

// ParseVideoStreamID will return a new VideoStreamID parsed from it's string representation
// The string is expected in a valid uuid.UUID format
func ParseVideoStreamID(s string) (VideoStreamID, error) {
	id, err := uuid.Parse(s)
	return VideoStreamID(id), err
}
