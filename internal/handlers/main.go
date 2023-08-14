package handlers

import (
	"context"
	"github.com/go-telegram/bot"
	"google-gen/internal/repo"
)

type H struct {
	User     repo.UserRepoI
	Label    repo.LabelRepo
	Question repo.QuestionRepo
	Payment  repo.PaymentRepo
	Choice   repo.ChoicesRepo
}

func New(ctx context.Context, telegramToken string, userService repo.UserRepoI, quesService repo.QuestionRepo, label repo.LabelRepo, choice repo.ChoicesRepo) {
	api := NewHandle(&H{User: userService, Question: quesService, Label: label, Choice: choice})
	opts := []bot.Option{
		bot.WithMiddlewares(showMessageWithUserID, showMessageWithUserName),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, api.handler),
		bot.WithDefaultHandler(api.urlStartHandler),
		bot.WithMessageTextHandler("/help", bot.MatchTypeExact, api.helpHandler),
		bot.WithMessageTextHandler("/confirm", bot.MatchTypeExact, api.helpHandler),
	}
	b, _ := bot.New(telegramToken, opts...)

	b.Start(ctx)
}
