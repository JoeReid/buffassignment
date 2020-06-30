package model

import "github.com/google/uuid"

type BuffStore interface {
	GetBuff(BuffID) (*Buff, error)
	ListBuff(offset, limit int) ([]Buff, error)
	ListBuffForStream(stream VideoStreamID, offset, limit int) ([]Buff, error)

	CreateBuff(Buff) error
	UpdateBuff(BuffID, Buff) error
	DeleteBuff(BuffID) error
}

type BuffID uuid.UUID

type Buff struct {
	ID       BuffID
	Stream   VideoStreamID
	Question string
	Answers  []Answer
}

func (b BuffID) String() string {
	return uuid.UUID(b).String()
}

func ParseBuffID(s string) (BuffID, error) {
	id, err := uuid.Parse(s)
	return BuffID(id), err
}

type Answer struct {
	ID      AnswerID
	Text    string
	Correct bool
}

type AnswerID uuid.UUID

func (a AnswerID) String() string {
	return uuid.UUID(a).String()
}

func ParseAnswerID(s string) (AnswerID, error) {
	id, err := uuid.Parse(s)
	return AnswerID(id), err
}
