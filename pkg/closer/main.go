package closer

import (
	"context"
	"sync"

	"github.com/VandiKond/StocksBack/pkg/logger"
	"github.com/VandiKond/vanerrors"
)

// Errors
const (
	ShutdownCancelled = "shutdown cancelled"
	GotSomeErrors     = "got some errors"
)

// Func for closing 
type Fn func() error

// The closer
type Closer struct {
	logger *logger.Logger
	mu     sync.Mutex
	funcs  []Fn
}

// Ads a new closer
func (c *Closer) Add(fn Fn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, fn)
}

// Closes the closer
func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		errs = make([]error, 0, len(c.funcs))
		wg   sync.WaitGroup
	)

	// Closes all funcs
	for _, f := range c.funcs {
		wg.Add(1)
		go func(f Fn) {
			defer wg.Done()

			if err := f(); err != nil {
				errs = append(errs, err)
			}

		}(f)
	}

	// Wait for finish of all
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		break
	case <-ctx.Done():
		return vanerrors.NewWrap(ShutdownCancelled, ctx.Err(), vanerrors.EmptyHandler)
	}

	// Adds errors
	if len(errs) > 0 {
		for _, err := range errs {
			c.logger.Errorln(err)
		}
		return vanerrors.NewSimple(GotSomeErrors)
	}

	return nil
}

func New(logger *logger.Logger) *Closer {
	return &Closer{logger: logger}
}
