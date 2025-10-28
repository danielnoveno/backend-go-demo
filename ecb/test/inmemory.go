package test

import (
	"context"
	"sync"
	"time"
)

type InMemoryRepository struct {
	mu      sync.RWMutex
	records []ScanRecord
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{}
}

func (r *InMemoryRepository) LoadHistory(ctx context.Context, line string, limit int) ([]HistoryRow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []HistoryRow
	for i := len(r.records) - 1; i >= 0; i-- {
		rec := r.records[i]
		if rec.Line != line {
			continue
		}
		result = append(result, HistoryRow{
			ID:      int64(i),
			SN:      rec.SN,
			FGType:  rec.FGType,
			Line:    rec.Line,
			Created: time.Now(),
		})
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}

func (r *InMemoryRepository) StoreScan(ctx context.Context, record ScanRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records = append(r.records, record)
	return nil
}
