package repo

import (
	"context"
	"google-gen/internal/model"
)

type LabelRepo interface {
	Create(ctx context.Context, name model.Label) error
	Update(ctx context.Context, name model.UpdateLabel) error
	GetByQuestionUrl(ctx context.Context, url string) ([]model.Label, error)
	GetAll(ctx context.Context) ([]model.Label, error)
}
