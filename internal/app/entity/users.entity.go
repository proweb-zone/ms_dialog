package entity

import "time"

type Users struct {
	ID         uint      `json:"id"`
	First_name string    `json:"first_name"`
	Last_name  string    `json:"last_name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Birth_date time.Time `json:"birth_date"`
	Gender     string    `json:"gender"`
	Hobby      string    `json:"hobby"`
	City       string    `json:"city"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Users) TableName() string {
	return "users"
}
