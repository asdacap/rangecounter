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
			"dateRange", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewBasicDateCounter(dateRange, backend)
			},
		}, {
			"tree-1", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewRangeTreeIntCounter(backend, 1, 1), dateRange)
			},
		}, {
			"tree-2", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewRangeTreeIntCounter(backend, 2, 1), dateRange)
			},
		}, {
			"tree-2-2", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewRangeTreeIntCounter(backend, 2, 2), dateRange)
			},
		}, {
			"tree-4", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewRangeTreeIntCounter(backend, 4, 1), dateRange)
			},
		}, {
			"tree-8", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewRangeTreeIntCounter(backend, 8, 1), dateRange)
			},
		}, {
			"tree-16", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewRangeTreeIntCounter(backend, 16, 1), dateRange)
			},
		}, {
			"tree-8-2", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewRangeTreeIntCounter(backend, 8, 2), dateRange)
			},
		}, {
			"tree-4-4", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewRangeTreeIntCounter(backend, 4, 4), dateRange)
			},
		}, {
			"tree-to-second", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewIntRangeTranslator(NewBasicIntRangeCounter(backend), dateRange, Seconds), dateRange)
			},
		}, {
			"tree-8-1-to-second", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewIntRangeTranslator(NewRangeTreeIntCounter(backend, 8, 1), dateRange, Seconds), dateRange)
			},
		}, {
			"tree-16-1-to-second", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewIntRangeTranslator(NewRangeTreeIntCounter(backend, 16, 1), dateRange, Seconds), dateRange)
			},
		}, {
			"tree-32-1-to-second", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewIntRangeTranslator(NewRangeTreeIntCounter(backend, 32, 1), dateRange, Seconds), dateRange)
			},
		}, {
			"tree-8-2-to-second", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewIntRangeTranslator(NewRangeTreeIntCounter(backend, 8, 2), dateRange, Seconds), dateRange)
			},
		}, {
			"tree-16-2-to-second", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewIntRangeTranslator(NewRangeTreeIntCounter(backend, 16, 2), dateRange, Seconds), dateRange)
			},
		}, {
			"tree-4-4-to-second", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewIntRangeTranslator(NewRangeTreeIntCounter(backend, 4, 4), dateRange, Seconds), dateRange)
			},
		}, {
			"tree-8-4-to-second", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewIntRangeTranslator(NewRangeTreeIntCounter(backend, 8, 4), dateRange, Seconds), dateRange)
			},
		}, {
			"tree-4-8-to-second", func(dateRange DateRange, backend Backend) DateRangeCounter {
				return NewIntBackedDateRange(NewIntRangeTranslator(NewRangeTreeIntCounter(backend, 4, 8), dateRange, Seconds), dateRange)
			},
		},
	}
	for _, d := range tests {
		b.Run(d.name, func(b *testing.B) {
			for _, counter := range counterToTest {
				b.Run(counter.name, func(b *testing.B) {
					rangeToTest := Minute
					round := 10000

					for bi := 0; bi < b.N; bi++ {
						backend := NewBenchmarkBackend()
						dateCounter := counter.factory(rangeToTest, backend)
						rand.Seed(0)
						ctx := context.Background()

						for i := 0; i < round; i++ {
							offset := rand.Int() % d.maxIncrementDateOffset
							err := dateCounter.Increment(ctx, rangeToTest.incrementDateForce(offset, baseDate), 1)
							if err != nil {
								b.Fail()
							}
						}
						b.ReportMetric(float64(backend.incrementKeyTouched)/float64(round), "incrementKeyTouched")

						for i := 0; i < round; i++ {
							startOffset := rand.Int() % d.maxQueryDateOffset
							bucket := (rand.Int() % d.maxQueryRange)+1
							_, err := dateCounter.QuerySum(ctx, rangeToTest.incrementDateForce(startOffset, baseDate), bucket)
							if err != nil {
								b.Fail()
							}
						}
						b.ReportMetric(float64(backend.queryKeyTouched)/float64(round), "queryKeyTouched")
						b.ReportMetric(float64(len(backend.store))/float64(round), "keyUsed")
					}
				})
			}
		})
	}
}
