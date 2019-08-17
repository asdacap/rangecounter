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
		startIncrement         time.Time
		maxIncrementDateOffset int
		startQuery             time.Time
		maxQueryDateOffset     int
		maxQueryRange          int
	}{
		{
			name:                   "range 5",
			startIncrement:         baseDate,
			maxIncrementDateOffset: 100000,
			startQuery:             baseDate,
			maxQueryDateOffset:     100000,
			maxQueryRange:          5,
		},
		{
			name:                   "range 20",
			startIncrement:         baseDate,
			maxIncrementDateOffset: 100000,
			startQuery:             baseDate,
			maxQueryDateOffset:     100000,
			maxQueryRange:          20,
		},
		{
			name:                   "range 100",
			startIncrement:         baseDate,
			maxIncrementDateOffset: 100000,
			startQuery:             baseDate,
			maxQueryDateOffset:     100000,
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
			"intBacked-to-second", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewIntRangeTranslator(NewBasicIntRangeCounter(backend), dateRange, Seconds), dateRange)
			},
		}, {
			"intRangeTreeBacked-8-4-to-second", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewIntRangeTranslator(NewRangeTreeIntCounter(backend, 12, 2), dateRange, Seconds), dateRange)
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
					rangeToTest := Minute
					backend := NewBenchmarkBackend()
					dateCounter := counter.factory(rangeToTest, backend)
					rand.Seed(0)
					ctx := context.Background()

					for i := 0; i < b.N; i++ {
						offset := rand.Int() % d.maxIncrementDateOffset
						err := dateCounter.Increment(ctx, rangeToTest.incrementDateForce(offset, baseDate), 1)
						if err != nil {
							b.Fail()
						}
					}
					b.ReportMetric(float64(backend.incrementKeyTouched)/float64(b.N), "incrementKeyTouched")

					for i := 0; i < b.N; i++ {
						startOffset := rand.Int() % d.maxQueryDateOffset
						bucket := rand.Int() % d.maxQueryRange
						_, err := dateCounter.QuerySum(ctx, rangeToTest.incrementDateForce(startOffset, baseDate), bucket)
						if err != nil {
							b.Fail()
						}
					}
					b.ReportMetric(float64(backend.queryKeyTouched)/float64(b.N), "queryKeyTouched")
				})
			}
		})
	}
}
