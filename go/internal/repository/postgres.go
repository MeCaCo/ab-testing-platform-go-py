package repository

import (
	"ab-testing-platform-go/internal/model"
	"context"
)

// PostgresRepository — заглушка
type PostgresRepository struct{}

func NewPostgresRepository(host string, port int, user, password, dbName string) *PostgresRepository {
	// TODO: реализовать подключение к PostgreSQL
	return &PostgresRepository{}
}

// SaveEvent — заглушка
func (p *PostgresRepository) SaveEvent(ctx context.Context, event *model.Event) error {
	// TODO: реализовать сохранение события в PostgreSQL
	return nil
}

// BatchSaveEvents — заглушка
func (p *PostgresRepository) BatchSaveEvents(ctx context.Context, events []*model.Event) error {
	// TODO: реализовать батчевое сохранение
	return nil
}
