package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/JoeReid/apiutils/tracer"
	"github.com/JoeReid/buffassignment/internal/config"
	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/JoeReid/buffassignment/internal/model/postgres"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
)

func main() {
	if err := tracer.InitTracer("Buff DB seed tool"); err != nil {
		tracer.UntracedLogf("failed to configure tracer: %e", err)
		os.Exit(1)
	}
	defer tracer.Close() // try to flush the traces before we exit

	os.Exit(populate())
}

func populate() (exitcode int) {
	sp, ctx := opentracing.StartSpanFromContext(context.Background(), "seed postgres database")
	defer sp.Finish()

	gofakeit.Seed(0)

	tracer.Log(sp, "read db config from environment")
	dc, err := config.DBConfig()
	if err != nil {
		tracer.SetError(sp, err)
		return 1
	}

	tracer.Log(sp, "build a new postgres store")
	store, err := postgres.NewStore(
		postgres.SetDBUser(dc.DBUser),
		postgres.SetDBPassword(dc.DBPassword),
		postgres.SetDBHostname(dc.DBHost),
		postgres.SetDBPort(dc.DBPort),
		postgres.SetDBName(dc.DBName),
		postgres.SetConnectTimeout(dc.DBConnectTimeout),
	)
	if err != nil {
		tracer.Log(sp, "failed to configure postgres store")
		tracer.SetError(sp, err)
		return 1
	}

	// Populate DB with 100 streams, each with 10 buffs
	for i := 0; i < 100; i++ {
		u, err := uuid.Parse(gofakeit.UUID())
		if err != nil {
			tracer.Log(sp, "failed to create new uuid")
			tracer.SetError(sp, err)
			return 1
		}
		vID := model.VideoStreamID(u)

		// start sometime in the last week, updated between then and now
		week := time.Hour * 24 * 7
		now := time.Now()
		startdate := gofakeit.DateRange(now.Add(-week), now)
		updated := gofakeit.DateRange(startdate, now)

		tracer.Log(sp, "create video stream entry")
		if err := store.CreateVideoStream(ctx, model.VideoStream{
			ID:        vID,
			Title:     fmt.Sprintf("%s %s stream", gofakeit.Adverb(), gofakeit.Adjective()),
			CreatedAt: startdate,
			UpdatedAt: updated,
		}); err != nil {
			tracer.Log(sp, "failed to create video stream")
			tracer.SetError(sp, err)
			return 1
		}

		// create buffs
		for j := 0; j < 10; j++ {
			u, err := uuid.Parse(gofakeit.UUID())
			if err != nil {
				tracer.Log(sp, "failed to create new uuid")
				tracer.SetError(sp, err)
				return 1
			}
			bID := model.BuffID(u)

			ans := []model.Answer{}

			// Create answers
			for k := 0; k < 5; k++ {
				u, err := uuid.Parse(gofakeit.UUID())
				if err != nil {
					tracer.Log(sp, "failed to create new uuid")
					tracer.SetError(sp, err)
					return 1
				}
				aID := model.AnswerID(u)

				ans = append(ans, model.Answer{
					ID:      aID,
					Text:    gofakeit.Noun(),
					Correct: k == 0,
				})
			}

			tracer.Log(sp, "create buff entry")
			if err := store.CreateBuff(ctx, model.Buff{
				ID:       bID,
				Stream:   vID,
				Question: gofakeit.Question(),
				Answers:  ans,
			}); err != nil {
				tracer.Log(sp, "failed to create buff")
				tracer.SetError(sp, err)
				return 1
			}
		}
	}
	return 0
}
