package CRUD

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestQuestionRepo_GetByUrl(t *testing.T) {
	repo := NewQuestion(db)
	resp, err := repo.GetByUrl(context.Background(), "https://docs.google.com/forms/d/1tP_rq1HHEksobsVTkk5y7ywMq-_uyjQGwbF2NkFB7A8/viewform?edit_requested=true")
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	log.Print(resp)
}
func TestQuestionRepo_Get(t *testing.T) {
	repo := NewQuestion(db)
	resp, err := repo.Get(context.Background(), "aa2506aa-dd1f-4373-9b2f-9cc9ce120f06")
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	log.Print("len", len(resp), resp)
	marshal, err := json.Marshal(resp)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
}
