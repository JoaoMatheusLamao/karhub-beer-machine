package beer_test

import (
	"context"
	"errors"
	"testing"

	"karhub-beer-machine/internal/application/beer"
	domain "karhub-beer-machine/internal/domain/beer"
)

type beerStyleRepositoryMock struct {
	styles []domain.BeerStyle
	err    error
}

func (m *beerStyleRepositoryMock) Create(style domain.BeerStyle) error {
	return nil
}

func (m *beerStyleRepositoryMock) Update(style domain.BeerStyle) error {
	return nil
}

func (m *beerStyleRepositoryMock) Delete(id string) error {
	return nil
}

func (m *beerStyleRepositoryMock) FindByID(id string) (domain.BeerStyle, error) {
	return domain.BeerStyle{}, nil
}

func (m *beerStyleRepositoryMock) FindAll() ([]domain.BeerStyle, error) {
	return m.styles, m.err
}

type spotifyGatewayMock struct {
	playlist beer.Playlist
	err      error
}

func (m *spotifyGatewayMock) FindPlaylistByStyle(
	ctx context.Context,
	styleName string,
) (beer.Playlist, error) {
	return m.playlist, m.err
}

/*
	TESTS
*/

func TestNewFindBestBeerStyleUseCase(t *testing.T) {
	tests := []struct {
		name       string
		repository domain.BeerStyleRepository
		spotify    beer.SpotifyGateway
		wantNil    bool
	}{
		{
			name:       "valid dependencies",
			repository: &beerStyleRepositoryMock{},
			spotify:    &spotifyGatewayMock{},
			wantNil:    false,
		},
		{
			name:       "nil repository",
			repository: nil,
			spotify:    &spotifyGatewayMock{},
			wantNil:    false,
		},
		{
			name:       "nil spotify gateway",
			repository: &beerStyleRepositoryMock{},
			spotify:    nil,
			wantNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := beer.NewFindBestBeerStyleUseCase(tt.repository, tt.spotify)

			if tt.wantNil && got != nil {
				t.Errorf("expected nil use case, got %v", got)
			}

			if !tt.wantNil && got == nil {
				t.Errorf("expected non-nil use case, got nil")
			}
		})
	}
}

func TestFindBestBeerStyleUseCase_Execute(t *testing.T) {
	tests := []struct {
		name        string
		repository  domain.BeerStyleRepository
		spotify     beer.SpotifyGateway
		temperature float64
		wantStyle   string
		wantErr     bool
	}{
		{
			name: "repository error",
			repository: &beerStyleRepositoryMock{
				err: errors.New("db error"),
			},
			spotify:     &spotifyGatewayMock{},
			temperature: 0,
			wantErr:     true,
		},
		{
			name: "spotify error",
			repository: &beerStyleRepositoryMock{
				styles: []domain.BeerStyle{
					{Name: "IPA", MinTemp: -7, MaxTemp: 10},
				},
			},
			spotify: &spotifyGatewayMock{
				err: errors.New("playlist not found"),
			},
			temperature: -5,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := beer.NewFindBestBeerStyleUseCase(tt.repository, tt.spotify)

			output, err := useCase.Execute(
				context.Background(),
				beer.FindBestBeerStyleInput{Temperature: tt.temperature},
			)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tt.wantErr && output.BeerStyle != tt.wantStyle {
				t.Errorf("expected beer style %s, got %s", tt.wantStyle, output.BeerStyle)
			}
		})
	}
}
