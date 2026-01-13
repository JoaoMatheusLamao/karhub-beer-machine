package beer_test

import (
	"testing"

	"karhub-beer-machine/internal/application/beer"
)

func TestDeleteBeerStyleUseCase(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "delete existing style",
			id:      "1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &beerStyleRepoMock{}
			uc := beer.NewDeleteBeerStyleUseCase(repo)

			err := uc.Execute(tt.id)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
	}
}
