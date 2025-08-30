package repository

import (
	"database/sql"
	"fmt"
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

	// row := a.conn.QueryRow("SELECT token FROM auth WHERE token = $1", token)

	// errToken := row.Scan(&auth.Token)
	// if errToken != nil {
	// 	return nil, fmt.Errorf("Error get token in DB")
	// }

	// return &auth, nil
}
