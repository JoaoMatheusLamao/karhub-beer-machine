package beer

import (
	domain "karhub-beer-machine/internal/domain/beer"
)

// UpdateBeerStyleInput represents the input for updating a beer style.
type UpdateBeerStyleInput struct {
	ID      string
	Name    string
	MinTemp float64
	MaxTemp float64
}

// UpdateBeerStyleUseCase handles updating beer styles.
type UpdateBeerStyleUseCase struct {
	repository domain.BeerStyleRepository
}

// NewUpdateBeerStyleUseCase creates a new UpdateBeerStyleUseCase.
func NewUpdateBeerStyleUseCase(
	repository domain.BeerStyleRepository,
) *UpdateBeerStyleUseCase {
	return &UpdateBeerStyleUseCase{
		repository: repository,
	}
}

// Execute runs the use case.
func (uc *UpdateBeerStyleUseCase) Execute(input UpdateBeerStyleInput) error {
	style, err := domain.NewBeerStyle(
		input.ID,
		input.Name,
		input.MinTemp,
		input.MaxTemp,
	)
	if err != nil {
		return err
	}

	return uc.repository.Update(style)
}
