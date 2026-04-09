// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"ab-testing-platform-go/config"
	"ab-testing-platform-go/internal/handler"
	"ab-testing-platform-go/internal/repository"
	"ab-testing-platform-go/internal/service"
	"ab-testing-platform-go/pkg/logger"
)

func main() {
	// 1. Загружаем конфигурацию
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Инициализируем логгер
	log := logger.NewLogger("go-edge-service")

	// 3. Инициализируем репозитории (Redis, PostgreSQL)
	redisRepo := repository.NewRedisRepository(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	dbRepo := repository.NewPostgresRepository(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	// 4. Инициализируем BatcherService
	batcherService := service.NewBatcherService(dbRepo)

	// 5. Инициализируем AssignerService (передаём batcher!)
	assignerService := service.NewAssignerService(redisRepo, dbRepo, batcherService, log)

	// 6. Создаём Gin роутер
	r := gin.Default()

	// 7. Настраиваем CORS (чтобы клиентский JS мог обращаться к API)
	configCORS := cors.DefaultConfig()
	configCORS.AllowAllOrigins = true
	configCORS.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	configCORS.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(configCORS))

	// 8. Регистрируем хендлеры
	assignHandler := handler.NewAssignHandler(assignerService)
	assignHandler.RegisterRoutes(r)

	// 9. Запускаем сервер в отдельной горутине
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		log.Info("Starting server on port " + cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(fmt.Sprintf("Server failed to start: %v", err))
		}
	}()

	// 10. Ожидаем сигнал остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// 11. Плавно останавливаем batcher перед остановкой сервера
	log.Info("Shutting down batcher...")
	batcherService.Stop()

	// 12. Плавно завершаем работу сервера
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	log.Info("Server exited")
}
