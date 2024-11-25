package entity

import (
	"time"
)

// UserEntity schema
type UserEntity struct {
	ID        uint
	Name      string
	Email     string
	Password  string
	BirthDay  *time.Time
	Gender    string
	PhotoURL  string `db:"photo_url"`
	Active    bool
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
