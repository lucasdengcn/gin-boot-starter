package models

import "time"

// UserSignUp request input model
type UserSignUp struct {
	Name     string     `json:"name" binding:"required"`
	BirthDay *time.Time `json:"birthday,omitempty" time_format:"2006-01-02" binding:"required"`
	Gender   string     `json:"gender" binding:"required,gender" example:"male, female"`
	PhotoURL string     `json:"photo_url,omitempty" binding:"required,http_url"`
	Email    string     `json:"email" binding:"required,email"`
	Password string     `json:"password" binding:"required"`
}

// UserInfo response output model
type UserInfo struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	BirthDay  *time.Time `json:"birthday,omitempty" time_format:"2006-01-02"`
	Gender    string     `json:"gender,omitempty" example:"male, female"`
	PhotoURL  string     `json:"photo_url,omitempty"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at,omitempty" time_format:"2006-01-02"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" time_format:"2006-01-02"`
}

// UserInfoUpdate request input model
type UserInfoUpdate struct {
	ID       uint       `json:"id" binding:"required"`
	Name     string     `json:"name" binding:"required"`
	BirthDay *time.Time `json:"birthday" time_format:"2006-01-02" binding:"required"`
	Gender   string     `json:"gender" binding:"required,gender" example:"male, female"`
	PhotoURL string     `json:"photo_url" binding:"required,http_url"`
	Email    string     `json:"email" binding:"required,email"`
	Status   string     `json:"status" binding:"required,enum=on off"`
}
