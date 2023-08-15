package service

import (
	"context"
	"google-gen/internal/model"
	"google-gen/internal/postgres/repo"
	"log"
)

func NewQuestion(repo repo.QuestionRepo) *Question {
	return &Question{repo: repo}
}

type Question struct {
	repo repo.QuestionRepo
}

func (q *Question) Create(ctx context.Context, ques model.Question) (string, error) {
	// Логирование начала выполнения метода
	log.Printf("Create method called with Question: %+v", ques)

	// Ваш код создания вопроса
	resp, err := q.repo.Create(ctx, ques)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// Логирование окончания выполнения метода
	log.Printf("Create method finished %s", resp)

	// Возвращение результата
	return resp, nil
}

func (q *Question) Get(ctx context.Context, string2 string) ([]model.RespQuestion, error) {
	// Логирование начала выполнения метода
	log.Printf("Get method called with string2: %s", string2)

	// Ваш код получения вопроса
	resp, err := q.repo.Get(ctx, string2)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Логирование окончания выполнения метода
	log.Printf("Get method finished %v", resp)

	// Возвращение результата
	return resp, nil
}

func (q *Question) GetByUrl(ctx context.Context, Url string) ([]model.RespQuestionDb, error) {
	// Логирование начала выполнения метода
	log.Printf("GetByUrl method called with URL: %s", Url)

	// Ваш код получения вопроса по URL
	resp, err := q.repo.GetByUrl(ctx, Url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Логирование окончания выполнения метода
	log.Printf("GetByUrl method finished %v", resp)

	// Возвращение результата
	return resp, nil
}
