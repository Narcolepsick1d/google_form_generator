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
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
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

var (
	Proba []model.StateProb
	Next  bool
)

func (h *H) urlStartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	if helper.IsGoogleFormsLink(update.Message.Text) {
		time.Sleep(1 * time.Second)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Проверьте вашу ссылку  ",
		})
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
		ctx, cancel := context.WithCancel(context.Background())
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			select {
			case <-signals:
				fmt.Println("Отмена операции по сигналу")
				cancel()
			}
		}()
		for _, r := range resp {
			SetNext(false)
			select {
			case <-ctx.Done():
				fmt.Println("Операция отменена")
				return
			default:
				tr := model.PropCh{
					IsMulti:   true,
					RespQuest: r,
				}
				reqT, err := json.Marshal(tr)
				if err != nil {
					log.Println(err, "error while unmarshal")
					return
				}
				fal := model.PropCh{
					IsMulti:   false,
					RespQuest: r,
				}
				reqF, err := json.Marshal(fal)
				if err != nil {
					log.Println(err, "error while unmarshal")
					return
				}
				kb := inline.New(b).
					Row().
					Button("Мульти выбор", reqT, h.updateChoices).
					Button("Один вариает ответа", reqF, h.updateChoices).
					Row()
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:      update.Message.Chat.ID,
					Text:        r.Name + ":",
					ReplyMarkup: kb,
				})
				if !GetNext() {
					fmt.Println("Ожидание подтверждения от пользователя...")
					// Здесь ожидание, пока пользователь не подтвердит или отменит
					// вы можете вызвать функцию, которая будет ожидать ввода пользователя и возвращать true/false
				}
			}

		}
		chswithProb, err := h.Question.Get(ctx, qId)
		if err != nil {
			log.Print("1", err)
			return
		}
		for _, m := range chswithProb {
			gg := make([]string, 0, 10)
			for _, p := range m.Choices {
				str := p.Choice + " " + strconv.Itoa(p.Probability) + "%"
				gg = append(gg, str)
			}
			yesno := ""
			if m.IsMulti {
				yesno = "есть мультвыбор"
			} else {
				yesno = "один вариант ответа"
			}
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   m.Name + "\n" + strings.Join(gg, "; ") + "\n " + yesno,
			})
		}
		marshal, err := json.Marshal(chswithProb)
		if err != nil {
			return
		}
		kb := inline.New(b).
			Row().
			Button("ВЕРНО", marshal, h.AmountCather).
			Button("Неверно", []byte(""), supportHandler)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        fmt.Sprintf("Итого: %v %s ", len(labelguid), "вопросов \n"),
			ReplyMarkup: kb,
		})

	} else {
		if len(Proba) == 0 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("Процесс прошел не по плану. Пожалуйста пройдите заново :("),
			})
			return
		}
		if !helper.IsProb(update.Message.Text) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("вы не правильно ввели проценты. Пример: \n 10,20,30 \n для соответсвуйщих вариантов ответа"),
			})
			return
		}
		numberStrings := strings.Split(update.Message.Text, ",")
		numbers := make([]int, 0, 10)
		sum := 0
		var bcc error
		for _, s := range numberStrings {
			num, err := strconv.Atoi(s)
			log.Println(err)
			bcc = err
			sum += num
		}
		fmt.Println(sum)
		if bcc != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("вы не правильно ввели проценты. Пример: \n 10,20,30 \n для соответсвуйщих вариантов ответа"),
			})
			return
		}
		if sum >= 100 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("Сумма процентов выбора перевалила за 100 генерация данных может стать не точными"),
			})
		}
		for _, numStr := range numberStrings {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				log.Println(err)
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   fmt.Sprintf("вы не правильно ввели проценты. Пример: \n 10,20,30 \n для соответсвуйщих вариантов ответа"),
				})
				return
			}
			numbers = append(numbers, num)
		}
		if len(Proba) != len(numbers) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("Количество вариантов ответа не совпадает с количеством процентов"),
			})
			return
		}
		for i, c := range Proba {
			err := h.Choice.Update(ctx, model.UpdateChoices{Id: c.Id, Probability: numbers[i]})
			if err != nil {
				log.Printf("error while creating user %v", err)
				return
			}
		}
		Proba = make([]model.StateProb, 0, 12)
		SetNext(true)
	}
}
func SetNext(bool2 bool) {
	Next = bool2
}
func GetNext() bool {
	for {
		if Next == true {
			break
		}
	}
	return Next
}
func (h *H) AmountCather(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	kb := inline.New(b).
		Row().
		Button("50", data, h.DodosAttacker50).
		Button("125", data, h.DodosAttacker125).
		Row().
		Button("260", data, h.DodosAttacker260).
		Button("500", data, h.DodosAttacker500)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        "Выберите количество заполнений",
		ReplyMarkup: kb,
	})
}
func (h *H) DodosAttacker50(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	var final []model.RespQuestion
	err := json.Unmarshal(data, &final)
	if err != nil {
		log.Print(err)
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "Процесс начался ждите ...",
	})
	urll := helper.ReplaceUrl(final[0].Url)

	for i := 0; i < 50; i++ {
		err = helper.Dodos(final, urll)
		if err != nil {
			log.Print(err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Chat.ID,
				Text:   "Ошибка. Используйте команду для связи с админом /help",
			})
		}

	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "Процесс завершился",
	})

}
func (h *H) DodosAttacker125(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	var final []model.RespQuestion
	err := json.Unmarshal(data, &final)
	if err != nil {
		log.Print(err)
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "Процесс начался ждите ...",
	})
	urll := helper.ReplaceUrl(final[0].Url)
	for i := 0; i < 125; i++ {
		err = helper.Dodos(final, urll)
		if err != nil {
			log.Print(err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Chat.ID,
				Text:   "Ошибка. Используйте команду для связи с админом /help",
			})
		}

	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "Процесс завершился",
	})

}

