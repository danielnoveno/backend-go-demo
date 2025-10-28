package config

import (
	"context"
	"sync"
)

// Settings mirrors the rows stored in Ecbconfig.
type Settings struct {
	LocalIP  string
	SimoIP   string
	UseWLAN  bool
	Theme    string
	LineType string
	Location string
	LineName string
	Tacktime int
}

// Repository abstracts DB access.
type Repository interface {
	LoadSettings(ctx context.Context) (Settings, error)
	SaveSettings(ctx context.Context, settings Settings) error
}

// Service caches settings and provides goroutine safe access.
type Service struct {
	repo Repository

	mu       sync.RWMutex
	settings Settings
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Load(ctx context.Context) (Settings, error) {
	s.mu.RLock()
	if s.settings != (Settings{}) {
		defer s.mu.RUnlock()
		return s.settings, nil
	}
	s.mu.RUnlock()

	settings, err := s.repo.LoadSettings(ctx)
	if err != nil {
		return Settings{}, err
	}
	s.mu.Lock()
	s.settings = settings
	s.mu.Unlock()
	return settings, nil
}

func (s *Service) Save(ctx context.Context, settings Settings) error {
	if err := s.repo.SaveSettings(ctx, settings); err != nil {
		return err
	}
	s.mu.Lock()
	s.settings = settings
	s.mu.Unlock()
	return nil
}
