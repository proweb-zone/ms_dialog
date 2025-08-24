package entity

import "time"

type Friends struct {
	ID        uint      `json:"id"`
	User_id   int       `json:"user_id`
	Friend_id int       `json:"friend_id`
	CreatedAt time.Time `json:"created_at"`
}

func (Friends) TableName() string {
	return "friends"
}
