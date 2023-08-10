package CRUD

import (
	"context"
	"github.com/jmoiron/sqlx"
	"google-gen/internal/model"
)

type LabelRepo struct {
	db *sqlx.DB
}

func NewLabel(db *sqlx.DB) *LabelRepo {
	return &LabelRepo{db: db}
}
func (l LabelRepo) Create(ctx context.Context, name model.Label) (string, error) {
	var guid string
	q := "INSERT INTO LABEL (entry , name ,question_id) values ($1,$2,(select id from questions where guid =$3)) returning guid"
	err := l.db.Get(&guid, q, name.Entry, name.Name, name.QuestionId)
	if err != nil {
		return "", err
	}
	return guid, nil
}
func (l LabelRepo) Update(ctx context.Context, name model.UpdateLabel) error {
	q := "Update label Set  is_multi=$1 where guid =$2"
	_, err := l.db.Exec(q, name.IsMulti, name.Id)
	if err != nil {
		return err
	}

	return nil
}
func (l LabelRepo) GetByQuestionUrl(ctx context.Context, url string) ([]model.Label, error) {
	return nil, nil
}
func (l LabelRepo) GetAll(ctx context.Context) ([]model.Label, error) { return nil, nil }
