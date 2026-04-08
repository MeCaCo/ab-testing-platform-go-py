package logger

import (
	"log"
	"time"
)

// SimpleLogger — простой логгер для начала
type SimpleLogger struct {
	port string
}

func NewLogger(port string) *SimpleLogger {
	return &SimpleLogger{port: port}
}

func (l *SimpleLogger) Info(msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[%s] [INFO] [Port %s] %s\n", timestamp, l.port, msg)
}

func (l *SimpleLogger) Error(msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[%s] [ERROR] [Port %s] %s\n", timestamp, l.port, msg)
}

func (l *SimpleLogger) Warn(msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[%s] [WARN] [Port %s] %s\n", timestamp, l.port, msg)
}
