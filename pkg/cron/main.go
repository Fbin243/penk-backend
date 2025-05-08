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

// RunAtTimestampAndReschedule schedules a task to run when the current time passes a specified Unix timestamp
// and reschedules it based on the new timestamp returned by the task
func (c *Cron) RunAtTimestampAndReschedule(task func() *int64, targetTimePtr *int64) {
	// Run a check every minute to see if we've passed the target time
	_, _ = c.AddFunc("@every 1m", func() {
		// Check if current time has passed the target time
		if targetTimePtr != nil && time.Now().Unix() >= *targetTimePtr {
			targetTimePtr = task()
		}
	})
	c.Start()
}
