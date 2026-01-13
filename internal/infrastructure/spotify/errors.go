package spotify

import "errors"

// ErrPlaylistNotFound is returned when no playlist is found for a beer style.
var ErrPlaylistNotFound = errors.New("playlist not found for beer style")
