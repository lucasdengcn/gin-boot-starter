package entity

import (
	"database/sql"
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
	PhotoURL  *sql.NullString `db:"photo_url"`
	Active    bool
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
