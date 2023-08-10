package CRUD

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google-gen/internal/model"
	"testing"
)

func TestLabelRepo_Update(t *testing.T) {
	repo := NewLabel(db)
	err := repo.Update(context.Background(), model.UpdateLabel{Id: "6d52f3a5-2a65-4002-b039-7c3974c3a1e8",
		IsMulti: true})
	assert.NoError(t, err)
}
