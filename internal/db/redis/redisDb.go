package redis

import (
	"context"
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

func (r *RedisDb) AddMsg(userID int, msg string) error {
	timeNow := time.Now().Unix()
	key := fmt.Sprintf("user:%d:messages", userID)
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

func (r *RedisDb) GetMessages(userIDs []int) ([]*entity.Posts, error) {
	posts := make([]*entity.Posts, 0)
	ctx := context.Background()
	for _, userID := range userIDs {
		key := fmt.Sprintf("user:%d:messages", userID)

		// Получаем последние 1000 сообщений
		postsListRedis, err := r.Client.ZRevRangeWithScores(ctx, key, 0, 999).Result()
		if err != nil {
			return nil, err
		}

		for itemId, postItem := range postsListRedis {
			//timestamp := time.Now().Unix()
			posts = append(posts, &entity.Posts{ID: itemId, User_id: userID, Text: postItem.Member.(string)})
		}

	}
	return posts, nil
}
