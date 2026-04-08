package service

import (
	"ab-testing-platform-go/internal/model"
	"context"
)

// BatcherService отвечает за микро-батчинг событий (собирает события БД пачками)
type BatcherService struct {
	// TODO: реализовать каналы и горутины для батчинга
}

func NewBatcherService() *BatcherService {
	return &BatcherService{}
}

// AddEvent добавляет событие в очередь на отправку
func (b *BatcherService) AddEvent(ctx context.Context, event *model.Event) error {
	// TODO: реализовать отправку в канал
	return nil
}

// Start запускает фоновую горутину для батчинга
func (b *BatcherService) Start() {
	// TODO: реализовать
}
