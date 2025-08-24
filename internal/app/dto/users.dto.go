package dto

type UsersRequestDto struct {
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Birth_date string `json:"birth_date"`
	Gender     string `json:"gender"`
	Hobby      string `json:"hobby"`
	City       string `json:"city"`
}

type UsersResponseDto struct {
	User_id uint `json:"user_id"`
}

func (d *UsersRequestDto) Validate() error {
	tagMsg := map[string]string{
		"required": "обязательно для заполнения",
		"numeric":  "должно быть числом",
	}

	return validate(d, tagMsg).Error
}
