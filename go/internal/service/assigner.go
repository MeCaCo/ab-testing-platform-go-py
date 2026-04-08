// go/internal/service/assigner.go
package service

import (
	"ab-testing-platform-go/internal/model"
	"ab-testing-platform-go/internal/repository"
	"ab-testing-platform-go/pkg/logger"
	"context"
	"fmt"
	"time"
)

type AssignerService struct {
	redisRepo *repository.RedisRepository
	dbRepo    *repository.PostgresRepository
	logger    *logger.SimpleLogger
}

// ПРАВИЛЬНАЯ версия: 3 параметра с типами
func NewAssignerService(redisRepo *repository.RedisRepository, dbRepo *repository.PostgresRepository, log *logger.SimpleLogger) *AssignerService {
	return &AssignerService{
		redisRepo: redisRepo,
		dbRepo:    dbRepo,
		logger:    log,
	}
}

// AssignVariant определяет, в какую группу (A или B) попадёт пользователь
// ВНИМАНИЕ: вызываем AssignVariant из пакета service!
func (s *AssignerService) AssignVariant(ctx context.Context, testID, userID string) (string, error) {
	cacheKey := fmt.Sprintf("assignment:%s:%s", testID, userID)

	// 1. Проверяем кэш (Redis)
	cachedVariant, err := s.redisRepo.GetAssignment(ctx, cacheKey)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to get assignment from redis: %v", err))
		return "", fmt.Errorf("failed to get assignment from redis: %w", err)
	}
	if cachedVariant != "" {
		s.logger.Info(fmt.Sprintf("Cache hit for %s", cacheKey))
		return cachedVariant, nil
	}

	// 2. Если не в кэше, используем хэширование ИЗ ТОГО ЖЕ ПАКЕТА service
	variant := AssignVariant(testID, userID) // <- ВАЖНО: вызываем функцию из пакета service

	// 3. Сохраняем в кэш на 24 часа
	err = s.redisRepo.SetAssignment(ctx, cacheKey, variant, 24*time.Hour)
	if err != nil {
		s.logger.Warn(fmt.Sprintf("Could not cache assignment: %v", err))
		// Не фатально, продолжаем
	}

	s.logger.Info(fmt.Sprintf("Assigned variant '%s' for %s", variant, cacheKey))
	return variant, nil
}

// RecordEvent сохраняет событие (показ, клик, конверсия) в БД
func (s *AssignerService) RecordEvent(ctx context.Context, event *model.Event) error {
	// Пока просто передаём в DB repo (позже реализуем батчер)
	return s.dbRepo.SaveEvent(ctx, event)
}
