package CRUD

import (
	"context"
	"github.com/jmoiron/sqlx"
	"google-gen/internal/model"
)

type ChoicesRepo struct {
	db *sqlx.DB
}

func NewChoice(db *sqlx.DB) *ChoicesRepo {
	return &ChoicesRepo{db: db}
}
func (c *ChoicesRepo) Create(ctx context.Context, choice model.Choices) error {
	q := "INSERT INTO Choices ( label_id, choice) values ((select id from label where guid =$1),$2) "
	_, err := c.db.Exec(q, choice.LabelId, choice.Choice)
	if err != nil {
		return err
	}
	return nil
}
func (c *ChoicesRepo) Update(ctx context.Context, choice model.UpdateChoices) error {
	q := `update choices set probability=$1,
		is_multi=$2 where guid =$3 `
	_, err := c.db.Exec(q, choice.Choice, choice.Probability, choice.Id)
	if err != nil {
		return err
	}
	return nil
}
func (c *ChoicesRepo) GetByLabelId(ctx context.Context, questionId string) ([]model.Choices, error) {
	return nil, nil
}
