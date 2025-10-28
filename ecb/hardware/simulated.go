package hardware

import (
	"context"
	"sync"
	"time"
)

// SimulatedAdapter is used for development without wiringPi.
type SimulatedAdapter struct {
	mu sync.RWMutex

	status Status
}

func NewSimulatedAdapter() *SimulatedAdapter {
	return &SimulatedAdapter{
		status: Status{Line: 0, Timestamp: time.Now()},
	}
}

func (a *SimulatedAdapter) Initialize(ctx context.Context) error {
	a.mu.Lock()
	a.status = Status{Line: 0, Timestamp: time.Now()}
	a.mu.Unlock()
	return nil
}

func (a *SimulatedAdapter) Start(ctx context.Context, line int) error {
	a.mu.Lock()
	a.status.Line = line
	a.status.UnderTest = true
	a.status.Pass = false
	a.status.Fail = false
	a.status.Timestamp = time.Now()
	a.mu.Unlock()
	go a.completeTest(line)
	return nil
}

func (a *SimulatedAdapter) Reset(ctx context.Context, line int) error {
	a.mu.Lock()
	a.status.Line = line
	a.status.Pass = false
	a.status.Fail = false
	a.status.UnderTest = false
	a.status.Timestamp = time.Now()
	a.mu.Unlock()
	return nil
}

func (a *SimulatedAdapter) SetActiveLine(ctx context.Context, line int) error {
	a.mu.Lock()
	a.status.Line = line
	a.status.Timestamp = time.Now()
	a.mu.Unlock()
	return nil
}

func (a *SimulatedAdapter) ReadStatus(ctx context.Context) (Status, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.status, nil
}

func (a *SimulatedAdapter) completeTest(line int) {
	time.Sleep(2 * time.Second)
	a.mu.Lock()
	defer a.mu.Unlock()
	a.status.Line = line
	a.status.UnderTest = false
	a.status.Pass = true
	a.status.Fail = false
	a.status.Timestamp = time.Now()
}
