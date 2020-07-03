package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/JoeReid/apiutils/tracer"
	"github.com/JoeReid/buffassignment/internal/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

// videoStream is the DB representation of the structure
type videoStream struct {
	ID      uuid.UUID
	Title   string
	Created time.Time
	Updated time.Time
}

// GetVideoStream returns a model.VideoStream by it's id
func (s *Store) GetVideoStream(ctx context.Context, id model.VideoStreamID) (*model.VideoStream, error) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "postgres get video stream")
	defer sp.Finish()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	tracer.Log(sp, "build sql")
	q, v, err := psql.Select(videoStreamFields...).From(videoStreamTable).Where("id = ?", uuid.UUID(id)).Limit(1).ToSql()
	if err != nil {
		tracer.SetError(sp, err)
		return nil, err
	}

	tracer.Logf(sp, "sql: %v", q)

	vid := videoStream{}
	if err := s.db.Get(&vid, q, v...); err != nil {
		tracer.SetError(sp, err)
		return nil, err
	}

	return &model.VideoStream{
		ID:        model.VideoStreamID(vid.ID),
		Title:     vid.Title,
		CreatedAt: vid.Created,
		UpdatedAt: vid.Updated,
	}, nil
}

// ListVideoStream returns a slice of model.VideoStream using offset and limit semantics
func (s *Store) ListVideoStream(ctx context.Context, offset, limit int) ([]model.VideoStream, error) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "postgres list video stream")
	defer sp.Finish()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	tracer.Log(sp, "build sql")
	qb := psql.Select(videoStreamFields...).From(videoStreamTable)

	if offset != 0 {
		qb = qb.Offset(uint64(offset))
	}
	if limit != 0 {
		qb = qb.Limit(uint64(limit))
	}

	q, v, err := qb.ToSql()
	if err != nil {
		tracer.SetError(sp, err)
		return nil, err
	}

	tracer.Logf(sp, "sql: %v", q)

	vids := make([]videoStream, 0)
	if err := s.db.Select(&vids, q, v...); err != nil {
		tracer.SetError(sp, err)
		return nil, err
	}

	mdlVids := make([]model.VideoStream, 0, len(vids))
	for _, vid := range vids {
		mdlVids = append(mdlVids, model.VideoStream{
			ID:        model.VideoStreamID(vid.ID),
			Title:     vid.Title,
			CreatedAt: vid.Created,
			UpdatedAt: vid.Updated,
		})
	}
	return mdlVids, nil
}

// CreateVideoStream adds a new VideoStream object into the postgres store
func (s *Store) CreateVideoStream(ctx context.Context, vid model.VideoStream) error {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "postgres create video stream")
	defer sp.Finish()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	tracer.Log(sp, "build sql")
	q, v, err := psql.Insert(videoStreamTable).Columns(videoStreamFields...).Values(
		uuid.UUID(vid.ID), vid.Title, vid.CreatedAt, vid.UpdatedAt,
	).ToSql()
	if err != nil {
		tracer.SetError(sp, err)
		return err
	}

	tracer.Logf(sp, "sql: %v", q)

	_, err = s.db.Exec(q, v...)
	if err != nil {
		tracer.SetError(sp, err)
	}

	return err
}

// UpdateVideoStream replaces the VideoStream with ID model.VideoStreamID with the given object
// This method is not yet implemented
func (s *Store) UpdateVideoStream(context.Context, model.VideoStreamID, model.VideoStream) error {
	return errors.New("not implemented")
}

// DeleteVideoStream deleted the VideoStream with ID model.VideoStreamID
// This method is not yet implemented
func (s *Store) DeleteVideoStream(context.Context, model.VideoStreamID) error {
	return errors.New("not implemented")
}
