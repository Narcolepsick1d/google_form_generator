package repo

import (
	"context"
	"google-gen/internal/model"
)

type Users struct {
	repo UserRepoI
}

type UserRepoI interface {
	Create(ctx context.Context, user model.UserTg) error
	Get(ctx context.Context, telegramId string) error
}