func (h *H) DodosAttacker260(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	var final []model.RespQuestion
	err := json.Unmarshal(data, &final)
	if err != nil {
		log.Print(err)
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "Процесс начался ждите ...",
	})
	urll := helper.ReplaceUrl(final[0].Url)

	for i := 0; i < 260; i++ {
		err = helper.Dodos(final, urll)
		if err != nil {
			log.Print(err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Chat.ID,
				Text:   "Ошибка. Используйте команду для связи с админом /help",
			})
		}

	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "Процесс завершился",
	})

}

func (h *H) DodosAttacker500(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	var final []model.RespQuestion
	err := json.Unmarshal(data, &final)
	if err != nil {
		log.Print(err)
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "Процесс начался ждите ...",
	})
	urll := helper.ReplaceUrl(final[0].Url)

	for i := 0; i < 500; i++ {
		err = helper.Dodos(final, urll)
		if err != nil {
			log.Print(err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Chat.ID,
				Text:   "Ошибка. Используйте команду для связи с админом /help",
			})
		}

	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "Процесс завершился",
	})

}

func (h *H) helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Для помощи напишите @real_eye. Укажите причину проблемы",
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
func (h *H) updateChoices(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	var r model.PropCh
	err := json.Unmarshal(data, &r)
	if err != nil {
		log.Println(err, "error while unmarshal")
		return
	}

	err = h.Label.Update(context.Background(), model.UpdateLabel{Id: r.RespQuest.Id, IsMulti: r.IsMulti})
	if err != nil {
		log.Println(err)
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Chat.ID,
		Text:   "Напишите вероятность через запятую соответвенно к вариантам ответа",
	})
	strs := make([]string, 0, 10)
	for _, chc := range r.RespQuest.Choices {
		strs = append(strs, chc.Choice)
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Chat.ID,
		Text:   r.RespQuest.Name + ": \n" + strings.Join(strs, ", "),
	})
	for _, k := range r.RespQuest.Choices {
		Proba = append(Proba, model.StateProb{Id: k.Id})
	}

}

/*
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
		Proba = append(Proba, model.StateProb{Id: r.Choices[0].Id}, model.StateProb{Id: r.Choices[1].Id})
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
		Proba = append(Proba, model.StateProb{Id: r.Choices[0].Id}, model.StateProb{Id: r.Choices[1].Id})
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
		Proba = append(Proba, model.StateProb{Id: r.Choices[0].Id}, model.StateProb{Id: r.Choices[1].Id})
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
		Proba = append(Proba, model.StateProb{Id: r.Choices[0].Id}, model.StateProb{Id: r.Choices[1].Id})
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
		Proba = append(Proba, model.StateProb{Id: r.Choices[0].Id}, model.StateProb{Id: r.Choices[1].Id})
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
		Proba = append(Proba, model.StateProb{Id: r.Choices[0].Id}, model.StateProb{Id: r.Choices[1].Id})
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
		Proba = append(Proba, model.StateProb{Id: r.Choices[0].Id}, model.StateProb{Id: r.Choices[1].Id})
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
		Proba = append(Proba, model.StateProb{Id: r.Choices[0].Id}, model.StateProb{Id: r.Choices[1].Id})
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
		Proba = append(Proba, model.StateProb{Id: r.Choices[0].Id}, model.StateProb{Id: r.Choices[1].Id})
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
		Proba = append(Proba, model.StateProb{Id: r.Choices[0].Id}, model.StateProb{Id: r.Choices[1].Id})
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
		Proba = append(Proba, model.StateProb{Id: r.Choices[0].Id}, model.StateProb{Id: r.Choices[1].Id})
	}
}
*/
