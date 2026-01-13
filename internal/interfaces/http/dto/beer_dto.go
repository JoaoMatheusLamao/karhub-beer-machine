package dto

// ---------- Requests ----------

// CreateBeerStyleRequest represents the HTTP payload to create a beer style.
type CreateBeerStyleRequest struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	MinTemp float64 `json:"minTemp"`
	MaxTemp float64 `json:"maxTemp"`
}

// UpdateBeerStyleRequest represents the HTTP payload to update a beer style.
type UpdateBeerStyleRequest struct {
	Name    string  `json:"name"`
	MinTemp float64 `json:"minTemp"`
	MaxTemp float64 `json:"maxTemp"`
}

// FindBestBeerStyleRequest represents the HTTP payload to find the best beer style.
type FindBestBeerStyleRequest struct {
	Temperature float64 `json:"temperature"`
}

// ---------- Responses ----------

// BeerStyleResponse represents a beer style in HTTP responses.
type BeerStyleResponse struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	MinTemp float64 `json:"minTemp"`
	MaxTemp float64 `json:"maxTemp"`
}

// TrackResponse represents a track in a playlist response.
type TrackResponse struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Link   string `json:"link"`
}

// PlaylistResponse represents a playlist in HTTP responses.
type PlaylistResponse struct {
	Name   string          `json:"name"`
	Tracks []TrackResponse `json:"tracks"`
}

// FindBestBeerStyleResponse represents the response for best beer style endpoint.
type FindBestBeerStyleResponse struct {
	BeerStyle string           `json:"beerStyle"`
	Playlist  PlaylistResponse `json:"playlist"`
}
