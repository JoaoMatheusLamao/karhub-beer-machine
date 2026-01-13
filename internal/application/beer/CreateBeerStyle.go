package beer

import (
	domain "karhub-beer-machine/internal/domain/beer"
)

// CreateBeerStyleInput represents the input for creating a beer style.
type CreateBeerStyleInput struct {
	ID      string
	Name    string
	MinTemp float64
	MaxTemp float64
}

// CreateBeerStyleUseCase handles creation of beer styles.
type CreateBeerStyleUseCase struct {
	repository domain.BeerStyleRepository
}

// NewCreateBeerStyleUseCase creates a new CreateBeerStyleUseCase.
func NewCreateBeerStyleUseCase(
	repository domain.BeerStyleRepository,
) *CreateBeerStyleUseCase {
	return &CreateBeerStyleUseCase{
		repository: repository,
	}
}

// Execute runs the use case.
func (uc *CreateBeerStyleUseCase) Execute(input CreateBeerStyleInput) error {
	style, err := domain.NewBeerStyle(
		input.ID,
		input.Name,
		input.MinTemp,
		input.MaxTemp,
	)
	if err != nil {
		return err
	}

	return uc.repository.Create(style)
}
