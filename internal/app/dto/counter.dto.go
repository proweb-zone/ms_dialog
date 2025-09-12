package dto

type CounterRequest struct {
	UserID string `json:"user_id"`
	ChatID string `json:"chat_id,omitempty"`
}

type CounterResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
