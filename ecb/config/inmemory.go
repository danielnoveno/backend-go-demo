package config

import "context"

type InMemoryRepository struct {
	settings Settings
}

func NewInMemoryRepository(settings Settings) *InMemoryRepository {
	return &InMemoryRepository{settings: settings}
}

func (r *InMemoryRepository) LoadSettings(ctx context.Context) (Settings, error) {
	return r.settings, nil
}

func (r *InMemoryRepository) SaveSettings(ctx context.Context, settings Settings) error {
	r.settings = settings
	return nil
}
