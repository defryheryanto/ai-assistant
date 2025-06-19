package tools

import (
	"context"
	"sync"

	"github.com/tmc/langchaingo/llms"
)

// ContextWindow defines how conversational history is stored and retrieved.
type ContextWindow interface {
	GetHistory(ctx context.Context, id string) ([]llms.MessageContent, error)
	SaveHistory(ctx context.Context, id string, history []llms.MessageContent) error
}

// InMemoryContextWindow is a simple in-memory implementation of ContextWindow.
type InMemoryContextWindow struct {
	mu      sync.RWMutex
	storage map[string][]llms.MessageContent
}

// NewInMemoryContextWindow creates a new in-memory context window.
func NewInMemoryContextWindow() *InMemoryContextWindow {
	return &InMemoryContextWindow{storage: make(map[string][]llms.MessageContent)}
}

// GetHistory returns the history associated with the given id.
func (m *InMemoryContextWindow) GetHistory(ctx context.Context, id string) ([]llms.MessageContent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	history, ok := m.storage[id]
	if !ok {
		return []llms.MessageContent{}, nil
	}
	copyHist := make([]llms.MessageContent, len(history))
	copy(copyHist, history)
	return copyHist, nil
}

// SaveHistory saves the conversation history for the given id.
func (m *InMemoryContextWindow) SaveHistory(ctx context.Context, id string, history []llms.MessageContent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	copyHist := make([]llms.MessageContent, len(history))
	copy(copyHist, history)
	m.storage[id] = copyHist
	return nil
}
