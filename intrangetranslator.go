package rangecounter

import "context"

type intRangeTranslator struct {
	innerCounter IntRangeCounter
	factor       int64
}

func (i *intRangeTranslator) Increment(ctx context.Context, at int64, by int64) error {
	return i.innerCounter.Increment(ctx, at*i.factor, by)
}

func (i *intRangeTranslator) QuerySum(ctx context.Context, from, to int64) (int64, error) {
	return i.innerCounter.QuerySum(ctx, from*i.factor, (to*i.factor)+i.factor-1)
}

func NewIntRangeTranslator(innerCounter IntRangeCounter, fromDateRange, toDateRange DateRange) IntRangeCounter {
	if toDateRange.getDuration().Nanoseconds() > fromDateRange.getDuration().Nanoseconds() {
		panic("to date range must be smaller than from date range")
	}
	factor := fromDateRange.getDuration().Nanoseconds() / toDateRange.getDuration().Nanoseconds()
	return &intRangeTranslator{
		innerCounter: innerCounter,
		factor:       factor,
	}
}
