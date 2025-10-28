package state

import (
	"context"
	"sync"
	"time"
)

// AppState centralises transient values that used to live in Laravel's session/globals.
type AppState struct {
	mu sync.RWMutex

	Theme          string
	Mode           string
	LineActive     int
	WaitForTest    bool
	Testing        bool
	CurrentMenuURL string

	LastModeChange time.Time
}

func New() *AppState {
	return &AppState{
		Theme:      "default",
		Mode:       "LIVE",
		LineActive: 0,
	}
}

func (s *AppState) WithLock(ctx context.Context, fn func(*AppState) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fn(s)
}

func (s *AppState) WithRead(ctx context.Context, fn func(*AppState)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fn(s)
}
