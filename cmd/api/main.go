package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"karhub-beer-machine/internal/application/beer"
	domain "karhub-beer-machine/internal/domain/beer"
	cacheinfra "karhub-beer-machine/internal/infrastructure/cache"
	"karhub-beer-machine/internal/infrastructure/persistence/memory"
	spotifyinfra "karhub-beer-machine/internal/infrastructure/spotify"
	httpapi "karhub-beer-machine/internal/interfaces/http"
	"karhub-beer-machine/internal/interfaces/http/handlers"
)

func main() {
	ctx := context.Background()

	repo := mustCreateRepository()
	spotifyGateway := mustCreateSpotifyGateway(ctx)

	useCases := buildUseCases(repo, spotifyGateway)
	handler := buildHTTPHandler(useCases)

	server := buildHTTPServer(handler)

	log.Printf("HTTP server running on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

/*
	---------- Builders ----------
*/

func mustCreateRepository() domain.BeerStyleRepository {
	repo, err := memory.NewBeerStyleRepository()
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}
	return repo
}

func mustCreateSpotifyGateway(ctx context.Context) beer.SpotifyGateway {
	// Cache sempre existe, independente de Spotify real ou stub
	playlistCache, err := cacheinfra.NewRistrettoCache[string, beer.Playlist](
		1e5,   // counters
		1<<20, // ~1MB
	)
	if err != nil {
		log.Fatalf("failed to create cache: %v", err)
	}

	spotifyClient, err := spotifyinfra.NewSpotifyClient(ctx)
	if err != nil {
		log.Printf(
			"spotify unavailable (%v), falling back to stub",
			err,
		)

		stub := spotifyinfra.NewSpotifyStub()

		return spotifyinfra.NewCachedSpotifyGateway(
			stub,
			playlistCache,
			10*time.Minute,
		)
	}

	return spotifyinfra.NewCachedSpotifyGateway(
		spotifyClient,
		playlistCache,
		10*time.Minute,
	)
}

type useCases struct {
	create   *beer.CreateBeerStyleUseCase
	update   *beer.UpdateBeerStyleUseCase
	delete   *beer.DeleteBeerStyleUseCase
	list     *beer.ListBeerStylesUseCase
	findBest *beer.FindBestBeerStyleUseCase
}

func buildUseCases(
	repo domain.BeerStyleRepository,
	spotify beer.SpotifyGateway,
) useCases {
	return useCases{
		create:   beer.NewCreateBeerStyleUseCase(repo),
		update:   beer.NewUpdateBeerStyleUseCase(repo),
		delete:   beer.NewDeleteBeerStyleUseCase(repo),
		list:     beer.NewListBeerStylesUseCase(repo),
		findBest: beer.NewFindBestBeerStyleUseCase(repo, spotify),
	}
}

func buildHTTPHandler(uc useCases) *handlers.BeerHandler {
	return handlers.NewBeerHandler(
		uc.create,
		uc.update,
		uc.delete,
		uc.list,
		uc.findBest,
	)
}

func buildHTTPServer(handler *handlers.BeerHandler) *http.Server {
	mux := http.NewServeMux()
	httpapi.RegisterRoutes(mux, handler)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	return &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
