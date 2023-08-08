package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
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
		}
		next(ctx, b, update)
	}
}

func showMessageWithUserName(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil {
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
		Text:   "Привет этот бот умеет генерировать ответы для google forms с определенным процентажем на каждом вопросе просто скинь ссылку и начнем )",
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "ПОЖАЛУЙСТА ЗАПОЛНЯЙТЕ ПРАВИЛЬНО",
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

	for _, r := range resp {
		lena := len(r.Choices)
		if lena > 12 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Вопрос с 1 или больше 12 вариантов ответа не принимаются",
			})
		}
		switch lena {
		case 2:
			kb := inline.New(b).
				Row().
				Button(r.Choices[0].Choice, []byte(r.Choices[0].Id), updateChoices).
				Button(r.Choices[1].Choice, []byte(r.Choices[1].Id), updateChoices).
				Row().
				Button("Мульти выбор", []byte(r.Id), updateChoices)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        r.Name + ":",
				ReplyMarkup: kb,
			})
		case 3:
			kb := inline.New(b).
				Row().
				Button(r.Choices[0].Choice, []byte(r.Choices[0].Id), updateChoices).
				Button(r.Choices[1].Choice, []byte(r.Choices[1].Id), updateChoices).
				Button(r.Choices[2].Choice, []byte(r.Choices[2].Id), updateChoices).
				Row().
				Button("Мульти выбор", []byte(r.Id), updateChoices)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        r.Name + ":",
				ReplyMarkup: kb,
			})
		case 4:
			kb := inline.New(b).
				Row().
				Button(r.Choices[0].Choice, []byte(r.Choices[0].Id), updateChoices).
				Button(r.Choices[1].Choice, []byte(r.Choices[1].Id), updateChoices).
				Row().
				Button(r.Choices[2].Choice, []byte(r.Choices[2].Id), updateChoices).
				Button(r.Choices[3].Choice, []byte(r.Choices[3].Id), updateChoices).
				Row().
				Button("Мульти выбор", []byte(r.Id), updateChoices)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        r.Name + ":  ",
				ReplyMarkup: kb,
			})
		case 5:
			kb := inline.New(b).
				Row().
				Button(r.Choices[0].Choice, []byte(r.Choices[0].Id), updateChoices).
				Button(r.Choices[1].Choice, []byte(r.Choices[1].Id), updateChoices).
				Row().
				Button(r.Choices[2].Choice, []byte(r.Choices[2].Id), updateChoices).
				Button(r.Choices[3].Choice, []byte(r.Choices[3].Id), updateChoices).
				Button(r.Choices[4].Choice, []byte(r.Choices[4].Id), updateChoices).
				Row().
				Button("Мульти выбор", []byte(r.Id), updateChoices)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        r.Name + ":  ",
				ReplyMarkup: kb,
			})
		case 6:
			kb := inline.New(b).
				Row().
				Button(r.Choices[0].Choice, []byte(r.Choices[0].Id), updateChoices).
				Button(r.Choices[1].Choice, []byte(r.Choices[1].Id), updateChoices).
				Button(r.Choices[2].Choice, []byte(r.Choices[2].Id), updateChoices).
				Row().
				Button(r.Choices[3].Choice, []byte(r.Choices[3].Id), updateChoices).
				Button(r.Choices[4].Choice, []byte(r.Choices[4].Id), updateChoices).
				Button(r.Choices[5].Choice, []byte(r.Choices[5].Id), updateChoices).
				Row().
				Button("Мульти выбор", []byte(r.Id), updateChoices)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        r.Name + ":  ",
				ReplyMarkup: kb,
			})
		case 7:
			kb := inline.New(b).
				Row().
				Button(r.Choices[0].Choice, []byte(r.Choices[0].Id), updateChoices).
				Button(r.Choices[1].Choice, []byte(r.Choices[1].Id), updateChoices).
				Button(r.Choices[2].Choice, []byte(r.Choices[2].Id), updateChoices).
				Row().
				Button(r.Choices[3].Choice, []byte(r.Choices[3].Id), updateChoices).
				Button(r.Choices[4].Choice, []byte(r.Choices[4].Id), updateChoices).
				Button(r.Choices[5].Choice, []byte(r.Choices[5].Id), updateChoices).
				Row().
				Button(r.Choices[6].Choice, []byte(r.Choices[6].Id), updateChoices).
				Row().
				Button("Мульти выбор", []byte(r.Id), updateChoices)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        r.Name + ":  ",
				ReplyMarkup: kb,
			})
		case 8:
			kb := inline.New(b).
				Row().
				Button(r.Choices[0].Choice, []byte(r.Choices[0].Id), updateChoices).
				Button(r.Choices[1].Choice, []byte(r.Choices[1].Id), updateChoices).
				Button(r.Choices[2].Choice, []byte(r.Choices[2].Id), updateChoices).
				Row().
				Button(r.Choices[3].Choice, []byte(r.Choices[3].Id), updateChoices).
				Button(r.Choices[4].Choice, []byte(r.Choices[4].Id), updateChoices).
				Button(r.Choices[5].Choice, []byte(r.Choices[5].Id), updateChoices).
				Row().
				Button(r.Choices[6].Choice, []byte(r.Choices[6].Id), updateChoices).
				Button(r.Choices[7].Choice, []byte(r.Choices[7].Id), updateChoices).
				Row().
				Button("Мульти выбор", []byte(r.Id), updateChoices)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        r.Name + ":  ",
				ReplyMarkup: kb,
			})
		case 9:
			kb := inline.New(b).
				Row().
				Button(r.Choices[0].Choice, []byte(r.Choices[0].Id), updateChoices).
				Button(r.Choices[1].Choice, []byte(r.Choices[1].Id), updateChoices).
				Button(r.Choices[2].Choice, []byte(r.Choices[2].Id), updateChoices).
				Row().
				Button(r.Choices[3].Choice, []byte(r.Choices[3].Id), updateChoices).
				Button(r.Choices[4].Choice, []byte(r.Choices[4].Id), updateChoices).
				Button(r.Choices[5].Choice, []byte(r.Choices[5].Id), updateChoices).
				Row().
				Button(r.Choices[6].Choice, []byte(r.Choices[6].Id), updateChoices).
				Button(r.Choices[7].Choice, []byte(r.Choices[7].Id), updateChoices).
				Button(r.Choices[8].Choice, []byte(r.Choices[8].Id), updateChoices).
				Row().
				Button("Мульти выбор", []byte(r.Id), updateChoices)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        r.Name + ":  ",
				ReplyMarkup: kb,
			})
		case 10:
			kb := inline.New(b).
				Row().
				Button(r.Choices[0].Choice, []byte(r.Choices[0].Id), updateChoices).
				Button(r.Choices[1].Choice, []byte(r.Choices[1].Id), updateChoices).
				Button(r.Choices[2].Choice, []byte(r.Choices[2].Id), updateChoices).
				Row().
				Button(r.Choices[3].Choice, []byte(r.Choices[3].Id), updateChoices).
				Button(r.Choices[4].Choice, []byte(r.Choices[4].Id), updateChoices).
				Button(r.Choices[5].Choice, []byte(r.Choices[5].Id), updateChoices).
				Row().
				Button(r.Choices[6].Choice, []byte(r.Choices[6].Id), updateChoices).
				Button(r.Choices[7].Choice, []byte(r.Choices[7].Id), updateChoices).
				Button(r.Choices[8].Choice, []byte(r.Choices[8].Id), updateChoices).
				Row().
				Button(r.Choices[9].Choice, []byte(r.Choices[9].Id), updateChoices).
				Row().
				Button("Мульти выбор", []byte(r.Id), updateChoices)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        r.Name + ":  ",
				ReplyMarkup: kb,
			})

		case 11:
			kb := inline.New(b).
				Row().
				Button(r.Choices[0].Choice, []byte(r.Choices[0].Id), updateChoices).
				Button(r.Choices[1].Choice, []byte(r.Choices[1].Id), updateChoices).
				Button(r.Choices[2].Choice, []byte(r.Choices[2].Id), updateChoices).
				Row().
				Button(r.Choices[3].Choice, []byte(r.Choices[3].Id), updateChoices).
				Button(r.Choices[4].Choice, []byte(r.Choices[4].Id), updateChoices).
				Button(r.Choices[5].Choice, []byte(r.Choices[5].Id), updateChoices).
				Row().
				Button(r.Choices[6].Choice, []byte(r.Choices[6].Id), updateChoices).
				Button(r.Choices[7].Choice, []byte(r.Choices[7].Id), updateChoices).
				Button(r.Choices[8].Choice, []byte(r.Choices[8].Id), updateChoices).
				Row().
				Button(r.Choices[9].Choice, []byte(r.Choices[9].Id), updateChoices).
				Button(r.Choices[10].Choice, []byte(r.Choices[10].Id), updateChoices).
				Row().
				Button("Мульти выбор", []byte(r.Id), updateChoices)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        r.Name + ":  ",
				ReplyMarkup: kb,
			})
		case 12:
			kb := inline.New(b).
				Row().
				Button(r.Choices[0].Choice, []byte(r.Choices[0].Id), updateChoices).
				Button(r.Choices[1].Choice, []byte(r.Choices[1].Id), updateChoices).
				Button(r.Choices[2].Choice, []byte(r.Choices[2].Id), updateChoices).
				Row().
				Button(r.Choices[3].Choice, []byte(r.Choices[3].Id), updateChoices).
				Button(r.Choices[4].Choice, []byte(r.Choices[4].Id), updateChoices).
				Button(r.Choices[5].Choice, []byte(r.Choices[5].Id), updateChoices).
				Row().
				Button(r.Choices[6].Choice, []byte(r.Choices[6].Id), updateChoices).
				Button(r.Choices[7].Choice, []byte(r.Choices[7].Id), updateChoices).
				Button(r.Choices[8].Choice, []byte(r.Choices[8].Id), updateChoices).
				Row().
				Button(r.Choices[9].Choice, []byte(r.Choices[9].Id), updateChoices).
				Button(r.Choices[10].Choice, []byte(r.Choices[10].Id), updateChoices).
				Button(r.Choices[11].Choice, []byte(r.Choices[11].Id), updateChoices).
				Row().
				Button("Мульти выбор", []byte(r.Id), updateChoices)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        r.Name + ":  ",
				ReplyMarkup: kb,
			})
		}

	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Итого: %v %s ", len(labelguid), "вопросов"),
	})
}
func updateChoices(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "You selected: " + string(data),
	})
	log.Println(string(data))
}
