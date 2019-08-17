package rangecounter

import "context"

type benchmarkingInMemoryBackend struct {
	store               map[string]int64
	queryCall           int
	queryCallFactor     int
	queryKeyTouched     int
	queryKeyFactor      int
	incrementCall       int
	incrementCallFactor int
	incrementKeyTouched int
	incrementKeyFactor  int
}

func NewBenchmarkBackend() Backend {
	return &benchmarkingInMemoryBackend{
		store:               map[string]int64{},
		queryCallFactor:     0,
		queryKeyFactor:      0,
		incrementCallFactor: 0,
		incrementKeyFactor:  0,
	}
}

func (b benchmarkingInMemoryBackend) Query(ctx context.Context, keys []string) ([]int64, error) {
	b.queryCall++
	b.queryKeyTouched += len(keys)

	for i := 0; i < b.queryCallFactor; i++ {
		load()
	}
	results := make([]int64, 0, len(keys))
	for _, key := range keys {
		for i := 0; i < b.queryKeyFactor; i++ {
			load()
		}
		results = append(results, b.store[key])
	}
	return results, nil
}

func (b benchmarkingInMemoryBackend) Increment(ctx context.Context, keys []string, values []int64) error {
	b.incrementCall++
	b.incrementKeyTouched += len(keys)

	for i := 0; i < b.incrementCallFactor; i++ {
		load()
	}
	for i := 0; i < len(keys); i++ {
		for i := 0; i < b.incrementKeyFactor; i++ {
			load()
		}
		key := keys[i]
		value := values[i]
		b.store[key] = b.store[key] + value
	}
	return nil
}

var something = 0

func load() {
	something++
}
