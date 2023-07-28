package repo

import (
	"context"
	"google-gen/internal/model"
)

type ChoicesRepo interface {
	Create(ctx context.Context, choice model.Choices) error
	Update(ctx context.Context, choice model.UpdateChoices) error
	GetByLabelId(ctx context.Context, questionId string) ([]model.Choices, error)
}
