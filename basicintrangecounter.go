package rangecounter

import (
	"context"
	"fmt"
)

type basicIntRangeCounter struct {
	backend Backend
}

func (birc *basicIntRangeCounter) QuerySum(ctx context.Context, from, to int64) (int64, error) {
	keys := []string{}
	for ; from <= to; from++ {
		keys = append(keys, fmt.Sprint(from))
	}
	ints, err := birc.backend.Query(ctx, keys)
	if err != nil {
		return 0, err
	}
	sum := int64(0)
	for _, it := range ints {
		sum += it
	}
	return sum, nil
}

func (birc *basicIntRangeCounter) Increment(ctx context.Context, at int64, by int64) error {
	return birc.backend.Increment(ctx, []string{fmt.Sprint(at)}, []int64{by})
}

func NewBasicIntRangeCounter(backend Backend) IntRangeCounter {
	return &basicIntRangeCounter{
		backend: backend,
	}
}
