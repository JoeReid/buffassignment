package model_test

import (
	"testing"

	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVideoStreamID(t *testing.T) {
	sentinelUUID := uuid.New()

	id, err := model.ParseVideoStreamID(sentinelUUID.String())
	require.NoError(t, err, "failed to parse video stream id")

	assert.EqualValues(t, sentinelUUID, id)
}

func TestVideoStreamIDString(t *testing.T) {
	sentinelUUID := uuid.New()

	id := model.VideoStreamID(sentinelUUID)
	assert.Equal(t, sentinelUUID.String(), id.String())

	_, err := uuid.Parse(id.String())
	require.NoError(t, err, "failed to parse UUID")
}
