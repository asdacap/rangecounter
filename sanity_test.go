package rangecounter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateRangeBehaviour(t *testing.T) {
	baseDate := time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local)

	type inputReq struct {
		dateOffset int
		by         int64
	}
	type queryReq struct {
		dateOffset  int
		bucketCount int
		expected    int64
	}
	tests := []struct {
		name    string
		inputs  []inputReq
		queries []queryReq
	}{
		{
			name: "happy path",
			inputs: []inputReq{
				{0, 1},
			},
			queries: []queryReq{
				{0, 1, 1},
			},
		},
		{
			name: "happy path 2",
			inputs: []inputReq{
				{0, 1},
				{1, 1},
			},
			queries: []queryReq{
				{1, 1, 1},
				{1, 2, 2},
			},
		},
		{
			name: "happy path 3",
			inputs: []inputReq{
				{0, 1},
				{0, 1},
				{1, 1},
				{3, 3},
			},
			queries: []queryReq{
				{0, 1, 2},
				{1, 2, 3},
				{9, 10, 6},
				{2, 1, 0},
				{3, 3, 4},
			},
		},
	}

	counterToTest := map[string]func(DateRange) DateRangeCounter{
		"dateRange": func(dateRange DateRange) DateRangeCounter {
			return NewBasicDateCounter(dateRange, NewInMemoryBackend())
		},
		"intBacked": func(dateRange DateRange) DateRangeCounter {
			return NewIntBackedDateRange(NewBasicIntRangeCounter(NewInMemoryBackend()), dateRange)
		},
		"intRangeTreeBacked1": func(dateRange DateRange) DateRangeCounter {
			return NewIntBackedDateRange(NewRangeTreeIntCounter(NewInMemoryBackend(), 1, 1), dateRange)
		},
		"intRangeTreeBacked": func(dateRange DateRange) DateRangeCounter {
			return NewIntBackedDateRange(NewRangeTreeIntCounter(NewInMemoryBackend(), 8, 1), dateRange)
		},
		"intRangeTreeBacked2": func(dateRange DateRange) DateRangeCounter {
			return NewIntBackedDateRange(NewRangeTreeIntCounter(NewInMemoryBackend(), 8, 2), dateRange)
		},
		"intRangeTreeBacked3": func(dateRange DateRange) DateRangeCounter {
			return NewIntBackedDateRange(NewRangeTreeIntCounter(NewInMemoryBackend(), 16, 3), dateRange)
		},
		"intRangeTreeBacked4": func(dateRange DateRange) DateRangeCounter {
			return NewIntBackedDateRange(NewIntRangeTranslator(NewRangeTreeIntCounter(NewInMemoryBackend(), 16, 3), dateRange, Seconds), dateRange)
		},
	}
	rangeToTests := []DateRange{Hour}
	for counterName, dateCounterFactory := range counterToTest {
		t.Run("counter "+counterName, func(t *testing.T) {
			for _, rangeToTest := range rangeToTests {
				t.Run("range "+fmt.Sprint(rangeToTest), func(t *testing.T) {
					for _, d := range tests {
						dateCounter := dateCounterFactory(rangeToTest)
						t.Run(d.name, func(t *testing.T) {
							ctx := context.Background()
							for _, inp := range d.inputs {
								err := dateCounter.Increment(ctx, rangeToTest.incrementDateForce(inp.dateOffset, baseDate), inp.by)
								assert.NoError(t, err)
							}

							for _, q := range d.queries {
								ans, err := dateCounter.QuerySum(ctx, rangeToTest.incrementDateForce(q.dateOffset, baseDate), q.bucketCount)
								assert.NoError(t, err)
								assert.EqualValues(t, q.expected, ans)
							}
						})
					}
				})
			}
		})
	}
}

func TestIntRangeCounterBehavior(t *testing.T) {
	type inputReq struct {
		at int64
		by int64
	}
	type queryReq struct {
		from     int64
		to       int64
		expected int64
	}
	tests := []struct {
		name    string
		inputs  []inputReq
		queries []queryReq
	}{
		{
			name: "happy path",
			inputs: []inputReq{
				{0, 1},
			},
			queries: []queryReq{
				{0, 1, 1},
			},
		},
		{
			name: "happy path 2",
			inputs: []inputReq{
				{0, 1},
				{1, 1},
			},
			queries: []queryReq{
				{1, 1, 1},
				{1, 2, 1},
			},
		},
		{
			name: "happy path 3",
			inputs: []inputReq{
				{0, 1},
				{0, 1},
				{1, 1},
				{3, 3},
			},
			queries: []queryReq{
				{0, 0, 2},
				{0, 1, 3},
				{0, 10, 6},
				{2, 2, 0},
				{1, 3, 4},
			},
		},
		{
			name: "happy path 4",
			inputs: []inputReq{
				{0, 1},
				{1, 1},
				{2, 1},
				{3, 1},
				{4, 1},
				{5, 1},
				{6, 1},
				{7, 1},
				{8, 1},
				{9, 1},
				{10, 1},
				{11, 1},
				{12, 1},
				{13, 1},
				{14, 1},
				{15, 1},
				{16, 1},
			},
			queries: []queryReq{
				{0, 0, 1},
				{0, 1, 2},
				{0, 2, 3},
				{0, 3, 4},
				{0, 4, 5},
				{1, 5, 5},
				{2, 6, 5},
				{3, 7, 5},
				{4, 8, 5},
				{5, 9, 5},
				{6, 10, 5},
				{7, 11, 5},
				{8, 12, 5},
				{9, 13, 5},
				{10, 14, 5},
				{11, 15, 5},
				{12, 16, 5},
				{13, 17, 4},
				{14, 18, 3},
				{15, 19, 2},
				{16, 20, 1},
			},
		},
	}

	counterToTest := map[string]func() IntRangeCounter{
		"intBacked": func() IntRangeCounter {
			return NewBasicIntRangeCounter(NewInMemoryBackend())
		},
		"intRangeTreeBacked": func() IntRangeCounter {
			return NewRangeTreeIntCounter(NewInMemoryBackend(), 8, 1)
		},
		"intRangeTreeBacked2": func() IntRangeCounter {
			return NewRangeTreeIntCounter(NewInMemoryBackend(), 8, 2)
		},
		"intRangeTreeBacked3": func() IntRangeCounter {
			return NewRangeTreeIntCounter(NewInMemoryBackend(), 16, 3)
		},
		"intRangeTreeBacked4": func() IntRangeCounter {
			return NewRangeTreeIntCounter(NewInMemoryBackend(), 50, 1)
		},
		"intRangeTreeBacked5": func() IntRangeCounter {
			return NewRangeTreeIntCounter(NewInMemoryBackend(), 50, 1)
		},
	}
	for counterName, intCounterFactory := range counterToTest {
		t.Run("counter "+counterName, func(t *testing.T) {
			for _, d := range tests {
				intCounter := intCounterFactory()
				t.Run(d.name, func(t *testing.T) {
					ctx := context.Background()
					for _, inp := range d.inputs {
						err := intCounter.Increment(ctx, inp.at, inp.by)
						assert.NoError(t, err)
					}

					for _, q := range d.queries {
						ans, err := intCounter.QuerySum(ctx, q.from, q.to)
						assert.NoError(t, err)
						assert.EqualValues(t, q.expected, ans)
					}
				})
			}
		})
	}
}
