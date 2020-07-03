package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/JoeReid/buffassignment/internal/config"
	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/JoeReid/buffassignment/internal/model/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetVideoStream(t *testing.T) {
	dc, err := config.DBConfig()
	require.NoError(t, err, "failed to configure DB connection")

	store, err := postgres.NewStore(
		postgres.SetDBUser(dc.DBUser),
		postgres.SetDBPassword(dc.DBPassword),
		postgres.SetDBHostname(dc.DBHost),
		postgres.SetDBPort(dc.DBPort),
		postgres.SetDBName(dc.DBName),
		postgres.SetConnectTimeout(dc.DBConnectTimeout),
	)
	require.NoError(t, err, "failed to create store")

	// get a single videostream
	v, err := store.ListVideoStream(context.TODO(), 0, 0)
	require.NoError(t, err, "failed to list video streams")

	// look for its uuid in the db
	b, err := store.GetVideoStream(context.TODO(), v[0].ID)
	require.NoError(t, err, "failed to get video stream")
	assert.NotEmpty(t, b, "the video stream should be populated with data")
}

func TestListVideoStream(t *testing.T) {
	dc, err := config.DBConfig()
	require.NoError(t, err, "failed to configure DB connection")

	store, err := postgres.NewStore(
		postgres.SetDBUser(dc.DBUser),
		postgres.SetDBPassword(dc.DBPassword),
		postgres.SetDBHostname(dc.DBHost),
		postgres.SetDBPort(dc.DBPort),
		postgres.SetDBName(dc.DBName),
		postgres.SetConnectTimeout(dc.DBConnectTimeout),
	)
	require.NoError(t, err, "failed to create store")

	v, err := store.ListVideoStream(context.TODO(), 0, 0)
	require.NoError(t, err, "failed to list video streams")
	assert.NotEmpty(t, v, "the video stream should be populated with data")
}

func TestCreateVideoStream(t *testing.T) {
	dc, err := config.DBConfig()
	require.NoError(t, err, "failed to configure DB connection")

	store, err := postgres.NewStore(
		postgres.SetDBUser(dc.DBUser),
		postgres.SetDBPassword(dc.DBPassword),
		postgres.SetDBHostname(dc.DBHost),
		postgres.SetDBPort(dc.DBPort),
		postgres.SetDBName(dc.DBName),
		postgres.SetConnectTimeout(dc.DBConnectTimeout),
	)
	require.NoError(t, err, "failed to create store")

	now := time.Now()
	v := model.VideoStream{
		ID:        model.VideoStreamID(uuid.New()),
		Title:     "a sepcial testing stream",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Create the video stream
	err = store.CreateVideoStream(context.TODO(), v)
	require.NoError(t, err, "failed to create video stream")

	// read it back and compare
	v2, err := store.GetVideoStream(context.TODO(), v.ID)
	require.NoError(t, err, "failed to get video stream")

	assert.Equal(t, v.Title, v2.Title)
}

func TestDeleteVideoStream(t *testing.T) {
	dc, err := config.DBConfig()
	require.NoError(t, err, "failed to configure DB connection")

	store, err := postgres.NewStore(
		postgres.SetDBUser(dc.DBUser),
		postgres.SetDBPassword(dc.DBPassword),
		postgres.SetDBHostname(dc.DBHost),
		postgres.SetDBPort(dc.DBPort),
		postgres.SetDBName(dc.DBName),
		postgres.SetConnectTimeout(dc.DBConnectTimeout),
	)
	require.NoError(t, err, "failed to create store")

	// Get a single videostream to remove
	v, err := store.ListVideoStream(context.TODO(), 0, 1)
	require.NoError(t, err, "failed to list video streams")

	err = store.DeleteVideoStream(context.TODO(), v[0].ID)

	// For now this feature is not implemented
	require.Error(t, err, "not implemented")
}

func TestUpdateVideoStream(t *testing.T) {
	dc, err := config.DBConfig()
	require.NoError(t, err, "failed to configure DB connection")

	store, err := postgres.NewStore(
		postgres.SetDBUser(dc.DBUser),
		postgres.SetDBPassword(dc.DBPassword),
		postgres.SetDBHostname(dc.DBHost),
		postgres.SetDBPort(dc.DBPort),
		postgres.SetDBName(dc.DBName),
		postgres.SetConnectTimeout(dc.DBConnectTimeout),
	)
	require.NoError(t, err, "failed to create store")

	// Get a single videostream to remove
	v, err := store.ListVideoStream(context.TODO(), 0, 1)
	require.NoError(t, err, "failed to list video streams")

	// Update it
	v[0].Title = "a sepcial testing stream"
	v[0].UpdatedAt = time.Now()

	err = store.UpdateVideoStream(context.TODO(), v[0].ID, v[0])

	// For now this feature is not implemented
	require.Error(t, err, "not implemented")
}
