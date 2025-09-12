package service

import (
	"context"
	"fmt"
	"sync"

	"ms_dialog/internal/app/repository"
)

type CounterService struct {
	repo       *repository.CounterRepository
	cache      map[string]int64
	cacheMutex sync.RWMutex
}

type Counter struct {
	UserID    string `json:"user_id"`
	ChatID    string `json:"chat_id,omitempty"`
	Count     int64  `json:"count"`
	UpdatedAt int64  `json:"updated_at"`
}

func NewCounterService(newRepo *repository.CounterRepository) *CounterService {
	return &CounterService{
		repo:  newRepo,
		cache: make(map[string]int64),
	}
}

func (s *CounterService) getKey(userID, chatID string) string {
	if chatID == "" {
		return fmt.Sprintf("unread:user:%s", userID)
	}
	return fmt.Sprintf("unread:user:%s:chat:%s", userID, chatID)
}

// Increment увеличивает счетчик непрочитанных сообщений
func (s *CounterService) Increment(ctx context.Context, userID, chatID string) error {
	key := s.getKey(userID, chatID)

	count, err := s.repo.IncrementRedis(ctx, key)
	if err != nil {
		return err
	}

	// Обновление кэша
	s.cacheMutex.Lock()
	s.cache[key] = count
	s.cacheMutex.Unlock()

	return nil
}

// Decrement уменьшает счетчик
func (s *CounterService) Decrement(ctx context.Context, userID, chatID string) error {
	key := s.getKey(userID, chatID)

	count, err := s.repo.Decrement(ctx, key)
	if err != nil {
		return err
	}

	// Не позволяем счетчику уйти в отрицательные значения
	if count < 0 {
		count = 0
		s.repo.Set(ctx, key, 0, 0)
	}

	s.cacheMutex.Lock()
	s.cache[key] = count
	s.cacheMutex.Unlock()

	return nil
}

// Get получает значение счетчика
func (s *CounterService) Get(ctx context.Context, userID, chatID string) (int64, error) {
	key := s.getKey(userID, chatID)

	// Проверяем кэш сначала
	s.cacheMutex.RLock()
	if count, exists := s.cache[key]; exists {
		s.cacheMutex.RUnlock()
		return count, nil
	}
	s.cacheMutex.RUnlock()

	// Если нет в кэше, идем в Redis
	count, err := s.repo.Get(ctx, key)
	if err != nil {
		// Ключ не существует, возвращаем 0
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	// Обновляем кэш
	s.cacheMutex.Lock()
	s.cache[key] = count
	s.cacheMutex.Unlock()

	return count, nil
}

// Reset сбрасывает счетчик
func (s *CounterService) Reset(ctx context.Context, userID, chatID string) error {
	key := s.getKey(userID, chatID)

	err := s.repo.Del(ctx, key)
	if err != nil {
		return err
	}

	// Удаляем из кэша
	s.cacheMutex.Lock()
	delete(s.cache, key)
	s.cacheMutex.Unlock()

	return nil
}
