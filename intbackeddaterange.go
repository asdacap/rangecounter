package rangecounter

import (
	"context"
	"time"
)

type intBackedDateRange struct {
	nativeRange  DateRange
	backingRange IntRangeCounter
}

func (ibdr *intBackedDateRange) QuerySum(ctx context.Context, at time.Time, bucketCount int) (int64, error) {
	durationNano := ibdr.nativeRange.getDuration().Nanoseconds()
	endIndex := at.UnixNano()/durationNano
	startIndex := endIndex - int64(bucketCount) + 1
	return ibdr.backingRange.QuerySum(ctx, startIndex, endIndex)
}

func (ibdr *intBackedDateRange) Increment(ctx context.Context, at time.Time, by int64) (error) {
	durationNano := ibdr.nativeRange.getDuration().Nanoseconds()
	index := at.UnixNano()/durationNano
	return ibdr.backingRange.Increment(ctx, index, by)
}

func NewIntBackedDateRange(backingRange IntRangeCounter, nativeRange DateRange) DateRangeCounter {
	return &intBackedDateRange{
		backingRange: backingRange,
		nativeRange: nativeRange,
	}
}
