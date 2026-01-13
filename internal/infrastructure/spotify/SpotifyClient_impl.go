package spotify

import (
	"context"
	"errors"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"

	"karhub-beer-machine/internal/application/beer"
)

// Client implements beer.SpotifyGateway using Spotify Web API.
type Client struct {
	client *spotify.Client
}

// NewSpotifyClient creates a Spotify client using Client Credentials flow,
// following the official example from the spotify/v2 repository.
func NewSpotifyClient(ctx context.Context) (*Client, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, errors.New("spotify credentials not set")
	}

	// OAuth2 client credentials configuration
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	// Retrieve token
	token, err := config.Token(ctx)
	if err != nil {
		return nil, err
	}

	// Create HTTP client authorized with the token
	httpClient := spotifyauth.New().Client(ctx, token)

	// Create Spotify client
	client := spotify.New(
		httpClient,
		spotify.WithRetry(true),
	)

	return &Client{client: client}, nil
}

// FindPlaylistByStyle searches for a public playlist containing the beer style name.
func (c *Client) FindPlaylistByStyle(
	ctx context.Context,
	styleName string,
) (beer.Playlist, error) {

	results, err := c.client.Search(
		ctx,
		styleName,
		spotify.SearchTypePlaylist,
		spotify.Limit(1),
	)
	if err != nil {
		return beer.Playlist{}, err
	}

	if results.Playlists == nil || len(results.Playlists.Playlists) == 0 {
		return beer.Playlist{}, ErrPlaylistNotFound
	}

	pl := results.Playlists.Playlists[0]

	playlist := beer.Playlist{
		Name:   pl.Name,
		Tracks: []beer.Track{},
	}

	// Fetch playlist tracks (limited for performance)
	tracks, err := c.client.GetPlaylistItems(
		ctx,
		pl.ID,
		spotify.Limit(10),
	)
	if err != nil {
		// Playlist without tracks is still acceptable
		return playlist, nil
	}

	for _, item := range tracks.Items {
		if item.Track.Track == nil {
			continue
		}

		track := beer.Track{
			Name: item.Track.Track.Name,
			Link: item.Track.Track.ExternalURLs["spotify"],
		}

		if len(item.Track.Track.Artists) > 0 {
			track.Artist = item.Track.Track.Artists[0].Name
		}

		playlist.Tracks = append(playlist.Tracks, track)
	}

	return playlist, nil
}
