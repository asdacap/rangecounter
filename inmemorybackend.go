package rangecounter

import "context"

type inMemoryBackend struct {
	store map[string]int64
}

func NewInMemoryBackend() Backend {
	return &inMemoryBackend{
		store: map[string]int64{},
	}
}

func (b inMemoryBackend) Query(ctx context.Context, keys []string) ([]int64, error) {
	results := make([]int64, 0, len(keys))
	for _, key := range keys {
		results = append(results, b.store[key])
	}
	return results, nil
}

func (b inMemoryBackend) Increment(ctx context.Context, keys []string, values []int64) error {
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		value := values[i]
		b.store[key] = b.store[key] + value
	}
	return nil
}
