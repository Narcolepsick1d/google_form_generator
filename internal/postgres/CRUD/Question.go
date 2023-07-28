package CRUD

import (
	"context"
	"github.com/jmoiron/sqlx"
	"google-gen/internal/model"
)

type QuestionRepo struct {
	db *sqlx.DB
}

func NewQuestion(db *sqlx.DB) *QuestionRepo {
	return &QuestionRepo{db: db}
}
func (q QuestionRepo) Create(ctx context.Context, ques model.Question) (string, error) {
	query := `Insert into questions (user_id,url) values ((select id from users where telegram_id=$1),$2) returning guid`
	var guid string
	err := q.db.Get(&guid, query, ques.UserId, ques.Link)
	if err != nil {
		return "", err
	}
	return guid, nil
}
func (q QuestionRepo) Get(ctx context.Context, string2 string) (model.Question, error) {
	return model.Question{}, nil
}
