package models

import "time"

type AuthTokens struct {
	AccessToken       string     `json:"access_token"`
	RefreshToken      string     `json:"refresh_token"`
	ExpireIn          *time.Time `json:"at_expire_in" time_format:"2006-01-01 17:23:10"`
	RefreshTokenExpIn *time.Time `json:"rt_expire_in" time_format:"2006-01-01 17:23:10"`
}
