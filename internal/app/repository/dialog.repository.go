package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"ms_dialog/internal/app/entity"

	"github.com/lib/pq"
)

type DialogRepository struct {
	conn *sql.DB
}

func NewDialogRepository(conn *sql.DB) *DialogRepository {
	return &DialogRepository{conn}
}

func (d *DialogRepository) SendMsgUser(newMsg *entity.Dialog) (*entity.Dialog, error) {
	ctx := context.Background()

	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := d.conn.PrepareContext(ctx, `INSERT INTO dialog (user_id_sender, user_id_recipient, msg, state) VALUES ($1, $2, $3, $4) RETURNING id`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, newMsg.User_id_sender, newMsg.User_id_recipient, newMsg.Msg, newMsg.State)

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var dialog entity.Dialog
	row.Scan(&dialog.ID)
	row = d.conn.QueryRow("SELECT id, user_id_sender, user_id_recipient, msg, state, created_at, updated_at FROM dialog WHERE id = $1", &dialog.ID)
	err = row.Scan(&dialog.ID, &dialog.User_id_sender, &dialog.User_id_recipient, &dialog.Msg, &dialog.State, &dialog.CreatedAt, &dialog.Updated_at)
	if err != nil {
		return nil, fmt.Errorf("scanning result: %w", err)
	}

	return &dialog, nil
}

func (d *DialogRepository) GetDialogList(userIdSender int, userIdRecepient int) (*[]entity.Dialog, error) {

	var idsSender []int
	idsSender = append(idsSender, userIdSender)
	idsSender = append(idsSender, userIdRecepient)

	idsSenderValues := make([]interface{}, 1)
	idsSenderValues[0] = pq.Array(idsSender)

	var idsRecipient []int
	idsRecipient = append(idsRecipient, userIdRecepient)
	idsRecipient = append(idsRecipient, userIdSender)

	idsRecipientValues := make([]interface{}, 1)
	idsRecipientValues[0] = pq.Array(idsRecipient)

	ctx := context.Background()

	query := "SELECT id, user_id_sender, user_id_recipient, msg, state, created_at, updated_at FROM dialog WHERE user_id_sender = ANY($1) AND user_id_recipient = ANY($2) AND state = true"

	rows, err := d.conn.QueryContext(ctx, query, idsSenderValues[0], idsRecipientValues[0])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dialogList []entity.Dialog
	for rows.Next() {
		var dialog entity.Dialog
		err := rows.Scan(&dialog.ID, &dialog.User_id_sender, &dialog.User_id_recipient, &dialog.Msg, &dialog.State, &dialog.CreatedAt, &dialog.Updated_at)
		if err != nil {
			return nil, err
		}
		dialogList = append(dialogList, dialog)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &dialogList, nil
}

func (d *DialogRepository) ActiveMsg(id int) error {
	ctx := context.Background()
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer tx.Rollback()

	stmt, err := d.conn.PrepareContext(ctx, `UPDATE dialog SET state = true WHERE id = $1`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	stmt.QueryRowContext(ctx, id)

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (d *DialogRepository) DeleteMsg(id int) error {
	query := "DELETE FROM dialog WHERE id = $1"

	result, err := d.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при получении количества удаленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("запись с ID %d не найдена", id)
	}

	return nil
}
