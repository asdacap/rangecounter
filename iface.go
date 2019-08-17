package rangecounter

import (
	"context"
	"time"
)

// Backend is where the values are actually stored.
// It should be a key-value store with support for pipelines (like redis) and incremenet.
type Backend interface {
	Query(ctx context.Context, keys []string) ([]int64, error)
	Increment(ctx context.Context, keys []string, values []int64) error
}

// IntRangeCounter query count stuff with int64 as its keys
type IntRangeCounter interface {
	QuerySum(ctx context.Context, from, to int64) (int64, error)
	Increment(ctx context.Context, at int64, by int64) (error)
}

// DateRangeCounter query count stuff with time.Time as its keys and rangeCount
// A range is a duration, for example second, minutes or hours, and it should be fixed for the implementation
// The bucketCount is the count of `range` before `at` (inclusive of `at`)
// The implementation should align to the boundary of range.
// On some implementation it should use IntRangeCounter under
type DateRangeCounter interface {
	QuerySum(ctx context.Context, at time.Time, bucketCount int) (int64, error)
	Increment(ctx context.Context, at time.Time, by int64) (error)
}

