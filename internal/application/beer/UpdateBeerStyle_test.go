package beer_test

import (
	"testing"

	"karhub-beer-machine/internal/application/beer"
)

func TestUpdateBeerStyleUseCase(t *testing.T) {
	tests := []struct {
		name    string
		input   beer.UpdateBeerStyleInput
		wantErr bool
	}{
		{
			name: "valid update",
			input: beer.UpdateBeerStyleInput{
				ID:      "1",
				Name:    "Imperial IPA",
				MinTemp: -8,
				MaxTemp: 12,
			},
			wantErr: false,
		},
		{
			name: "invalid update",
			input: beer.UpdateBeerStyleInput{
				ID:      "2",
				Name:    "",
				MinTemp: 5,
				MaxTemp: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &beerStyleRepoMock{}
			uc := beer.NewUpdateBeerStyleUseCase(repo)

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
