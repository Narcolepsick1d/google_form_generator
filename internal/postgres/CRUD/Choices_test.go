package CRUD

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google-gen/internal/model"
	"testing"
)

func TestChoicesRepo_Update(t *testing.T) {
	repo := NewChoice(db)
	err := repo.Update(context.Background(), model.UpdateChoices{Id: "313a0992-b03c-463f-9c32-1db294d35f40", Probability: 100})
	assert.NoError(t, err)
}
