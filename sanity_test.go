package rangecounter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBasicDateBehaviour(t *testing.T) {
	baseDate := time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local)

	type inputReq struct {
		dateOffset int
		by         int64
	}
	type queryReq struct {
		dateOffset  int
		bucketCount int
	}
	tests := []struct {
		name            string
		inputs          []inputReq
		queries         []queryReq
		expectedResults []int64
	}{
		{
			name: "happy path",
			inputs: []inputReq{
				{0, 1},
			},
			queries: []queryReq{
				{0, 1},
			},
			expectedResults: []int64{
				1,
			},
		},
		{
			name: "happy path 2",
			inputs: []inputReq{
				{0, 1},
				{1, 1},
			},
			queries: []queryReq{
				{1, 1},
				{1, 2},
			},
			expectedResults: []int64{
				1, 2,
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
				{0, 1},
				{1, 2},
				{9, 10},
				{2, 1},
				{3, 3},
			},
			expectedResults: []int64{
				2, 3, 6, 0, 4,
			},
		},
	}

	counterToTest := map[string]func(DateRange)DateRangeCounter{
		"dateRange": func(dateRange DateRange) DateRangeCounter {
			return NewBasicDateCounter(dateRange, NewInMemoryBackend())
		},
		"intBacked": func(dateRange DateRange) DateRangeCounter {
			return NewIntBackedDateRange(NewBasicIntRangeCounter(), dateRange)
		},
	}
	rangeToTests := []DateRange{Seconds, Minute, Hour}
	for counterName, dateCounterFactory := range counterToTest {
		t.Run("counter " + counterName, func(t *testing.T) {
			for _, rangeToTest := range rangeToTests {
				t.Run("range " + fmt.Sprint(rangeToTest), func(t *testing.T) {
					for _, d := range tests {
						dateCounter := dateCounterFactory(rangeToTest)
						t.Run(d.name, func(t *testing.T) {
							ctx := context.Background()
							for _, inp := range d.inputs {
								err := dateCounter.Increment(ctx, rangeToTest.incrementDateForce(inp.dateOffset, baseDate), inp.by)
								assert.NoError(t, err)
							}

							outs := []int64{}
							for _, q := range d.queries {
								ans, err := dateCounter.QuerySum(ctx, rangeToTest.incrementDateForce(q.dateOffset, baseDate), q.bucketCount)
								assert.NoError(t, err)
								outs = append(outs, ans)
							}

							assert.EqualValues(t, d.expectedResults, outs)
						})
					}
				})
			}
		})
	}
}
