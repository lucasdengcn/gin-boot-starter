package entity

import (
	"time"
)

// UserEntity schema
type UserEntity struct {
	ID        uint
	Name      string
	Email     string
	Password  string `db:"hashed_password"`
	BirthDay  *time.Time
	Gender    string
	PhotoURL  string `db:"photo_url"`
	Active    bool
	Roles     string
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
