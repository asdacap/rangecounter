package rangecounter

import "context"

type basicIntRangeCounter struct {
	store map[int64]int64
}

func (birc *basicIntRangeCounter) QuerySum(ctx context.Context, from, to int64) (int64, error) {
	sum := int64(0)
	for ;from<=to;from++ {
		sum += birc.store[from]
	}
	return sum, nil
}

func (birc *basicIntRangeCounter) Increment(ctx context.Context, at int64, by int64) (error) {
	birc.store[at] = birc.store[at] + by
	return nil
}

func NewBasicIntRangeCounter() IntRangeCounter {
	return &basicIntRangeCounter{
		store: map[int64]int64{},
	}
}
