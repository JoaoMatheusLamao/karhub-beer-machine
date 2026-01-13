package beer

import (
	"context"

	domain "karhub-beer-machine/internal/domain/beer"
)

// Playlist represents a simplified playlist model returned by the Spotify gateway.
type Playlist struct {
	Name   string
	Tracks []Track
}

// Track represents a music track inside a playlist.
type Track struct {
	Name   string
	Artist string
	Link   string
}

// SpotifyGateway defines the contract to retrieve playlists by beer style name.
// This is an application-level port.
type SpotifyGateway interface {
	FindPlaylistByStyle(ctx context.Context, styleName string) (Playlist, error)
}

// FindBestBeerStyleInput represents the input for the use case.
type FindBestBeerStyleInput struct {
	Temperature float64
}

// FindBestBeerStyleOutput represents the output of the use case.
type FindBestBeerStyleOutput struct {
	BeerStyle string
	Playlist  Playlist
}

// FindBestBeerStyleUseCase orchestrates the process of selecting the best beer style
// for a given temperature and retrieving a related playlist.
type FindBestBeerStyleUseCase struct {
	repository domain.BeerStyleRepository
	spotify    SpotifyGateway
}

// NewFindBestBeerStyleUseCase creates a new instance of the use case.
func NewFindBestBeerStyleUseCase(
	repository domain.BeerStyleRepository,
	spotify SpotifyGateway,
) *FindBestBeerStyleUseCase {
	return &FindBestBeerStyleUseCase{
		repository: repository,
		spotify:    spotify,
	}
}

// Execute runs the use case.
func (uc *FindBestBeerStyleUseCase) Execute(
	ctx context.Context,
	input FindBestBeerStyleInput,
) (FindBestBeerStyleOutput, error) {
	styles, err := uc.repository.FindAll()
	if err != nil {
		return FindBestBeerStyleOutput{}, err
	}

	bestStyle, err := domain.SelectBestStyle(styles, input.Temperature)
	if err != nil {
		return FindBestBeerStyleOutput{}, err
	}

	playlist, err := uc.spotify.FindPlaylistByStyle(ctx, bestStyle.Name)
	if err != nil {
		return FindBestBeerStyleOutput{}, err
	}

	return FindBestBeerStyleOutput{
		BeerStyle: bestStyle.Name,
		Playlist:  playlist,
	}, nil
}
