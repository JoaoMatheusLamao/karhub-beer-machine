package cache

import (
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

// RistrettoCache is a generic in-memory cache wrapper around ristretto v2.
type RistrettoCache[K ristretto.Key, V any] struct {
	cache *ristretto.Cache[K, V]
}

// NewRistrettoCache creates a new Ristretto cache instance.
func NewRistrettoCache[K ristretto.Key, V any](
	numCounters int64,
	maxCost int64,
) (*RistrettoCache[K, V], error) {

	c, err := ristretto.NewCache[K, V](&ristretto.Config[K, V]{
		NumCounters: numCounters,
		MaxCost:     maxCost,
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}

	return &RistrettoCache[K, V]{cache: c}, nil
}

func (r *RistrettoCache[K, V]) Get(key K) (V, bool) {
	return r.cache.Get(key)
}

func (r *RistrettoCache[K, V]) SetWithTTL(
	key K,
	value V,
	ttl time.Duration,
) bool {
	return r.cache.SetWithTTL(key, value, 1, ttl)
}
