package rangecounter

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
)


func (drange DateRange) incrementDateForce(multiple int, at time.Time) time.Time {
	res, err := drange.incrementDate(multiple, at)
	if err != nil {
		panic(err)
	}
	return res
}

type basicDateCounter struct {
	drange DateRange
	backend Backend
}

func NewBasicDateCounter(drange DateRange, backend Backend) DateRangeCounter {
	return &basicDateCounter{
		drange: drange,
		backend: backend,
	}
}

func (b *basicDateCounter) QuerySum(ctx context.Context, at time.Time, bucketCount int) (int64, error) {
	at, err := b.drange.alignDate(at)
	if err != nil {
		return 0, errors.Wrap(err, "unable to align date")
	}

	keys := []string{}
	for i := 0;i<bucketCount;i++ {
		keys = append(keys, b.getKey(at))
		at, err = b.drange.incrementDate(-1, at)
		if err != nil {
			return 0, errors.Wrap(err, "unable to decrement date")
		}
	}

	results, err := b.backend.Query(ctx, keys)
	if err != nil {
		return 0, errors.Wrap(err, "unable to query counters")
	}

	sum := int64(0)
	for _, res := range results {
		sum = sum + res
	}

	return sum, nil
}

func (b *basicDateCounter) Increment(ctx context.Context, at time.Time, by int64) error {
	at, err := b.drange.alignDate(at)
	if err != nil {
		return errors.Wrap(err, "unable to align date")
	}

	return b.backend.Increment(ctx, []string{b.getKey(at)}, []int64{by})
}

func (b *basicDateCounter) getKey(at time.Time) string {
	return fmt.Sprintf("%v:%v", b.drange, at.Unix())
}

func (b *basicDateCounter) String() string {
	return "basicDateCounter"
}
