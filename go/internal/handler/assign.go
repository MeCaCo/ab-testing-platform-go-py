// go/internal/handler/assign.go
package handler

import (
	"ab-testing-platform-go/internal/model"
	"ab-testing-platform-go/internal/service"
	"net/http"
	"time" // Импортируем пакет time

	"github.com/gin-gonic/gin"
)

// Убедись, что структура AssignHandler объявлена правильно
type AssignHandler struct {
	assignerService *service.AssignerService // <-- Вот это поле!
}

// Убедись, что NewAssignHandler возвращает правильную структуру
func NewAssignHandler(assignerService *service.AssignerService) *AssignHandler {
	return &AssignHandler{assignerService: assignerService}
}

// AssignVariant handles the assignment of a user to a test group (A or B)
// GET /api/v1/assign?test_id=123&user_id=user456
func (h *AssignHandler) AssignVariant(c *gin.Context) {
	testID := c.Query("test_id")
	userID := c.Query("user_id")

	if testID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "test_id and user_id are required",
		})
		return
	}

	ctx := c.Request.Context()
	// ПРАВИЛЬНО: вызываем метод у поля assignerService
	variant, err := h.assignerService.AssignVariant(ctx, testID, userID) // <-- СТРОКА 44 (примерно)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to assign variant",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"test_id":   testID,
		"user_id":   userID,
		"variant":   variant,    // "A" or "B"
		"timestamp": time.Now(), // <-- Заменили c.Now() на time.Now()
	})
}

// RecordEvent handles incoming events (impression, click, conversion)
// POST /api/v1/event
// Body: {"test_id": "...", "user_id": "...", "variant": "...", "type": "...", "value": 1.0}
func (h *AssignHandler) RecordEvent(c *gin.Context) {
	var event model.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	ctx := c.Request.Context()
	err := h.assignerService.RecordEvent(ctx, &event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to record event",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "event recorded",
		"event_id": event.ID,
	})
}

// HealthCheck проверяет, жив ли сервис
func (h *AssignHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "go-edge-service",
		"timestamp": time.Now(), // <-- Заменили c.Now() на time.Now()
	})
}

// RegisterRoutes регистрирует маршруты для распределения и событий
func (h *AssignHandler) RegisterRoutes(r *gin.Engine) {
	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/assign", h.AssignVariant)
		apiV1.POST("/event", h.RecordEvent) // <-- Добавлен маршрут для события
	}

	// Health check
	r.GET("/health", h.HealthCheck)
}
