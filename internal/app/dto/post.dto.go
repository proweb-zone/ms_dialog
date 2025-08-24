package dto

type PostRequestDto struct {
	ID         int    `json:"id"`
	User_id    int    `json:"user_id"`
	Text       string `json:"text"`
	Updated_at string `json:"updated_at"`
}

type PostResponseDto struct {
	Post_id int `json:"post_id"`
}

func (d *PostRequestDto) Validate() error {
	tagMsg := map[string]string{
		"required": "обязательно для заполнения",
		"numeric":  "должно быть числом",
	}

	return validate(d, tagMsg).Error
}
