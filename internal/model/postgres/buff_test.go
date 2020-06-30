package postgres_test

import (
	"testing"
	"time"

	"github.com/JoeReid/buffassignment/internal/config"
	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/JoeReid/buffassignment/internal/model/postgres"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBuff(t *testing.T) {
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

	// this uuid is predictably generated by the db seed process
	sentinelUUID, err := uuid.Parse(`167939cb-6627-46e9-95af-5a25367951ba`)
	require.NoError(t, err, "failed to parse uuid")

	b, err := store.GetBuff(model.BuffID(sentinelUUID))
	require.NoError(t, err, "failed to get buff")
	assert.NotEmpty(t, b, "the buff should be populated with data")
}

func TestListBuff(t *testing.T) {
	dc, err := config.DBConfig()
	require.NoError(t, err, "failed to configure DB connection")

	store, err := postgres.NewStore(
		postgres.SetDBUser(dc.DBUser),
		postgres.SetDBPassword(dc.DBPassword),
		postgres.SetDBHostname(dc.DBHost),
		postgres.SetDBPort(dc.DBPort),
		postgres.SetDBName(dc.DBName),
		postgres.SetConnectTimeout(dc.DBConnectTimeout),
		postgres.WithMaxConLifetime(time.Minute),
		postgres.WithMaxOpenCons(1),
		postgres.WithMaxIdleCons(1),
	)
	require.NoError(t, err, "failed to create store")

	b, err := store.ListBuff(0, 0)
	require.NoError(t, err, "failed to get buff")
	assert.NotEmpty(t, b, "the buff should be populated with data")
}

func TestListBuffForStream(t *testing.T) {
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

	// this uuid is predictably generated by the db seed process
	sentinelUUID, err := uuid.Parse(`d1e2c649-8185-4ad8-a81d-0d86d1e91e00`)
	require.NoError(t, err, "failed to parse uuid")

	b, err := store.ListBuffForStream(model.VideoStreamID(sentinelUUID), 0, 0)
	require.NoError(t, err, "failed to list buff")
	assert.NotEmpty(t, b, "the buff should be populated with data")
}

func TestCreateBuff(t *testing.T) {
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

	sentinelUUID := uuid.New()
	sentinelUUID2 := uuid.New()
	sentinelUUID3 := uuid.New()

	streamUUID, err := uuid.Parse(`d1e2c649-8185-4ad8-a81d-0d86d1e91e00`)
	require.NoError(t, err, "failed to parse uuid")

	b := model.Buff{
		ID:       model.BuffID(sentinelUUID),
		Stream:   model.VideoStreamID(streamUUID),
		Question: "What is the meaning of life, the universe, and everything?",
		Answers: []model.Answer{
			{ID: model.AnswerID(sentinelUUID2), Text: "42", Correct: true},
			{ID: model.AnswerID(sentinelUUID3), Text: "43", Correct: false},
		},
	}
	err = store.CreateBuff(b)
	require.NoError(t, err, "failed to create buff")
}

func TestDeleteBuff(t *testing.T) {
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

	sentinelUUID, err := uuid.Parse(`d1e2c649-8185-4ad8-a81d-0d86d1e91e00`)
	require.NoError(t, err, "failed to parse uuid")

	err = store.DeleteBuff(model.BuffID(sentinelUUID))

	// For now this feature is not implemented
	require.Error(t, err, "not implemented")
}

func TestUpdateBuff(t *testing.T) {
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

	sentinelUUID := uuid.New()
	sentinelUUID2 := uuid.New()
	sentinelUUID3 := uuid.New()

	streamUUID, err := uuid.Parse(`d1e2c649-8185-4ad8-a81d-0d86d1e91e00`)
	require.NoError(t, err, "failed to parse uuid")

	v := model.Buff{
		ID:       model.BuffID(sentinelUUID),
		Stream:   model.VideoStreamID(streamUUID),
		Question: "What is the meaning of life, the universe, and everything?",
		Answers: []model.Answer{
			{ID: model.AnswerID(sentinelUUID2), Text: "42", Correct: true},
			{ID: model.AnswerID(sentinelUUID3), Text: "43", Correct: false},
		},
	}
	err = store.UpdateBuff(model.BuffID(sentinelUUID), v)

	// For now this feature is not implemented
	require.Error(t, err, "not implemented")
}
