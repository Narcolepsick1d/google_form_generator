package repo

import (
	"context"
	"google-gen/internal/model"
)

type PaymentRepo interface {
	Create(ctx context.Context, payment model.Payment) error
	GetByUserId(ctx context.Context, userId string) (model.Payment, error)
	GetAll(ctx context.Context) ([]model.Payment, error)
}
