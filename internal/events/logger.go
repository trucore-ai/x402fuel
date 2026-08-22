package events

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/trucore-ai/x402fuel/internal/types"
)

type Logger struct {
	mu   sync.Mutex
	path string
	file *os.File
}

func NewLogger(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open event log: %w", err)
	}
	return &Logger{path: path, file: f}, nil
}

func (l *Logger) Log(event types.Event) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	data = append(data, '\n')
	if _, err := l.file.Write(data); err != nil {
		return fmt.Errorf("write event: %w", err)
	}
	return nil
}

func (l *Logger) Close() error {
	return l.file.Close()
}