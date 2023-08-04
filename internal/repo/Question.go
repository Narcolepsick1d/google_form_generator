package repo

import (
	"context"
	"google-gen/internal/model"
)

type Question struct {
	repo QuestionRepo
}
type QuestionRepo interface {
	Create(ctx context.Context, ques model.Question) (string, error)
	Get(ctx context.Context, string2 string) ([]model.RespQuestion, error)
	GetByUrl(ctx context.Context, Url string) ([]model.RespQuestion, error)
}
