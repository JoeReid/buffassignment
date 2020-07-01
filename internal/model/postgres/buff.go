package postgres

import (
	"errors"

	"github.com/JoeReid/apiutils/tracer"
	"github.com/JoeReid/buffassignment/internal/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

// question is the DB representation of the structure
type question struct {
	ID     uuid.UUID
	Stream uuid.UUID
	Text   string
}

// answer is the DB representation of the structure
type answer struct {
	ID       uuid.UUID
	Question uuid.UUID
	Text     string
	Correct  bool
}

// GetBuff returns a model.Buff by it's id
func (s *Store) GetBuff(id model.BuffID) (*model.Buff, error) {
	// TODO: replace with context method
	sp := opentracing.GlobalTracer().StartSpan("Postgres:Get Buff")
	defer sp.Finish()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, v, err := psql.Select(buffFields...).From(questionTable).Join(
		answerTable+" ON questions.id = answers.question",
	).Where("questions.id = ?", uuid.UUID(id)).Limit(1).ToSql()
	if err != nil {
		tracer.Log(sp, "failed to build sql query")
		tracer.SetError(sp, err)
		return nil, err
	}

	res, err := s.db.Queryx(q, v...)
	if err != nil {
		tracer.Log(sp, "failed to run query")
		tracer.SetError(sp, err)
		return nil, err
	}

	mdlBuff := model.Buff{Answers: make([]model.Answer, 0)}
	for res.Next() {
		ques := question{}
		ans := answer{}

		if err := res.Scan(
			&ques.ID, &ques.Stream, &ques.Text,
			&ans.ID, &ans.Question, &ans.Text, &ans.Correct,
		); err != nil {
			tracer.Log(sp, "failed to scan results")
			tracer.SetError(sp, err)
			return nil, err
		}

		mdlBuff.ID = model.BuffID(ques.ID)
		mdlBuff.Stream = model.VideoStreamID(ques.Stream)
		mdlBuff.Question = ques.Text
		mdlBuff.Answers = append(mdlBuff.Answers, model.Answer{
			ID:      model.AnswerID(ans.ID),
			Text:    ans.Text,
			Correct: ans.Correct,
		})
	}

	return &mdlBuff, nil
}

// ListBuff returns a slice of model.Buff using offset and limit semantics
func (s *Store) ListBuff(offset, limit int) ([]model.Buff, error) {
	// TODO: replace with context method
	sp := opentracing.GlobalTracer().StartSpan("Postgres:List Buff")
	defer sp.Finish()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	qb := psql.Select(buffFields...).From(questionTable).Join(
		answerTable + " ON questions.id = answers.question",
	)

	if offset != 0 {
		qb = qb.Offset(uint64(offset))
	}
	if limit != 0 {
		qb = qb.Limit(uint64(limit))
	}

	q, v, err := qb.ToSql()
	if err != nil {
		tracer.Log(sp, "failed to build sql query")
		tracer.SetError(sp, err)
		return nil, err
	}

	res, err := s.db.Queryx(q, v...)
	if err != nil {
		return nil, err
	}

	mdlBuffs := make(map[uuid.UUID]model.Buff)

	for res.Next() {
		ques := question{}
		ans := answer{}

		if err := res.Scan(
			&ques.ID, &ques.Stream, &ques.Text,
			&ans.ID, &ans.Question, &ans.Text, &ans.Correct,
		); err != nil {
			return nil, err
		}

		// Create blank object if not exist
		mdlBuff, ok := mdlBuffs[ques.ID]
		if !ok {
			mdlBuff = model.Buff{Answers: make([]model.Answer, 0)}
		}

		mdlBuff.ID = model.BuffID(ques.ID)
		mdlBuff.Stream = model.VideoStreamID(ques.Stream)
		mdlBuff.Question = ques.Text
		mdlBuff.Answers = append(mdlBuff.Answers, model.Answer{
			ID:      model.AnswerID(ans.ID),
			Text:    ans.Text,
			Correct: ans.Correct,
		})
		mdlBuffs[ques.ID] = mdlBuff
	}

	rtn := make([]model.Buff, 0, len(mdlBuffs))
	for _, v := range mdlBuffs {
		rtn = append(rtn, v)
	}
	return rtn, nil
}

// ListBuffForStream returns a slice of model.Buff using offset and limit semantics
// Where all the returned buffs are ascociated with the given model.VideoStreamID
func (s *Store) ListBuffForStream(stream model.VideoStreamID, offset, limit int) ([]model.Buff, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	qb := psql.Select(buffFields...).From(questionTable).Join(
		answerTable+" ON questions.id = answers.question",
	).Where("questions.stream = ?", uuid.UUID(stream))

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

	res, err := s.db.Queryx(q, v...)
	if err != nil {
		return nil, err
	}

	mdlBuffs := make(map[uuid.UUID]model.Buff)

	for res.Next() {
		ques := question{}
		ans := answer{}

		if err := res.Scan(
			&ques.ID, &ques.Stream, &ques.Text,
			&ans.ID, &ans.Question, &ans.Text, &ans.Correct,
		); err != nil {
			return nil, err
		}

		// Create blank object if not exist
		mdlBuff, ok := mdlBuffs[ques.ID]
		if !ok {
			mdlBuff = model.Buff{Answers: make([]model.Answer, 0)}
		}

		mdlBuff.ID = model.BuffID(ques.ID)
		mdlBuff.Stream = model.VideoStreamID(ques.Stream)
		mdlBuff.Question = ques.Text
		mdlBuff.Answers = append(mdlBuff.Answers, model.Answer{
			ID:      model.AnswerID(ans.ID),
			Text:    ans.Text,
			Correct: ans.Correct,
		})
		mdlBuffs[ques.ID] = mdlBuff
	}

	rtn := make([]model.Buff, 0, len(mdlBuffs))
	for _, v := range mdlBuffs {
		rtn = append(rtn, v)
	}
	return rtn, nil
}

// CreateBuff adds a new buff object into the postgres store
func (s *Store) CreateBuff(buff model.Buff) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, v, err := psql.Insert(questionTable).Columns(questionFields...).Values(
		uuid.UUID(buff.ID), uuid.UUID(buff.Stream), buff.Question,
	).ToSql()
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(q, v...)
	if err != nil {
		// No need to check the error here,
		// just make a best attempt to clean up the transaction
		// nolint:errcheck
		defer tx.Rollback()
		return err
	}

	for _, ans := range buff.Answers {
		q2, v2, err := psql.Insert(answerTable).Columns(answerFields...).Values(
			uuid.UUID(ans.ID), uuid.UUID(buff.ID), ans.Text, ans.Correct,
		).ToSql()
		if err != nil {
			// No need to check the error here,
			// just make a best attempt to clean up the transaction
			// nolint:errcheck
			defer tx.Rollback()
			return err
		}

		_, err = tx.Exec(q2, v2...)
		if err != nil {
			// No need to check the error here,
			// just make a best attempt to clean up the transaction
			// nolint:errcheck
			defer tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// UpdateBuff replaces the Buff with ID model.BuffID with the given object
// This method is not yet implemented
func (s *Store) UpdateBuff(model.BuffID, model.Buff) error {
	return errors.New("not implemented")
}

// DeleteBuff deleted the Buff with ID model.BuffID
// This method is not yet implemented
func (s *Store) DeleteBuff(model.BuffID) error {
	return errors.New("not implemented")
}
