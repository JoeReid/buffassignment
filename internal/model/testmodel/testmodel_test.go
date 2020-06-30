package testmodel_test

import (
	"testing"

	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/JoeReid/buffassignment/internal/model/testmodel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockGetVideoStream(t *testing.T) {
	store := testmodel.NewModelMock()
	store.On("GetVideoStream", mock.Anything).Return(&model.VideoStream{}, nil)

	v, err := store.GetVideoStream(model.VideoStreamID(uuid.New()))
	assert.Equal(t, nil, err)
	assert.Equal(t, &model.VideoStream{}, v)
}

func TestMockListVideoStream(t *testing.T) {
	store := testmodel.NewModelMock()
	store.On("ListVideoStream", mock.Anything, mock.Anything).Return([]model.VideoStream{}, nil)

	v, err := store.ListVideoStream(0, 0)
	assert.Equal(t, nil, err)
	assert.Equal(t, []model.VideoStream{}, v)
}

func TestMockCreateVideoStream(t *testing.T) {
	store := testmodel.NewModelMock()
	store.On("CreateVideoStream", mock.Anything).Return(nil)

	err := store.CreateVideoStream(model.VideoStream{})
	assert.Equal(t, nil, err)
}

func TestMockUpdateVideoStream(t *testing.T) {
	store := testmodel.NewModelMock()
	store.On("UpdateVideoStream", mock.Anything, mock.Anything).Return(nil)

	err := store.UpdateVideoStream(model.VideoStreamID(uuid.New()), model.VideoStream{})
	assert.Equal(t, nil, err)
}

func TestMockDeleteVideoStream(t *testing.T) {
	store := testmodel.NewModelMock()
	store.On("DeleteVideoStream", mock.Anything).Return(nil)

	err := store.DeleteVideoStream(model.VideoStreamID(uuid.New()))
	assert.Equal(t, nil, err)
}

func TestMockGetBuff(t *testing.T) {
	store := testmodel.NewModelMock()
	store.On("GetBuff", mock.Anything).Return(&model.Buff{}, nil)

	v, err := store.GetBuff(model.BuffID(uuid.New()))
	assert.Equal(t, nil, err)
	assert.Equal(t, &model.Buff{}, v)
}

func TestMockListBuff(t *testing.T) {
	store := testmodel.NewModelMock()
	store.On("ListBuff", mock.Anything, mock.Anything).Return([]model.Buff{}, nil)

	v, err := store.ListBuff(0, 0)
	assert.Equal(t, nil, err)
	assert.Equal(t, []model.Buff{}, v)
}

func TestMockListForStreamBuff(t *testing.T) {
	store := testmodel.NewModelMock()
	store.On("ListBuffForStream", mock.Anything, mock.Anything, mock.Anything).Return([]model.Buff{}, nil)

	v, err := store.ListBuffForStream(model.VideoStreamID(uuid.New()), 0, 0)
	assert.Equal(t, nil, err)
	assert.Equal(t, []model.Buff{}, v)
}

func TestMockCreateBuff(t *testing.T) {
	store := testmodel.NewModelMock()
	store.On("CreateBuff", mock.Anything).Return(nil)

	err := store.CreateBuff(model.Buff{})
	assert.Equal(t, nil, err)
}

func TestMockUpdateBuff(t *testing.T) {
	store := testmodel.NewModelMock()
	store.On("UpdateBuff", mock.Anything, mock.Anything).Return(nil)

	err := store.UpdateBuff(model.BuffID(uuid.New()), model.Buff{})
	assert.Equal(t, nil, err)
}

func TestMockDeleteBuff(t *testing.T) {
	store := testmodel.NewModelMock()
	store.On("DeleteBuff", mock.Anything).Return(nil)

	err := store.DeleteBuff(model.BuffID(uuid.New()))
	assert.Equal(t, nil, err)
}
