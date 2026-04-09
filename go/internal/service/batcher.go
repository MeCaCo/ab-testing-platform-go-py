package service

import (
	"context"
	"log"
	"sync"
	"time"

	"ab-testing-platform-go/internal/model"
	"ab-testing-platform-go/internal/repository"
)

const (
	BatchSize    = 100             // Сколько событий в пачке
	BatchTimeout = 1 * time.Second // Сколько ждать, прежде чем отправить пачку
)

// BatcherService отвечает за микро-батчинг событий (собирает события в БД пачками)
type BatcherService struct {
	dbRepo    *repository.PostgresRepository
	eventChan chan *model.Event  // Канал для получения событий от хендлеров
	wg        sync.WaitGroup     // Для ожидания завершения горутины
	ctx       context.Context    // Контекст для остановки горутины
	cancel    context.CancelFunc // Функция для отмены контекста
}

func NewBatcherService(dbRepo *repository.PostgresRepository) *BatcherService {
	ctx, cancel := context.WithCancel(context.Background())
	bs := &BatcherService{
		dbRepo:    dbRepo,
		eventChan: make(chan *model.Event, 1000), // Буферизованный канал на 1000 событий
		ctx:       ctx,
		cancel:    cancel,
	}

	// Запускаем горутину, которая будет собирать и отправлять пачки
	bs.wg.Add(1)
	go bs.batchLoop()

	return bs
}

// AddEvent добавляет событие в очередь на отправку
func (b *BatcherService) AddEvent(ctx context.Context, event *model.Event) error {
	select {
	case b.eventChan <- event: // Отправляем событие в канал
		return nil
	case <-ctx.Done(): // Если контекст отменили (например, во время запроса)
		return ctx.Err()
	case <-b.ctx.Done(): // Если сервис останавливается
		return b.ctx.Err()
	}
}

// batchLoop — основная горутина, которая собирает и отправляет пачки
func (b *BatcherService) batchLoop() {
	defer b.wg.Done()

	var batch []*model.Event
	ticker := time.NewTicker(BatchTimeout)
	defer ticker.Stop()

	for {
		select {
		case event := <-b.eventChan: // Получаем событие из канала
			batch = append(batch, event)

			// Если пачка набрана, отправляем
			if len(batch) >= BatchSize {
				b.sendBatch(batch)
				batch = nil                // Очищаем пачку
				ticker.Reset(BatchTimeout) // Сбрасываем таймер
			}

		case <-ticker.C: // Если таймер истёк, отправляем оставшиеся события
			if len(batch) > 0 {
				b.sendBatch(batch)
				batch = nil
			}

		case <-b.ctx.Done(): // Если сервис останавливается
			// Отправляем оставшиеся события перед остановкой
			if len(batch) > 0 {
				b.sendBatch(batch)
			}
			return
		}
	}
}

// sendBatch отправляет пачку событий в репозиторий (в БД)
func (b *BatcherService) sendBatch(batch []*model.Event) {
	start := time.Now()
	err := b.dbRepo.BatchSaveEvents(context.Background(), batch)
	if err != nil {
		log.Printf("[BATCHER] Error saving batch of %d events: %v", len(batch), err)
	} else {
		log.Printf("[BATCHER] Successfully saved batch of %d events in %v", len(batch), time.Since(start))
	}
}

// Stop gracefully shuts down the batcher
func (b *BatcherService) Stop() {
	b.cancel()         // Отменяем контекст, чтобы batchLoop завершился
	b.wg.Wait()        // Ждём, пока горутина завершится
	close(b.eventChan) // Закрываем канал
}
