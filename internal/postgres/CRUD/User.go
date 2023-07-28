package CRUD

import (
	"context"

	"github.com/jmoiron/sqlx"
	"google-gen/internal/model"
	"google-gen/internal/repo"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUsers(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}
func (u UserRepo) Create(ctx context.Context, user model.UserTg) error {
	query := `INSERT INTO USERS (Telegram_Id,NickName,FirstName,LastName)
			values ($1,$2,$3,$4) on conflict (Telegram_Id) Do nothing`
	_, err := u.db.Exec(query,
		user.TelegramId,
		user.UserName,
		user.FirstName,
		user.LastName)
	if err != nil {
		return err
	}
	return nil
}

func (u UserRepo) Get(ctx context.Context, telegramId string) error {
	//TODO implement me
	panic("implement me")
}

func NewUserRepo(db *sqlx.DB) repo.UserRepoI {
	return UserRepo{db: db}
}
