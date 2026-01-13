package memory_test

import (
	"testing"

	domain "karhub-beer-machine/internal/domain/beer"
	"karhub-beer-machine/internal/infrastructure/persistence/memory"
)

func TestBeerStyleRepository_CreateAndFindByID(t *testing.T) {
	repo, err := memory.NewBeerStyleRepository()
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	style := domain.BeerStyle{
		ID:      "1",
		Name:    "IPA",
		MinTemp: -7,
		MaxTemp: 10,
	}

	if err := repo.Create(style); err != nil {
		t.Fatalf("unexpected error on create: %v", err)
	}

	got, err := repo.FindByID("1")
	if err != nil {
		t.Fatalf("unexpected error on find: %v", err)
	}

	if got.Name != style.Name {
		t.Errorf("expected %s, got %s", style.Name, got.Name)
	}
}

func TestBeerStyleRepository_Update(t *testing.T) {
	repo, _ := memory.NewBeerStyleRepository()

	style := domain.BeerStyle{
		ID:      "1",
		Name:    "IPA",
		MinTemp: -7,
		MaxTemp: 10,
	}

	_ = repo.Create(style)

	style.Name = "Imperial IPA"

	if err := repo.Update(style); err != nil {
		t.Fatalf("unexpected error on update: %v", err)
	}

	got, _ := repo.FindByID("1")

	if got.Name != "Imperial IPA" {
		t.Errorf("expected updated name, got %s", got.Name)
	}
}

func TestBeerStyleRepository_Delete(t *testing.T) {
	repo, _ := memory.NewBeerStyleRepository()

	style := domain.BeerStyle{
		ID:      "1",
		Name:    "Stout",
		MinTemp: -5,
		MaxTemp: 5,
	}

	_ = repo.Create(style)

	if err := repo.Delete("1"); err != nil {
		t.Fatalf("unexpected error on delete: %v", err)
	}

	_, err := repo.FindByID("1")
	if err != domain.ErrBeerStyleNotFound {
		t.Errorf("expected ErrBeerStyleNotFound, got %v", err)
	}
}

func TestBeerStyleRepository_FindAll(t *testing.T) {
	repo, _ := memory.NewBeerStyleRepository()

	styles := []domain.BeerStyle{
		{ID: "1", Name: "IPA", MinTemp: -7, MaxTemp: 10},
		{ID: "2", Name: "Pilsner", MinTemp: -2, MaxTemp: 4},
	}

	for _, s := range styles {
		_ = repo.Create(s)
	}

	all, err := repo.FindAll()
	if err != nil {
		t.Fatalf("unexpected error on find all: %v", err)
	}

	if len(all) != len(styles) {
		t.Errorf("expected %d styles, got %d", len(styles), len(all))
	}
}

func TestBeerStyleRepository_UpdateNotFound(t *testing.T) {
	repo, _ := memory.NewBeerStyleRepository()

	err := repo.Update(domain.BeerStyle{
		ID:   "missing",
		Name: "Ghost Beer",
	})

	if err != domain.ErrBeerStyleNotFound {
		t.Errorf("expected ErrBeerStyleNotFound, got %v", err)
	}
}

func TestBeerStyleRepository_DeleteNotFound(t *testing.T) {
	repo, _ := memory.NewBeerStyleRepository()

	err := repo.Delete("missing")

	if err != domain.ErrBeerStyleNotFound {
		t.Errorf("expected ErrBeerStyleNotFound, got %v", err)
	}
}
