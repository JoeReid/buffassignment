package postgres

import (
	"errors"
	"time"

	"github.com/JoeReid/buffassignment/internal/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// videoStream is the DB representation of the structure
type videoStream struct {
	ID      uuid.UUID
	Title   string
	Created time.Time
	Updated time.Time
}

// GetVideoStream returns a model.VideoStream by it's id
func (s *Store) GetVideoStream(id model.VideoStreamID) (*model.VideoStream, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, v, err := psql.Select(videoStreamFields...).From(videoStreamTable).Where("id = ?", uuid.UUID(id)).Limit(1).ToSql()
	if err != nil {
		return nil, err
	}

	vid := videoStream{}
	if err := s.db.Get(&vid, q, v...); err != nil {
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
func (s *Store) ListVideoStream(offset, limit int) ([]model.VideoStream, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	qb := psql.Select(videoStreamFields...).From(videoStreamTable)

	if offset != 0 {
		qb = qb.Offset(uint64(offset))
	}
	if limit != 0 {
		qb = qb.Limit(uint64(limit))
	}

	q, v, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	vids := make([]videoStream, 0)
	if err := s.db.Select(&vids, q, v...); err != nil {
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
func (s *Store) CreateVideoStream(vid model.VideoStream) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, v, err := psql.Insert(videoStreamTable).Columns(videoStreamFields...).Values(
		uuid.UUID(vid.ID), vid.Title, vid.CreatedAt, vid.UpdatedAt,
	).ToSql()
	if err != nil {
		return err
	}
	_, err = s.db.Exec(q, v...)
	return err
}

// UpdateVideoStream replaces the VideoStream with ID model.VideoStreamID with the given object
// This method is not yet implemented
func (s *Store) UpdateVideoStream(model.VideoStreamID, model.VideoStream) error {
	return errors.New("not implemented")
}

// DeleteVideoStream deleted the VideoStream with ID model.VideoStreamID
// This method is not yet implemented
func (s *Store) DeleteVideoStream(model.VideoStreamID) error {
	return errors.New("not implemented")
}
