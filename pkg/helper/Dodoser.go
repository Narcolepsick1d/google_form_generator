package helper

import (
	"errors"
	"google-gen/internal/model"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func Dodos(info []model.RespQuestion, urll string) error {
	resp := make([]model.FinalEntity, 0, 75)
	for _, j := range info {
		if j.IsMulti {
			totalProbability := 0

			for _, choice := range j.Choices {
				totalProbability += choice.Probability
			}

			selectedItems := make([]string, 0)

			randomCount := rand.Intn(len(j.Choices)) + 1

			// Создаем карту для отслеживания уже выбранных строк
			selectedMap := make(map[string]bool)

			for i := 0; i < randomCount; {
				randomNumber := rand.Intn(totalProbability) + 1

				cumulativeProb := 0

				for _, choice := range j.Choices {
					cumulativeProb += choice.Probability
					if randomNumber <= cumulativeProb {
						if !selectedMap[choice.Choice] {
							selectedItems = append(selectedItems, choice.Choice)
							selectedMap[choice.Choice] = true
							i++
						}
						break
					}
				}
			}
			resp = append(resp, model.FinalEntity{EntryId: j.EntryId, Choice: selectedItems})
		} else {
			totalProbability := 0

			for _, choice := range j.Choices {
				totalProbability += choice.Probability
			}

			randomNumber := rand.Intn(totalProbability) + 1

			cumulativeProb := 0

			for _, choice := range j.Choices {
				cumulativeProb += choice.Probability
				if randomNumber <= cumulativeProb {
					chs := make([]string, 0)
					chs = append(chs, choice.Choice)
					resp = append(resp, model.FinalEntity{EntryId: j.EntryId, Choice: chs})
					break
				}
			}
		}
	}
	nigger := make(map[string][]string)
	for _, k := range resp {
		nigger["entry."+k.EntryId] = k.Choice
	}
	time.Sleep(2 * time.Second)
	resp = make([]model.FinalEntity, 0, 75)
	t, err := http.PostForm(urll, nigger)
	if err != nil {
		log.Print(err)
		return err
	}
	if t.StatusCode != 200 {
		return errors.New("Что то не так с опросником")
	}
	return nil
}
func ReplaceUrl(url string) string {
	index := strings.Index(url, "/viewform")
	n := url[:index] + "/formResponse"
	return n
}
