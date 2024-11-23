package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserEntity(t *testing.T) {
	ue := UserEntity{}
	assert.Nil(t, ue.BirthDay)
	assert.Equal(t, uint(0), ue.ID)
	assert.Empty(t, ue.Name)
	assert.Empty(t, ue.Password)
	assert.Empty(t, ue.Gender)
	//
	assert.Nil(t, ue.CreatedAt)
	assert.Nil(t, ue.UpdatedAt)
}
