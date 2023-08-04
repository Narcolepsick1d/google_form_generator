package CRUD

import (
	"context"
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
