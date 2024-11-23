package models

import "time"

// UserSignUp request input model
type UserSignUp struct {
	Name     string     `json:"name" binding:"required"`
	BirthDay *time.Time `json:"birthday" time_format:"2006-01-02" time_utc:"1" binding:"required"`
	Gender   string     `json:"gender" binding:"required"`
	PhotoURL string     `json:"photo_url" binding:"required"`
}

// UserInfo response output model
type UserInfo struct {
	ID       uint       `json:"id"`
	Name     string     `json:"name"`
	BirthDay *time.Time `json:"birthday" time_format:"2006-01-02"`
	Gender   string     `json:"gender"`
	PhotoURL string     `json:"photo_url"`
}
