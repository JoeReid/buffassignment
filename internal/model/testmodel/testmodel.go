package testmodel

import (
	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/stretchr/testify/mock"
)

var _ model.Store = &modelMock{}

type modelMock struct {
	mock.Mock
}

func (m *modelMock) GetVideoStream(v model.VideoStreamID) (*model.VideoStream, error) {
	args := m.MethodCalled("GetVideoStream", v)
	return args.Get(0).(*model.VideoStream), args.Error(1)
}

func (m *modelMock) ListVideoStream(offset, limit int) ([]model.VideoStream, error) {
	args := m.MethodCalled("ListVideoStream", offset, limit)
	return args.Get(0).([]model.VideoStream), args.Error(1)
}

func (m *modelMock) CreateVideoStream(v model.VideoStream) error {
	args := m.MethodCalled("CreateVideoStream", v)
	return args.Error(0)
}

func (m *modelMock) UpdateVideoStream(i model.VideoStreamID, v model.VideoStream) error {
	args := m.MethodCalled("UpdateVideoStream", i, v)
	return args.Error(0)
}

func (m *modelMock) DeleteVideoStream(v model.VideoStreamID) error {
	args := m.MethodCalled("DeleteVideoStream", v)
	return args.Error(0)
}

func (m *modelMock) GetBuff(b model.BuffID) (*model.Buff, error) {
	args := m.MethodCalled("GetBuff", b)
	return args.Get(0).(*model.Buff), args.Error(1)
}

func (m *modelMock) ListBuff(offset, limit int) ([]model.Buff, error) {
	args := m.MethodCalled("ListBuff", offset, limit)
	return args.Get(0).([]model.Buff), args.Error(1)
}

func (m *modelMock) ListBuffForStream(stream model.VideoStreamID, offset, limit int) ([]model.Buff, error) {
	args := m.MethodCalled("ListBuffForStream", stream, offset, limit)
	return args.Get(0).([]model.Buff), args.Error(1)
}

func (m *modelMock) CreateBuff(b model.Buff) error {
	args := m.MethodCalled("CreateBuff", b)
	return args.Error(0)
}

func (m *modelMock) UpdateBuff(i model.BuffID, b model.Buff) error {
	args := m.MethodCalled("UpdateBuff", i, b)
	return args.Error(0)
}

func (m *modelMock) DeleteBuff(b model.BuffID) error {
	args := m.MethodCalled("DeleteBuff", b)
	return args.Error(0)
}

func NewModelMock() *modelMock { return &modelMock{} }
