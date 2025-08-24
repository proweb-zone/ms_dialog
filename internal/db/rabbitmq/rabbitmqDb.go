package rabbitmq

import (
	"context"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMq struct {
	Conn *amqp.Connection
}

func InitRabbitMqDb() (*RabbitMq, error) {
	conn, err := amqp.Dial("amqp://guest:guest@172.33.0.6:5672/")
	if err != nil {
		return nil, err
	}

	return &RabbitMq{
		Conn: conn,
	}, nil
}

func (r *RabbitMq) AddMsg(userId int, message string) error {
	ch, err := r.Conn.Channel()
	if err != nil {
		return err
	}

	userIdStr := "userid_" + strconv.FormatInt(int64(userId), 10)

	// Определение очереди
	q, errQueue := ch.QueueDeclare(
		userIdStr, // Имя очереди
		false,     // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // args
	)
	if errQueue != nil {
		return errQueue
	}

	// Создание сообщения
	msg := []byte(message)

	// Отправка сообщения
	errQueue = ch.PublishWithContext(context.Background(),
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	if errQueue != nil {
		return errQueue
	}

	return nil
}
