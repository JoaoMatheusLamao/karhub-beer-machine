package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"karhub-beer-machine/internal/application/beer"
	domain "karhub-beer-machine/internal/domain/beer"
	"karhub-beer-machine/internal/interfaces/http/dto"

	"github.com/google/uuid"
)

type BeerHandler struct {
	createUC   *beer.CreateBeerStyleUseCase
	updateUC   *beer.UpdateBeerStyleUseCase
	deleteUC   *beer.DeleteBeerStyleUseCase
	listUC     *beer.ListBeerStylesUseCase
	findBestUC *beer.FindBestBeerStyleUseCase
}

func NewBeerHandler(
	createUC *beer.CreateBeerStyleUseCase,
	updateUC *beer.UpdateBeerStyleUseCase,
	deleteUC *beer.DeleteBeerStyleUseCase,
	listUC *beer.ListBeerStylesUseCase,
	findBestUC *beer.FindBestBeerStyleUseCase,
) *BeerHandler {
	return &BeerHandler{
		createUC:   createUC,
		updateUC:   updateUC,
		deleteUC:   deleteUC,
		listUC:     listUC,
		findBestUC: findBestUC,
	}
}

/*
POST /beer-styles
*/
func (h *BeerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBeerStyleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err := h.createUC.Execute(beer.CreateBeerStyleInput{
		ID:      uuid.New().String(),
		Name:    req.Name,
		MinTemp: req.MinTemp,
		MaxTemp: req.MaxTemp,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

/*
PUT /beer-styles/{id}
*/
func (h *BeerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/beer-styles/")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateBeerStyleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err := h.updateUC.Execute(beer.UpdateBeerStyleInput{
		ID:      id,
		Name:    req.Name,
		MinTemp: req.MinTemp,
		MaxTemp: req.MaxTemp,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
DELETE /beer-styles/{id}
*/
func (h *BeerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/beer-styles/")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	if err := h.deleteUC.Execute(id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
GET /beer-styles
*/
func (h *BeerHandler) List(w http.ResponseWriter, _ *http.Request) {
	styles, err := h.listUC.Execute()
	if err != nil {
		h.handleError(w, err)
		return
	}

	resp := make([]dto.BeerStyleResponse, 0, len(styles))
	for _, s := range styles {
		resp = append(resp, dto.BeerStyleResponse{
			ID:      s.ID,
			Name:    s.Name,
			MinTemp: s.MinTemp,
			MaxTemp: s.MaxTemp,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

/*
POST /beer-styles/best
*/
func (h *BeerHandler) FindBest(w http.ResponseWriter, r *http.Request) {
	var req dto.FindBestBeerStyleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	out, err := h.findBestUC.Execute(
		r.Context(),
		beer.FindBestBeerStyleInput{Temperature: req.Temperature},
	)
	if err != nil {
		h.handleError(w, err)
		return
	}

	playlist := dto.PlaylistResponse{
		Name:   out.Playlist.Name,
		Tracks: []dto.TrackResponse{},
	}

	for _, t := range out.Playlist.Tracks {
		playlist.Tracks = append(playlist.Tracks, dto.TrackResponse{
			Name:   t.Name,
			Artist: t.Artist,
			Link:   t.Link,
		})
	}

	resp := dto.FindBestBeerStyleResponse{
		BeerStyle: out.BeerStyle,
		Playlist:  playlist,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *BeerHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidBeerStyle):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, domain.ErrBeerStyleNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrEmptyBeerStyleList):
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
