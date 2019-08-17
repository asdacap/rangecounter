package rangecounter

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkBasicDateBehaviour(b *testing.B) {
	baseDate := time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local)
	tests := []struct {
		name                   string
		incrementFactor        int
		startIncrement         time.Time
		maxIncrementDateOffset int
		queryFactor            int
		startQuery             time.Time
		maxQueryDateOffset     int
		maxQueryRange          int
	}{
		{
			name:                   "range 5",
			incrementFactor:        1,
			startIncrement:         baseDate,
			maxIncrementDateOffset: 1000,
			queryFactor:            1,
			startQuery:             baseDate,
			maxQueryDateOffset:     1000,
			maxQueryRange:          5,
		},
		{
			name:                   "range 5 high increment",
			incrementFactor:        10,
			startIncrement:         baseDate,
			maxIncrementDateOffset: 1000,
			queryFactor:            1,
			startQuery:             baseDate,
			maxQueryDateOffset:     1000,
			maxQueryRange:          5,
		},
		{
			name:                   "range 5 high query",
			incrementFactor:        1,
			startIncrement:         baseDate,
			maxIncrementDateOffset: 1000,
			queryFactor:            10,
			startQuery:             baseDate,
			maxQueryDateOffset:     1000,
			maxQueryRange:          5,
		},
		{
			name:                   "range 20",
			incrementFactor:        1,
			startIncrement:         baseDate,
			maxIncrementDateOffset: 1000,
			queryFactor:            1,
			startQuery:             baseDate,
			maxQueryDateOffset:     1000,
			maxQueryRange:          20,
		},
		{
			name:                   "range 20 high increment",
			incrementFactor:        10,
			startIncrement:         baseDate,
			maxIncrementDateOffset: 1000,
			queryFactor:            1,
			startQuery:             baseDate,
			maxQueryDateOffset:     1000,
			maxQueryRange:          20,
		},
		{
			name:                   "range 20 high query",
			incrementFactor:        1,
			startIncrement:         baseDate,
			maxIncrementDateOffset: 1000,
			queryFactor:            10,
			startQuery:             baseDate,
			maxQueryDateOffset:     1000,
			maxQueryRange:          20,
		},
		{
			name:                   "range 100",
			incrementFactor:        1,
			startIncrement:         baseDate,
			maxIncrementDateOffset: 1000,
			queryFactor:            1,
			startQuery:             baseDate,
			maxQueryDateOffset:     1000,
			maxQueryRange:          100,
		},
		{
			name:                   "range 100 high increment",
			incrementFactor:        10,
			startIncrement:         baseDate,
			maxIncrementDateOffset: 1000,
			queryFactor:            1,
			startQuery:             baseDate,
			maxQueryDateOffset:     1000,
			maxQueryRange:          100,
		},
		{
			name:                   "range 100 high query",
			incrementFactor:        1,
			startIncrement:         baseDate,
			maxIncrementDateOffset: 1000,
			queryFactor:            10,
			startQuery:             baseDate,
			maxQueryDateOffset:     1000,
			maxQueryRange:          100,
		},
	}

	counterToTest := []struct {
		name    string
		factory func(dateRange DateRange, backend Backend) DateRangeCounter
	}{
		{
			"intBacked", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewBasicIntRangeCounter(backend), dateRange)
			},
		}, {
			"intRangeTreeBacked-2", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewRangeTreeIntCounter(backend, 2, 1), dateRange)
			},
		}, {
			"intRangeTreeBacked-2-2", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewRangeTreeIntCounter(backend, 2, 2), dateRange)
			},
		}, {
			"intRangeTreeBacked-4", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewRangeTreeIntCounter(backend, 4, 1), dateRange)
			},
		}, {
			"intRangeTreeBacked-8", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewRangeTreeIntCounter(backend, 8, 1), dateRange)
			},
		}, {
			"intRangeTreeBacked-16", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewRangeTreeIntCounter(backend, 16, 1), dateRange)
			},
		}, {
			"intRangeTreeBacked-8-2", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewRangeTreeIntCounter(backend, 8, 2), dateRange)
			},
		}, {
			"intRangeTreeBacked-4-4", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewRangeTreeIntCounter(backend, 4, 4), dateRange)
			},
		}, {
			"intRangeTreeBacked-12-2-to-second", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewIntRangeTranslator(NewRangeTreeIntCounter(backend, 12, 2), dateRange, Seconds), dateRange)
			},
		},
	}
	for _, d := range tests {
		b.Run(d.name, func(b *testing.B) {
			for _, counter := range counterToTest {
				b.Run("counter "+counter.name, func(b *testing.B) {
					rangeToTest := Hour
					dateCounter := counter.factory(rangeToTest, NewBenchmarkBackend())
					rand.Seed(0)
					ctx := context.Background()

					for i := 0; i < b.N; i++ {
						r := rand.Int() % (d.incrementFactor + d.queryFactor)
						if r < d.incrementFactor {
							offset := rand.Int() % d.maxIncrementDateOffset
							err := dateCounter.Increment(ctx, rangeToTest.incrementDateForce(offset, baseDate), 1)
							if err != nil {
								b.Fail()
							}
						} else {
							startOffset := rand.Int() % d.maxQueryDateOffset
							bucket := rand.Int() % d.maxQueryRange
							_, err := dateCounter.QuerySum(ctx, rangeToTest.incrementDateForce(startOffset, baseDate), bucket)
							if err != nil {
								b.Fail()
							}
						}
					}
				})
			}
		})
	}
}
