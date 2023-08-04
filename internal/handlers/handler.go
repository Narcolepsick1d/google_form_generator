package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"google-gen/internal/model"
	"google-gen/pkg/helper"
	"log"
	"strconv"
)

func NewHandle(opt *H) *H {
	return &H{
		User:     opt.User,
		Label:    opt.Label,
		Choice:   opt.Choice,
		Payment:  opt.Payment,
		Question: opt.Question,
	}
}
func showMessageWithUserID(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil {
			log.Printf("%d say: %s", update.Message.From.ID, update.Message.Text)
		}
		next(ctx, b, update)
	}
}

func showMessageWithUserName(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil {
			log.Printf("%s say: %s", update.Message.From.FirstName, update.Message.Text)
		}
		next(ctx, b, update)
	}
}

func (h *H) handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	err := h.User.Create(ctx, model.UserTg{
		TelegramId: strconv.Itoa(int(update.Message.From.ID)),
		UserName:   update.Message.From.Username,
		LastName:   update.Message.From.LastName,
		FirstName:  update.Message.From.FirstName,
	})
	if err != nil {
		log.Printf("error while creating user %v", err)
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
func (h *H) urlStartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if !helper.IsGoogleFormsLink(update.Message.Text) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Проверьте вашу ссылку  ",
		})
		return
	}
	qId, err := h.Question.Create(ctx, model.Question{
		UserId: strconv.Itoa(int(update.Message.From.ID)),
		Link:   update.Message.Text,
	})
	if err != nil {
		log.Printf("error while adding url %v", err)
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "сканируем ваш опросник :)",
	})
	s := helper.ExampleScrape(update.Message.Text)
	labels, htmls := helper.GetLabel(s)
	labelguid := make([]string, 0, 10)
	for _, l := range labels {
		lguid, err := h.Label.Create(ctx, model.Label{
			Name:       l.Name,
			Entry:      l.Entry,
			QuestionId: qId,
		})
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "не удалось записать ключи для опросника попробуйте снова :(",
			})
		}
		labelguid = append(labelguid, lguid)
	}
	chs := helper.GetChoices(htmls)
	for i := 0; i < len(chs); i++ {
		for j := 0; j < len(chs[i]); j++ {
			err := h.Choice.Create(ctx, model.Choices{Choice: chs[i][j], LabelId: labelguid[i]})
			if err != nil {
				return
			}
		}
	}
	resp, err := h.Question.Get(ctx, qId)
	if err != nil {
		log.Print("1", err)
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Ваши вопросы и все возможные ответы",
	})
	cant := 0
	for i, r := range resp {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   r.Name + ":  " + r.Choice,
		})
		cant = i
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Итого: %v %s %v", cant, "вариантов ответов на все вопросы, а их", len(labelguid)),
	})
	fmt.Print(cant)
}
