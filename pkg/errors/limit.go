package errors

import "fmt"

var (
	ErrorMetricLimitReached   = fmt.Errorf("metric limit reached")
	ErrorPropertyLimitReached = fmt.Errorf("property limit reached")
)
