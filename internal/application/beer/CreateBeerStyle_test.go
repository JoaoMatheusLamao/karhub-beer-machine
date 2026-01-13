package beer_test

import (
	"testing"

	"karhub-beer-machine/internal/application/beer"
	domain "karhub-beer-machine/internal/domain/beer"
)

type beerStyleRepoMock struct {
	created []domain.BeerStyle
	err     error
}

func (m *beerStyleRepoMock) Create(style domain.BeerStyle) error {
	m.created = append(m.created, style)
	return m.err
}

func (m *beerStyleRepoMock) Update(style domain.BeerStyle) error {
	return nil
}
func (m *beerStyleRepoMock) Delete(id string) error {
	return nil
}
func (m *beerStyleRepoMock) FindByID(id string) (domain.BeerStyle, error) {
	return domain.BeerStyle{}, nil
}
func (m *beerStyleRepoMock) FindAll() ([]domain.BeerStyle, error) {
	return nil, nil
}

func TestCreateBeerStyleUseCase(t *testing.T) {
	tests := []struct {
		name    string
		input   beer.CreateBeerStyleInput
		wantErr bool
	}{
		{
			name: "valid beer style",
			input: beer.CreateBeerStyleInput{
				ID:      "1",
				Name:    "IPA",
				MinTemp: -7,
				MaxTemp: 10,
			},
			wantErr: false,
		},
		{
			name: "invalid beer style",
			input: beer.CreateBeerStyleInput{
				ID:      "2",
				Name:    "",
				MinTemp: 1,
				MaxTemp: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &beerStyleRepoMock{}
			uc := beer.NewCreateBeerStyleUseCase(repo)

			err := uc.Execute(tt.input)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
