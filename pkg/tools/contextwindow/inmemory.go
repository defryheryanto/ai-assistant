package contextwindow

import (
	"context"
	"sync"

	"github.com/tmc/langchaingo/llms"
)

// InMemoryContextWindow is a simple in-memory implementation of tools.ContextWindow.
type InMemoryContextWindow struct {
	mu      sync.RWMutex
	storage map[string][]llms.MessageContent
	limit   int
}

type Option func(*InMemoryContextWindow)

// WithLimit sets the maximum number of history messages retained for each context ID.
func WithLimit(limit int) Option {
	return func(m *InMemoryContextWindow) {
		if limit > 0 {
			m.limit = limit
		}
	}
}

// NewInMemoryContextWindow creates a new in-memory context window.
func NewInMemoryContextWindow(opts ...Option) *InMemoryContextWindow {
	cw := &InMemoryContextWindow{storage: make(map[string][]llms.MessageContent)}
	for _, opt := range opts {
		opt(cw)
	}
	return cw
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
	if m.limit > 0 && len(copyHist) > m.limit {
		copyHist = copyHist[len(copyHist)-m.limit:]
	}
	m.storage[id] = copyHist
	return nil
}
