package spotify

import (
	"context"

	"karhub-beer-machine/internal/application/beer"
)

// Stub implements beer.SpotifyGateway as a fallback.
type Stub struct{}

// NewSpotifyStub creates a Spotify stub gateway.
func NewSpotifyStub() *Stub {
	return &Stub{}
}

func (s *Stub) FindPlaylistByStyle(
	_ context.Context,
	styleName string,
) (beer.Playlist, error) {
	return beer.Playlist{
		Name: styleName + " Playlist (stub)",
		Tracks: []beer.Track{
			{
				Name:   "Stub Song",
				Artist: "Stub Artist",
				Link:   "https://open.spotify.com",
			},
		},
	}, nil
}
