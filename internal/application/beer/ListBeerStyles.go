package beer

import domain "karhub-beer-machine/internal/domain/beer"

// ListBeerStylesUseCase handles listing all beer styles.
type ListBeerStylesUseCase struct {
	repository domain.BeerStyleRepository
}

// NewListBeerStylesUseCase creates a new ListBeerStylesUseCase.
func NewListBeerStylesUseCase(
	repository domain.BeerStyleRepository,
) *ListBeerStylesUseCase {
	return &ListBeerStylesUseCase{
		repository: repository,
	}
}

// Execute runs the use case.
func (uc *ListBeerStylesUseCase) Execute() ([]domain.BeerStyle, error) {
	return uc.repository.FindAll()
}
