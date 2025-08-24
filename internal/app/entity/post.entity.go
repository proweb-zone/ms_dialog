package entity

import "time"

type Posts struct {
	ID        int       `json:"id"`
	User_id   int       `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Posts) TableName() string {
	return "posts"
}
