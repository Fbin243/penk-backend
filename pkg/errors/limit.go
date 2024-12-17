package errors

import "fmt"

var (
	ErrorCharacterLimitReached = fmt.Errorf("character limit reached")
	ErrorMetricLimitReached    = fmt.Errorf("metric limit reached")
	ErrorPropertyLimitReached  = fmt.Errorf("property limit reached")
)
