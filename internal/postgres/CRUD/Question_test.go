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
func TestQuestionRepo_Get(t *testing.T) {
	repo := NewQuestion(db)
	resp, err := repo.Get(context.Background(), "f8031f95-a316-46ea-8d44-38a920053dcd")
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	log.Print("len", len(resp), resp)

}
