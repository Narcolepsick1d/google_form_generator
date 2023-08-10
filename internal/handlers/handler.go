package handlers

import (
	"context"
	"encoding/json"
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
	for _, r := range resp {
		lena := len(r.Choices)
		if lena > 12 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Вопрос с 1 или больше 12 вариантов ответа не принимаются",
			})
			return
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Ваши вопросы \n Пожалуйста выберите какие из них имеют мульти-выбор(multiple choices)",
	})
	for _, r := range resp {
		req, err := json.Marshal(r)
		if err != nil {
			log.Println(err, "error while unmarshal")
			return
		}
		kb := inline.New(b).
			Row().
			Button("Мульти выбор", req, h.updateChoicesTrue).
			Button("Один вариает ответа", req, h.updateChoicesFalse).
			Row()
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        r.Name + ":",
			ReplyMarkup: kb,
		})
	}
	kb := inline.New(b).
		Row().
		Button("ВЕРНО", []byte("ok"), supportHandler).
		Button("Неверно", []byte(""), supportHandler)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        fmt.Sprintf("Итого: %v %s ", len(labelguid), "вопросов \n"),
		ReplyMarkup: kb,
	})
}
func supportHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	if string(data) == "ok" {
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "Для помощи напишите @real_eye. Укажите причину проблемы",
	})
}
func (h *H) updateChoicesTrue(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	var r model.RespQuestion
	err := json.Unmarshal(data, &r)
	if err != nil {
		log.Println(err, "error while unmarshal")
		return
	}

	err = h.Label.Update(context.Background(), model.UpdateLabel{Id: r.Id, IsMulti: true})
	if err != nil {
		log.Println(err)
		return
	}
	switch len(r.Choices) {
	case 2:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":",
			ReplyMarkup: kb,
		})
	case 3:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Button(r.Choices[2].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":",
			ReplyMarkup: kb,
		})
	case 4:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[2].Choice, data, h.updateProb).
			Button(r.Choices[3].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":",
			ReplyMarkup: kb,
		})
	case 5:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[2].Choice, data, h.updateProb).
			Button(r.Choices[3].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})
	case 6:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Button(r.Choices[2].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[3].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb).
			Button(r.Choices[5].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})
	case 7:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[2].Choice, data, h.updateProb).
			Button(r.Choices[3].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[4].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb).
			Button(r.Choices[5].Choice, data, h.updateProb).
			Button(r.Choices[6].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})
	case 8:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[2].Choice, data, h.updateProb).
			Button(r.Choices[3].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[5].Choice, data, h.updateProb).
			Button(r.Choices[6].Choice, data, h.updateProb).
			Button(r.Choices[7].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})
	case 9:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Button(r.Choices[2].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[3].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb).
			Button(r.Choices[5].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[6].Choice, data, h.updateProb).
			Button(r.Choices[7].Choice, data, h.updateProb).
			Button(r.Choices[8].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})
	case 10:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Button(r.Choices[2].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[3].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb).
			Button(r.Choices[5].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[6].Choice, data, h.updateProb).
			Button(r.Choices[7].Choice, data, h.updateProb).
			Button(r.Choices[8].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[9].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})

	case 11:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Button(r.Choices[2].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[3].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb).
			Button(r.Choices[5].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[6].Choice, data, h.updateProb).
			Button(r.Choices[7].Choice, data, h.updateProb).
			Button(r.Choices[8].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[9].Choice, data, h.updateProb).
			Button(r.Choices[10].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})
	case 12:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Button(r.Choices[2].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[3].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb).
			Button(r.Choices[5].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[6].Choice, data, h.updateProb).
			Button(r.Choices[7].Choice, data, h.updateProb).
			Button(r.Choices[8].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[9].Choice, data, h.updateProb).
			Button(r.Choices[10].Choice, data, h.updateProb).
			Button(r.Choices[11].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})
	}

}
func (h *H) updateChoicesFalse(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	var r model.RespQuestion
	err := json.Unmarshal(data, &r)
	if err != nil {
		log.Println(err, "error while unmarshal")
		return
	}
	switch len(r.Choices) {
	case 2:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":",
			ReplyMarkup: kb,
		})
	case 3:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Button(r.Choices[2].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":",
			ReplyMarkup: kb,
		})
	case 4:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[2].Choice, data, h.updateProb).
			Button(r.Choices[3].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":",
			ReplyMarkup: kb,
		})
	case 5:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[2].Choice, data, h.updateProb).
			Button(r.Choices[3].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})
	case 6:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Button(r.Choices[2].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[3].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb).
			Button(r.Choices[5].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})
	case 7:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[2].Choice, data, h.updateProb).
			Button(r.Choices[3].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[4].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb).
			Button(r.Choices[5].Choice, data, h.updateProb).
			Button(r.Choices[6].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})
	case 8:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[2].Choice, data, h.updateProb).
			Button(r.Choices[3].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[5].Choice, data, h.updateProb).
			Button(r.Choices[6].Choice, data, h.updateProb).
			Button(r.Choices[7].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})
	case 9:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Button(r.Choices[2].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[3].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb).
			Button(r.Choices[5].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[6].Choice, data, h.updateProb).
			Button(r.Choices[7].Choice, data, h.updateProb).
			Button(r.Choices[8].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})
	case 10:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Button(r.Choices[2].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[3].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb).
			Button(r.Choices[5].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[6].Choice, data, h.updateProb).
			Button(r.Choices[7].Choice, data, h.updateProb).
			Button(r.Choices[8].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[9].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})

	case 11:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Button(r.Choices[2].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[3].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb).
			Button(r.Choices[5].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[6].Choice, data, h.updateProb).
			Button(r.Choices[7].Choice, data, h.updateProb).
			Button(r.Choices[8].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[9].Choice, data, h.updateProb).
			Button(r.Choices[10].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})
	case 12:
		kb := inline.New(b).
			Row().
			Button(r.Choices[0].Choice, data, h.updateProb).
			Button(r.Choices[1].Choice, data, h.updateProb).
			Button(r.Choices[2].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[3].Choice, data, h.updateProb).
			Button(r.Choices[4].Choice, data, h.updateProb).
			Button(r.Choices[5].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[6].Choice, data, h.updateProb).
			Button(r.Choices[7].Choice, data, h.updateProb).
			Button(r.Choices[8].Choice, data, h.updateProb).
			Row().
			Button(r.Choices[9].Choice, data, h.updateProb).
			Button(r.Choices[10].Choice, data, h.updateProb).
			Button(r.Choices[11].Choice, data, h.updateProb)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Chat.ID,
			Text:        r.Name + ":  ",
			ReplyMarkup: kb,
		})
	}
}

func (h *H) updateProb(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	fmt.Println(string(data))
}
