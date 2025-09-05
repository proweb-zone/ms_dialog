package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"ms_dialog/internal/app/entity"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDb struct {
	Client *redis.Client
}

type Message struct {
	UserID    int
	Message   string
	Timestamp string
}

func InitRedisDb() (*RedisDb, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123123vv",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisDb{
		Client: client,
	}, nil
}

func (r *RedisDb) AddMsg(dialog *entity.Dialog) error {
	userID := dialog.User_id_sender
	userIDFriend := dialog.User_id_recipient
	msg := dialog.Msg

	timeNow := time.Now().Unix()
	key := fmt.Sprintf("user-%d:friend-%d:dialog", userID, userIDFriend)
	ctx := context.Background()

	prepairMsg := redis.Z{
		Score:  float64(timeNow),
		Member: msg,
	}

	_, err := r.Client.ZAdd(ctx, key, prepairMsg).Result()
	if err != nil {
		return fmt.Errorf("failed to add message: %w", err)
	}

	// Удаляем сообщения, если их больше 1000
	r.Client.ZRemRangeByRank(ctx, key, 0, -1001)

	return nil
}

func (r *RedisDb) GetMessages(userId int, userIdFriend int) (*[]*entity.Dialog, error) {
	buildDialogs := make([]*entity.Dialog, 0)
	ctx := context.Background()
	key := fmt.Sprintf("user-%d:friend-%d:dialog", userId, userIdFriend)

	// Получаем последние 1000 сообщений
	dialogsRedis, err := r.Client.ZRevRangeWithScores(ctx, key, 0, 999).Result()
	if err != nil {
		return nil, err
	}

	for itemId, dialog := range dialogsRedis {
		//timestamp := time.Now().Unix()
		buildDialogs = append(buildDialogs, &entity.Dialog{ID: itemId, User_id_sender: userId, State: true, Msg: dialog.Member.(string)})
	}

	return &buildDialogs, nil
}

func (r *RedisDb) UdfGetMessages(userId int, userIdFriend int) (*[]entity.Dialog, error) {
	ctx := context.Background()

	dialogList, err := r.Client.FCall(ctx, "get_messages",
		[]string{fmt.Sprintf("user-%d", userId), fmt.Sprintf("friend-%d", userIdFriend)}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	var buildDialogs []entity.Dialog
	if err := json.Unmarshal([]byte(dialogList.(string)), &buildDialogs); err != nil {
		return nil, fmt.Errorf("failed to parse messages: %w", err)
	}

	reBuildDialogs := make([]entity.Dialog, 0)
	for itemId, dialog := range buildDialogs {
		reBuildDialogs = append(reBuildDialogs, entity.Dialog{ID: itemId, User_id_sender: userId, User_id_recipient: userIdFriend, State: true, Msg: dialog.Msg})
	}

	return &reBuildDialogs, nil
}
