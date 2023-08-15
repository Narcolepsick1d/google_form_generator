package service

import (
	"context"
	"google-gen/internal/model"
	"google-gen/internal/postgres/repo"
	"log"
)

func NewLabel(repo repo.LabelRepo) *Label {
	return &Label{repo: repo}
}

type Label struct {
	repo repo.LabelRepo
}

func (l *Label) Create(ctx context.Context, name model.Label) (string, error) {
	// Логирование начала выполнения метода
	log.Printf("Create method called with Label: %+v", name)

	// Ваш код создания метки
	resp, err := l.repo.Create(ctx, name)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// Логирование окончания выполнения метода
	log.Printf("Create method finished, returned ID: %s", resp)

	// Возвращение результата
	return resp, nil
}

func (l *Label) Update(ctx context.Context, name model.UpdateLabel) error {
	// Логирование начала выполнения метода
	log.Printf("Update method called with UpdateLabel: %+v", name)

	// Ваш код обновления метки
	err := l.repo.Update(ctx, name)
	if err != nil {
		log.Print(err)
		return err
	}
	// Логирование окончания выполнения метода
	log.Printf("Update method finished")

	// Возвращение результата
	return nil
}
func (l *Label) GetByQuestionUrl(ctx context.Context, url string) ([]model.Label, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Label) GetAll(ctx context.Context) ([]model.Label, error) {
	//TODO implement me
	panic("implement me")
}
