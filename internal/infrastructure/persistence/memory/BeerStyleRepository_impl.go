package memory

import (
	"sync"

	"github.com/dgraph-io/ristretto/v2"

	domain "karhub-beer-machine/internal/domain/beer"
)

// BeerStyleRepositoryImpl is an in-memory implementation of BeerStyleRepository
// backed by Ristretto v2 (generic).
//
// Key   -> string (BeerStyle.ID)
// Value -> domain.BeerStyle
type BeerStyleRepositoryImpl struct {
	cache *ristretto.Cache[string, domain.BeerStyle]

	// keys is required because Ristretto does NOT support iteration.
	// We keep track of inserted keys to support FindAll().
	mu   sync.RWMutex
	keys map[string]struct{}
}

// NewBeerStyleRepository creates a new in-memory BeerStyleRepository using Ristretto v2.
func NewBeerStyleRepository() (*BeerStyleRepositoryImpl, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, domain.BeerStyle]{
		NumCounters: 1e5,
		MaxCost:     1 << 20, // ~1MB
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}

	return &BeerStyleRepositoryImpl{
		cache: cache,
		keys:  make(map[string]struct{}),
	}, nil
}

// Create stores a new beer style.
func (r *BeerStyleRepositoryImpl) Create(style domain.BeerStyle) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cache.Set(style.ID, style, 1)
	r.cache.Wait()

	r.keys[style.ID] = struct{}{}
	return nil
}

// Update updates an existing beer style.
func (r *BeerStyleRepositoryImpl) Update(style domain.BeerStyle) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.keys[style.ID]; !found {
		return domain.ErrBeerStyleNotFound
	}

	r.cache.Set(style.ID, style, 1)
	r.cache.Wait()
	return nil
}

// Delete removes a beer style by ID.
func (r *BeerStyleRepositoryImpl) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.keys[id]; !found {
		return domain.ErrBeerStyleNotFound
	}

	r.cache.Del(id)
	delete(r.keys, id)
	return nil
}

// FindByID retrieves a beer style by ID.
func (r *BeerStyleRepositoryImpl) FindByID(id string) (domain.BeerStyle, error) {
	style, found := r.cache.Get(id)
	if !found {
		return domain.BeerStyle{}, domain.ErrBeerStyleNotFound
	}

	return style, nil
}

// FindAll retrieves all beer styles stored in memory.
func (r *BeerStyleRepositoryImpl) FindAll() ([]domain.BeerStyle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	styles := make([]domain.BeerStyle, 0, len(r.keys))

	for id := range r.keys {
		if style, found := r.cache.Get(id); found {
			styles = append(styles, style)
		}
	}

	return styles, nil
}
