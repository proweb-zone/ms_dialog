package repository

import (
	"database/sql"
	"fmt"
	"ms_dialog/internal/app/entity"
)

type AuthRepository struct {
	conn *sql.DB
}

func NewAuthRepository(conn *sql.DB) *AuthRepository {
	return &AuthRepository{conn}
}

func (a *AuthRepository) CreateAuth(userId int, token string) error {
	query := `INSERT INTO auth (user_id, token) VALUES ($1, $2)`

	_, err := a.conn.Exec(query, userId, token)
	if err != nil {
		return fmt.Errorf("Error create token in DB")
	}

	return nil
}

func (a *AuthRepository) CheckToken(token string) (*entity.Auth, error) {
	row := a.conn.QueryRow("SELECT id, user_id, token, created_at FROM auth WHERE token = $1", token)

	var auth entity.Auth
	err := row.Scan(&auth.ID, &auth.User_id, &auth.Token, &auth.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &auth, nil
}
