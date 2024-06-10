package timetrack

import (
    "time"
)

type Timer struct {
    StartTime time.Time `json:"start"`
    EndTime   time.Time `json:"end"`
}

func (t *Timer) Start() {
    t.StartTime = time.Now()
}

func (t *Timer) Stop() {
    t.EndTime = time.Now()
}

func (t *Timer) Duration() time.Duration {
    return t.EndTime.Sub(t.StartTime)
}