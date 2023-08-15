package CRUD

import "C"
import (
	"context"
	"fmt"
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
	var (
		choices []model.RespQuestionDb
		resp    []model.RespQuestion
	)
	query := `select l.guid as label_id,c.guid,l.entry as entry_id,l.name,C.choice,C.Probability  from Questions q
    join Label L on q.Id = L.question_id
    join Choices C on L.id = C.label_id
	where q.guid=$1;`
	err := q.db.Select(&choices, query, string2)
	if err != nil {
		return nil, err
	}
	end := len(choices)
	fmt.Println(end)
	chs := make([]model.Choices, 0, 10)
	for i := 1; i < end; i++ {
		//if choices[i].EntryId != choices[i-1].EntryId {
		//}
		if choices[i].EntryId == choices[i-1].EntryId {
			chs = append(chs, model.Choices{Choice: choices[i-1].Choice, Id: choices[i-1].Id, Probability: choices[i-1].Probability})
		} else {
			chs = append(chs, model.Choices{Choice: choices[i-1].Choice, Id: choices[i-1].Id, Probability: choices[i-1].Probability})
			resp = append(resp, model.RespQuestion{
				EntryId: choices[i-1].EntryId,
				Name:    choices[i-1].Name,
				Choices: chs,
				Id:      choices[i-1].LabelId,
			})
			chs = make([]model.Choices, 0, 10)
		}
		if end-1 == i {
			chs = append(chs, model.Choices{Choice: choices[i].Choice, Id: choices[i].Id, Probability: choices[i].Probability})
			resp = append(resp, model.RespQuestion{
				EntryId: choices[i].EntryId,
				Name:    choices[i].Name,
				Choices: chs,
				Id:      choices[i].LabelId,
			})
		}
	}
	return resp, nil
}
func (q QuestionRepo) GetByUrl(ctx context.Context, url string) ([]model.RespQuestionDb, error) {
	var choices []model.RespQuestionDb
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
