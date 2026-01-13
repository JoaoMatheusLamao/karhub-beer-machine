package beer

import domain "karhub-beer-machine/internal/domain/beer"

// DeleteBeerStyleUseCase handles deletion of beer styles.
type DeleteBeerStyleUseCase struct {
	repository domain.BeerStyleRepository
}

// NewDeleteBeerStyleUseCase creates a new DeleteBeerStyleUseCase.
func NewDeleteBeerStyleUseCase(
	repository domain.BeerStyleRepository,
) *DeleteBeerStyleUseCase {
	return &DeleteBeerStyleUseCase{
		repository: repository,
	}
}

// Execute runs the use case.
func (uc *DeleteBeerStyleUseCase) Execute(id string) error {
	return uc.repository.Delete(id)
}
