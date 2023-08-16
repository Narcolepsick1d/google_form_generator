package handlers

import (
	"context"
	"github.com/go-telegram/bot"
	repo2 "google-gen/internal/postgres/repo"
)

type H struct {
	User     repo2.UserRepoI
	Label    repo2.LabelRepo
	Question repo2.QuestionRepo
	Payment  repo2.PaymentRepo
	Choice   repo2.ChoicesRepo
}

func New(ctx context.Context, telegramToken string, userService repo2.UserRepoI, quesService repo2.QuestionRepo, label repo2.LabelRepo, choice repo2.ChoicesRepo) {
	api := NewHandle(&H{User: userService, Question: quesService, Label: label, Choice: choice})
	opts := []bot.Option{
		bot.WithMiddlewares(showMessageWithUserID, showMessageWithUserName),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, api.handler),
		bot.WithDefaultHandler(api.urlStartHandler),
		bot.WithMessageTextHandler("/help", bot.MatchTypeExact, api.helpHandler),
		bot.WithMessageTextHandler("/confirm", bot.MatchTypeExact, api.confirmHandler),
	}
	b, _ := bot.New(telegramToken, opts...)

	b.Start(ctx)
}
