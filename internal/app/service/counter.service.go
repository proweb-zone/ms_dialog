package service

import (
	"context"
	"fmt"
	"sync"

	"ms_dialog/internal/app/entity"
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

// Increment увеличивает счетчик непрочитанных сообщений
func (s *CounterService) IncrementCounter(dialog *entity.Dialog) error {

	userID := dialog.User_id_sender
	friendID := dialog.User_id_recipient
	ctx := context.Background()

	key := s.getKey(userID, friendID)

	count, err := s.repo.IncrementCounter(ctx, key)
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
// func (s *CounterService) DecrementCounter(ctx context.Context, userID, chatID string) error {
// 	key := s.getKey(userID, chatID)

// 	count, err := s.repo.DecrementCounter(ctx, key)
// 	if err != nil {
// 		return err
// 	}

// 	// Не позволяем счетчику уйти в отрицательные значения
// 	if count < 0 {
// 		count = 0
// 		s.repo.SetCounter(ctx, key, 0, 0)
// 	}

// 	s.cacheMutex.Lock()
// 	s.cache[key] = count
// 	s.cacheMutex.Unlock()

// 	return nil
// }

// // Get получает значение счетчика
func (s *CounterService) GetCounter(userID, friendID int) (int64, error) {
	key := s.getKey(userID, friendID)
	ctx := context.Background()

	// Проверяем кэш сначала
	s.cacheMutex.RLock()
	if count, exists := s.cache[key]; exists {
		s.cacheMutex.RUnlock()
		return count, nil
	}
	s.cacheMutex.RUnlock()

	// Если нет в кэше, идем в Redis
	count, err := s.repo.GetCounter(ctx, key)
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

// // Reset сбрасывает счетчик
func (s *CounterService) ResetCounter(userID int, friendID int) error {
	key := s.getKey(userID, friendID)
	ctx := context.Background()

	err := s.repo.DelCounter(ctx, key)
	if err != nil {
		return err
	}

	// Удаляем из кэша
	s.cacheMutex.Lock()
	delete(s.cache, key)
	s.cacheMutex.Unlock()

	return nil
}

func (s *CounterService) SetCounter(userID int, friendID int, counter int64) {
	ctx := context.Background()
	key := s.getKey(userID, friendID)

	s.repo.SetCounter(ctx, key, counter, 0)
}

func (s *CounterService) getKey(userID int, friendID int) string {
	return fmt.Sprintf("user-id-%d:friend-id-%d:counter", userID, friendID)
}
