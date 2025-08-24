package entity

import "time"

type Auth struct {
	ID        uint      `json:"id"`
	User_id   int       `json:"user_id`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Auth) TableName() string {
	return "auth"
}
