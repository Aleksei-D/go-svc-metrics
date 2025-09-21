package errors

import "errors"

var (
	ErrInvalidMetricValue      = errors.New("invalid metric value")
	ErrInvalidCounterOperation = errors.New("invalid counter operation")
	ErrInvalidCGaugeOperation  = errors.New("invalid gauge operation")
	ErrInvalidMetricVType      = errors.New("invalid metric type")
)
