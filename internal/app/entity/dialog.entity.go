package entity

import "time"

type Dialog struct {
	ID                int       `json:"id"`
	User_id_sender    int       `json:"user_id_sender"`
	User_id_recipient int       `json:"user_id"`
	Msg               string    `json:"msg"`
	CreatedAt         time.Time `json:"created_at"`
	Updated_at        string    `json:"updated_at"`
}

func (Dialog) TableName() string {
	return "dialog"
}
