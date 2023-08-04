package CRUD

import "C"
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
	query := `Insert into questions (user_id,url) values ((select id from users where telegram_id=$1),$2)  returning guid;`
	var guid string
	err := q.db.Get(&guid, query, ques.UserId, ques.Link)
	if err != nil {
		return "", err
	}
	return guid, nil
}
func (q QuestionRepo) Get(ctx context.Context, string2 string) ([]model.RespQuestion, error) {
	var choices []model.RespQuestion
	query := `select q.guid,l.name,C.choice  from Questions q
    join Label L on q.Id = L.question_id
    join Choices C on L.id = C.label_id
	where q.guid=$1;`
	err := q.db.Select(&choices, query, string2)
	if err != nil {
		return nil, err
	}
	return choices, nil
}
func (q QuestionRepo) GetByUrl(ctx context.Context, url string) ([]model.RespQuestion, error) {
	var choices []model.RespQuestion
	query := `select q.guid,l.name,C.choice  from Questions q
    join Label L on q.Id = L.question_id
    join Choices C on L.id = C.label_id
	where q.url=$1;`
	err := q.db.Select(&choices, query, url)
	if err != nil {
		return nil, err
	}
	return choices, nil
}
