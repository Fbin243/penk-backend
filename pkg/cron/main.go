package cron

import (
	"fmt"

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
