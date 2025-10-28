package hardware

import (
	"context"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

// Status represents the digital signals sampled from the ECB rig.
type Status struct {
	Pass      bool
	Fail      bool
	UnderTest bool
	Line      int
	Timestamp time.Time
	Mode      string
}

// Adapter exposes the operations previously done via shell exec.
type Adapter interface {
	Initialize(ctx context.Context) error
	Start(ctx context.Context, line int) error
	Reset(ctx context.Context, line int) error
	SetActiveLine(ctx context.Context, line int) error
	ReadStatus(ctx context.Context) (Status, error)
}

// ExecAdapter shells out to the legacy gpio utility (wiringPi).
type ExecAdapter struct {
	mu sync.Mutex
}

func NewExecAdapter() *ExecAdapter {
	return &ExecAdapter{}
}

func (a *ExecAdapter) Initialize(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	// Equivalent wiringPi setup.
	if err := run(ctx, "gpio", "mode", "23", "out"); err != nil {
		return err
	}
	if err := run(ctx, "gpio", "mode", "24", "out"); err != nil {
		return err
	}
	if err := run(ctx, "gpio", "mode", "25", "out"); err != nil {
		return err
	}
	if err := run(ctx, "gpio", "mode", "28", "out"); err != nil {
		return err
	}
	if err := run(ctx, "gpio", "mode", "29", "out"); err != nil {
		return err
	}
	if err := run(ctx, "gpio", "write", "25", "0"); err != nil {
		return err
	}
	if err := run(ctx, "gpio", "write", "28", "1"); err != nil {
		return err
	}
	return run(ctx, "gpio", "write", "29", "1")
}

func (a *ExecAdapter) Start(ctx context.Context, line int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := run(ctx, "gpio", "write", "25", strconv.Itoa(line)); err != nil {
		return err
	}
	if err := run(ctx, "gpio", "write", "23", "1"); err != nil {
		return err
	}
	if err := run(ctx, "gpio", "write", "28", "0"); err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	if err := run(ctx, "gpio", "write", "23", "0"); err != nil {
		return err
	}
	if err := run(ctx, "gpio", "write", "28", "1"); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	return run(ctx, "gpio", "write", "24", "1")
}

func (a *ExecAdapter) Reset(ctx context.Context, line int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := run(ctx, "gpio", "write", "25", strconv.Itoa(line)); err != nil {
		return err
	}
	if err := run(ctx, "gpio", "write", "23", "1"); err != nil {
		return err
	}
	if err := run(ctx, "gpio", "write", "28", "0"); err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	if err := run(ctx, "gpio", "write", "23", "0"); err != nil {
		return err
	}
	if err := run(ctx, "gpio", "write", "28", "1"); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	if err := run(ctx, "gpio", "write", "23", "1"); err != nil {
		return err
	}
	if err := run(ctx, "gpio", "write", "28", "0"); err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	if err := run(ctx, "gpio", "write", "23", "0"); err != nil {
		return err
	}
	return run(ctx, "gpio", "write", "28", "1")
}

func (a *ExecAdapter) SetActiveLine(ctx context.Context, line int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return run(ctx, "gpio", "write", "25", strconv.Itoa(line))
}

func (a *ExecAdapter) ReadStatus(ctx context.Context) (Status, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	pass, err := readPin(ctx, 2)
	if err != nil {
		return Status{}, err
	}
	fail, err := readPin(ctx, 21)
	if err != nil {
		return Status{}, err
	}
	undertest, err := readPin(ctx, 22)
	if err != nil {
		return Status{}, err
	}
	line, err := readPin(ctx, 25)
	if err != nil {
		return Status{}, err
	}
	return Status{
		Pass:      pass == 0,
		Fail:      fail == 0,
		UnderTest: undertest == 0,
		Line:      line,
		Timestamp: time.Now(),
	}, nil
}

func run(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.Run()
}

func readPin(ctx context.Context, pin int) (int, error) {
	out, err := exec.CommandContext(ctx, "gpio", "read", strconv.Itoa(pin)).Output()
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(string(out[:1]))
	if err != nil {
		return 0, err
	}
	return value, nil
}
