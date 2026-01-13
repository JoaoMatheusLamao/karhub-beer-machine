package beer

// BeerStyleRepository defines the persistence contract for beer styles.
// This interface belongs to the domain layer and must be implemented
// by infrastructure adapters (e.g., memory, postgres).
type BeerStyleRepository interface {
	// Create persists a new beer style.
	Create(style BeerStyle) error

	// Update updates an existing beer style.
	Update(style BeerStyle) error

	// Delete removes a beer style by its identifier.
	Delete(id string) error

	// FindByID retrieves a beer style by its identifier.
	FindByID(id string) (BeerStyle, error)

	// FindAll retrieves all beer styles.
	FindAll() ([]BeerStyle, error)
}
