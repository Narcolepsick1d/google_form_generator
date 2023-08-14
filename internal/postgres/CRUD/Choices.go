package CRUD

import (
	"context"
	"fmt"
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
	args := make(map[string]interface{})
	args["guid"] = choice.Id
	q := `update choices set `
	if choice.Choice != "" {
		q += `choice = :choice , `
		args["choice"] = choice.Choice
	}
	if choice.Probability != 0 {
		q += `probability = :probability , `
		args["probability"] = choice.Probability
	}
	q += ` updated_at = CURRENT_TIMESTAMP where guid = :guid`
	fmt.Print(q)
	_, err := c.db.NamedExec(q, args)
	if err != nil {
		return err
	}
	return nil
}
func (c *ChoicesRepo) GetByLabelId(ctx context.Context, questionId string) ([]model.Choices, error) {
	return nil, nil
}
