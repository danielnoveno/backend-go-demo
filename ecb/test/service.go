package test

import (
	"context"
	"errors"
	"time"

	"deeply/ecb/config"
	"deeply/ecb/hardware"
	"deeply/ecb/state"
)

// DataRepository wraps persistence previously handled by Eloquent/Yajra.
type DataRepository interface {
	LoadHistory(ctx context.Context, line string, limit int) ([]HistoryRow, error)
	StoreScan(ctx context.Context, record ScanRecord) error
}

type HistoryRow struct {
	ID      int64
	SN      string
	FGType  string
	Line    string
	Created time.Time
}

type ScanRecord struct {
	SN       string
	FGType   string
	Line     string
	Step     int
	Result   string
	Metadata map[string]string
}

// Service coordinates scan workflow, hardware, and persistence.
type Service struct {
	dataRepo DataRepository
	cfg      *config.Service
	gpio     hardware.Adapter
	appState *state.AppState
}

func NewService(dataRepo DataRepository, cfg *config.Service, gpio hardware.Adapter, appState *state.AppState) *Service {
	return &Service{
		dataRepo: dataRepo,
		cfg:      cfg,
		gpio:     gpio,
		appState: appState,
	}
}

func (s *Service) LoadHistory(ctx context.Context, line string, limit int) ([]HistoryRow, error) {
	return s.dataRepo.LoadHistory(ctx, line, limit)
}

func (s *Service) InitializeRig(ctx context.Context) error {
	if err := s.gpio.Initialize(ctx); err != nil {
		return err
	}
	return s.appState.WithLock(ctx, func(state *state.AppState) error {
		state.LineActive = 0
		state.WaitForTest = false
		state.Testing = false
		return nil
	})
}

func (s *Service) StartTest(ctx context.Context, line int) error {
	if line < 0 {
		return errors.New("invalid line")
	}
	if err := s.gpio.Start(ctx, line); err != nil {
		return err
	}
	return s.appState.WithLock(ctx, func(state *state.AppState) error {
		state.LineActive = line
		state.WaitForTest = false
		state.Testing = true
		return nil
	})
}

func (s *Service) ResetRig(ctx context.Context, line int) error {
	if err := s.gpio.Reset(ctx, line); err != nil {
		return err
	}
	return s.appState.WithLock(ctx, func(state *state.AppState) error {
		state.LineActive = line
		state.Testing = false
		state.WaitForTest = false
		return nil
	})
}

func (s *Service) SetActiveLine(ctx context.Context, line int) error {
	if err := s.gpio.SetActiveLine(ctx, line); err != nil {
		return err
	}
	return s.appState.WithLock(ctx, func(state *state.AppState) error {
		state.LineActive = line
		return nil
	})
}

func (s *Service) PollStatus(ctx context.Context) (hardware.Status, error) {
	status, err := s.gpio.ReadStatus(ctx)
	if err != nil {
		return hardware.Status{}, err
	}
	return status, nil
}

func (s *Service) StoreScan(ctx context.Context, record ScanRecord) error {
	return s.dataRepo.StoreScan(ctx, record)
}
