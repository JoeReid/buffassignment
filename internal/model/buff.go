package model

import "github.com/google/uuid"

// BuffStore defines all the actions needed to implement a buff storage layer
// This could be implemented by:
//   - A relational database (for production)
//   - A mock implementation (for testing)
//   - An RPC backend (for unforeseen future developments)
//
// Genericising the storage actions in this way makes the code considerably
// easier to re-factor with respect to storage sub-systems, should they need to change
type BuffStore interface {
	GetBuff(BuffID) (*Buff, error)
	ListBuff(offset, limit int) ([]Buff, error)
	ListBuffForStream(stream VideoStreamID, offset, limit int) ([]Buff, error)

	CreateBuff(Buff) error
	UpdateBuff(BuffID, Buff) error
	DeleteBuff(BuffID) error
}

// BuffID is a uuid.UUID type
// It is defined as it's own type to make the use of IDs in the model type-safe
// E.g. you can't accidentally use a BuffID as a VideoStreamID
type BuffID uuid.UUID

// String provides access to the underlying uuid.String function
func (b BuffID) String() string {
	return uuid.UUID(b).String()
}

// ParseBuffID will return a new BuffID parsed from it's string representation
// The string is expected in a valid uuid.UUID format
func ParseBuffID(s string) (BuffID, error) {
	id, err := uuid.Parse(s)
	return BuffID(id), err
}

// Buff defines the abstract representation of the Buff type in the data model
//
// This is how we can think about a Buff in the application
// (separate from the database or API encoding representations)
type Buff struct {
	ID       BuffID
	Stream   VideoStreamID
	Question string
	Answers  []Answer
}

// Answer defines the abstract representation of the Answer type in the data model
// It is not meant to be used in isolation from a Buff type
type Answer struct {
	ID      AnswerID
	Text    string
	Correct bool
}

// AnswerID is a uuid.UUID type
// It is defined as it's own type to make the use of IDs in the model type-safe
// E.g. you can't accidentally use a BuffID as a VideoStreamID
type AnswerID uuid.UUID

// String provides access to the underlying uuid.String function
func (a AnswerID) String() string {
	return uuid.UUID(a).String()
}

// ParseAnswerID will return a new AnswerID parsed from it's string representation
// The string is expected in a valid uuid.UUID format
func ParseAnswerID(s string) (AnswerID, error) {
	id, err := uuid.Parse(s)
	return AnswerID(id), err
}
