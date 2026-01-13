package spotify

import (
	"context"
	"fmt"
	"time"

	"karhub-beer-machine/internal/application/beer"
	"karhub-beer-machine/internal/infrastructure/cache"
)

// CachedGateway decorates a SpotifyGateway with cache.
type CachedGateway struct {
	gateway beer.SpotifyGateway
	cache   *cache.RistrettoCache[string, beer.Playlist]
	ttl     time.Duration
}

// NewCachedSpotifyGateway creates a cached Spotify gateway.
func NewCachedSpotifyGateway(
	gateway beer.SpotifyGateway,
	cache *cache.RistrettoCache[string, beer.Playlist],
	ttl time.Duration,
) *CachedGateway {
	return &CachedGateway{
		gateway: gateway,
		cache:   cache,
		ttl:     ttl,
	}
}

func (c *CachedGateway) FindPlaylistByStyle(
	ctx context.Context,
	styleName string,
) (beer.Playlist, error) {

	key := fmt.Sprintf("spotify:playlist:%s", styleName)

	if playlist, found := c.cache.Get(key); found {
		return playlist, nil
	}

	playlist, err := c.gateway.FindPlaylistByStyle(ctx, styleName)
	if err != nil {
		return beer.Playlist{}, err
	}

	c.cache.SetWithTTL(key, playlist, c.ttl)
	return playlist, nil
}
