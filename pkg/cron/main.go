package cron

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	*cron.Cron
}

func NewCron() *Cron {
	return &Cron{cron.New()}
}

func (c *Cron) RunDaily(task func()) {
	_, _ = c.AddFunc("@daily", task)
	c.Start()
}

func (c *Cron) RunEverySeconds(task func(), seconds int) {
	_, _ = c.AddFunc(fmt.Sprintf("@every %ds", seconds), task)
	c.Start()
}

// RunOnce schedules a task to run once when the current time passes a specified Unix timestamp
func (c *Cron) RunOnce(task func(), targetTimePtr *int64) {
	// Run a check every minute to see if we've passed the target time
	_, _ = c.AddFunc("@every 1m", func() {
		now := time.Now().Unix()

		// Check if current time has passed the target time
		if now >= *targetTimePtr {
			// Execute the task
			task()

			// Update the target time to now + 1 minute
			*targetTimePtr = now + 60
		}
	})
	c.Start()
}
