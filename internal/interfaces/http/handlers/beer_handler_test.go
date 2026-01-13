package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"karhub-beer-machine/internal/application/beer"
	domain "karhub-beer-machine/internal/domain/beer"
	"karhub-beer-machine/internal/infrastructure/persistence/memory"
	httpapi "karhub-beer-machine/internal/interfaces/http"
	"karhub-beer-machine/internal/interfaces/http/handlers"
)

/*
	Spotify mock
*/

type spotifyMock struct {
	playlist beer.Playlist
	err      error
}

func (s *spotifyMock) FindPlaylistByStyle(
	_ context.Context,
	styleName string,
) (beer.Playlist, error) {
	return s.playlist, s.err
}

/*
	Helper to build test server
*/

func setupServer(t *testing.T) *httptest.Server {
	t.Helper()

	repo, err := memory.NewBeerStyleRepository()
	if err != nil {
		t.Fatalf("failed to create repo: %v", err)
	}

	// seed data
	_ = repo.Create(domain.BeerStyle{ID: "1", Name: "Dunkel", MinTemp: -8, MaxTemp: 2})
	_ = repo.Create(domain.BeerStyle{ID: "2", Name: "IPA", MinTemp: -7, MaxTemp: 10})

	spotify := &spotifyMock{
		playlist: beer.Playlist{Name: "IPA Party"},
	}

	// use cases
	createUC := beer.NewCreateBeerStyleUseCase(repo)
	updateUC := beer.NewUpdateBeerStyleUseCase(repo)
	deleteUC := beer.NewDeleteBeerStyleUseCase(repo)
	listUC := beer.NewListBeerStylesUseCase(repo)
	findBestUC := beer.NewFindBestBeerStyleUseCase(repo, spotify)

	handler := handlers.NewBeerHandler(
		createUC,
		updateUC,
		deleteUC,
		listUC,
		findBestUC,
	)

	mux := http.NewServeMux()
	httpapi.RegisterRoutes(mux, handler)

	return httptest.NewServer(mux)
}

/*
	TESTS
*/

func TestFindBestBeerStyleHTTP(t *testing.T) {
	server := setupServer(t)
	defer server.Close()

	tests := []struct {
		name           string
		body           any
		wantStatusCode int
		wantBeerStyle  string
	}{
		{
			name: "success",
			body: map[string]any{
				"temperature": -7,
			},
			wantStatusCode: http.StatusOK,
			wantBeerStyle:  "Dunkel",
		},
		{
			name:           "invalid body",
			body:           "invalid-json",
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			switch v := tt.body.(type) {
			case string:
				buf.WriteString(v)
			default:
				_ = json.NewEncoder(&buf).Encode(v)
			}

			resp, err := http.Post(
				server.URL+"/beer-styles/best",
				"application/json",
				&buf,
			)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatusCode {
				t.Fatalf(
					"expected status %d, got %d",
					tt.wantStatusCode,
					resp.StatusCode,
				)
			}

			if tt.wantStatusCode == http.StatusOK {
				var out struct {
					BeerStyle string `json:"beerStyle"`
				}

				if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if out.BeerStyle != tt.wantBeerStyle {
					t.Errorf(
						"expected beerStyle %s, got %s",
						tt.wantBeerStyle,
						out.BeerStyle,
					)
				}
			}
		})
	}
}
