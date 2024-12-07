package cron

import (
	"sync"
	"time"

	"github.com/VandiKond/StocksBack/pkg/logger"
)
// It is a cron data
type Cron struct {
	running bool
	mu      sync.Mutex
	period  time.Duration
	start   int
	fn      func() error
	logger  *logger.Logger
}

// Creates a new cron
func New(period time.Duration, start int, fn func() error, logger *logger.Logger) *Cron {
	// Checks the time
	if start >= 24 {
		return nil
	}

	return &Cron{
		running: false,
		mu:      sync.Mutex{},
		period:  period,
		start:   start,
		fn:      fn,
		logger:  logger,
	}
}

// Runs the cron
func (c *Cron) Run() {

	c.mu.Lock()

	defer func() {

		// Defer closing all
		c.mu.Unlock()
		c.running = false
	}()

	// Not allowed to run twice
	if c.running {
		return
	}

	// Run
	c.running = true
	go c.run()
}

func (c *Cron) run() {
	// If starts not now
	if c.start > 0 {
		// Now and day
		now := time.Now()
		day := now.Day()

		// Getting next
		next := time.Date(now.Year(), now.Month(), day, c.start, 0, 0, 0, now.Location())
		if next.Before(now) {
			next = next.AddDate(0, 0, 1)
		}

		// Waiting til the day is need
		time.Sleep(next.Sub(now))
	}

	for {
		// Running the function
		err := c.fn()
		if err != nil {
			c.logger.Errorln(err)
		}

		// Waiting till the next
		time.Sleep(c.period)
	}
}
