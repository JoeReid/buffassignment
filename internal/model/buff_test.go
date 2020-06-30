package model_test

import (
	"testing"

	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBuffID(t *testing.T) {
	sentinelUUID := uuid.New()

	id, err := model.ParseBuffID(sentinelUUID.String())
	require.NoError(t, err, "failed to parse buff id")

	assert.EqualValues(t, sentinelUUID, id)
}

func TestBuffIDString(t *testing.T) {
	sentinelUUID := uuid.New()

	id := model.BuffID(sentinelUUID)
	assert.Equal(t, sentinelUUID.String(), id.String())

	_, err := uuid.Parse(id.String())
	require.NoError(t, err, "failed to parse UUID")
}

func TestParseAnswerID(t *testing.T) {
	sentinelUUID := uuid.New()

	id, err := model.ParseAnswerID(sentinelUUID.String())
	require.NoError(t, err, "failed to parse answer id")

	assert.EqualValues(t, sentinelUUID, id)
}

func TestAnswerIDString(t *testing.T) {
	sentinelUUID := uuid.New()

	id := model.AnswerID(sentinelUUID)
	assert.Equal(t, sentinelUUID.String(), id.String())

	_, err := uuid.Parse(id.String())
	require.NoError(t, err, "failed to parse UUID")
}
