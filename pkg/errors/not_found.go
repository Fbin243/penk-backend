package errors

import "fmt"

var (
	ErrorMetricNotFound   = fmt.Errorf("metric not found")
	ErrorPropertyNotFound = fmt.Errorf("property not found")
)
