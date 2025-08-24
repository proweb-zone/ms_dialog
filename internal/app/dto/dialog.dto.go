package dto

import "time"

type DialogRequestDto struct {
	ID                int       `json:"id"`
	User_id_sender    int       `json:"user_id_sender"`
	User_id_recipient int       `json:"user_id_recipient"`
	Msg               string    `json:"msg"`
	CreatedAt         time.Time `json:"created_at"`
	Updated_at        string    `json:"updated_at"`
}

type DialogResponseDto struct {
	Dialog_id int `json:"post_id"`
}

func (d *DialogRequestDto) Validate() error {
	tagMsg := map[string]string{
		"required": "обязательно для заполнения",
		"numeric":  "должно быть числом",
	}

	return validate(d, tagMsg).Error
}
