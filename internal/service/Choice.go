package service

import (
	"context"
	"google-gen/internal/model"
	"google-gen/internal/postgres/repo"
	"log"
)

func NewChoice(repo repo.ChoicesRepo) *Choice {
	return &Choice{repo: repo}
}

type Choice struct {
	repo repo.ChoicesRepo
}

func (c *Choice) Create(ctx context.Context, choice model.Choices) error {
	// Логирование начала выполнения метода
	log.Printf("Create method called with Choice: %+v", choice)

	// Ваш код создания выбора
	err := c.repo.Create(ctx, choice)
	if err != nil {
		log.Println(err)
		return err
	}
	// Логирование окончания выполнения метода
	log.Printf("Create method finished")

	// Возвращение результата
	return nil
}

func (c *Choice) Update(ctx context.Context, choice model.UpdateChoices) error {
	// Логирование начала выполнения метода
	log.Printf("Update method called with UpdateChoices: %+v", choice)

	// Ваш код обновления выбора
	err := c.repo.Update(ctx, choice)
	if err != nil {
		log.Println(err)
		return err
	}
	// Логирование окончания выполнения метода
	log.Printf("Update method finished")

	// Возвращение результата
	return nil
}

func (c *Choice) GetByLabelId(ctx context.Context, questionId string) ([]model.Choices, error) {
	//TODO implement me
	panic("implement me")
}
