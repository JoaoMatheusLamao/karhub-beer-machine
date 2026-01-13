package beer_test

import (
	"testing"

	domain "karhub-beer-machine/internal/domain/beer"
)

func TestNewBeerStyle(t *testing.T) {
	tests := []struct {
		name      string
		styleName string
		minTemp   float64
		maxTemp   float64
		wantErr   bool
	}{
		{
			name:      "valid beer style",
			styleName: "IPA",
			minTemp:   -7,
			maxTemp:   10,
			wantErr:   false,
		},
		{
			name:      "empty name",
			styleName: "",
			minTemp:   -2,
			maxTemp:   4,
			wantErr:   true,
		},
		{
			name:      "invalid temperature range",
			styleName: "Pilsner",
			minTemp:   5,
			maxTemp:   -1,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewBeerStyle("id", tt.styleName, tt.minTemp, tt.maxTemp)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestBeerStyle_AverageTemperature(t *testing.T) {
	tests := []struct {
		name     string
		style    domain.BeerStyle
		expected float64
	}{
		{
			name: "average calculation",
			style: domain.BeerStyle{
				Name:    "Weissbier",
				MinTemp: -1,
				MaxTemp: 3,
			},
			expected: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			avg := tt.style.AverageTemperature()

			if avg != tt.expected {
				t.Errorf("expected average %.1f, got %.1f", tt.expected, avg)
			}
		})
	}
}

func TestBeerStyle_DistanceTo(t *testing.T) {
	tests := []struct {
		name       string
		style      domain.BeerStyle
		targetTemp float64
		expected   float64
	}{
		{
			name: "distance calculation",
			style: domain.BeerStyle{
				Name:    "IPA",
				MinTemp: -7,
				MaxTemp: 10,
			},
			targetTemp: -7,
			expected:   8.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance := tt.style.DistanceTo(tt.targetTemp)

			if distance != tt.expected {
				t.Errorf("expected distance %.1f, got %.1f", tt.expected, distance)
			}
		})
	}
}

func TestSelectBestStyle(t *testing.T) {
	tests := []struct {
		name        string
		styles      []domain.BeerStyle
		temperature float64
		wantStyle   string
		wantErr     bool
	}{
		{
			name: "select best style by distance",
			styles: []domain.BeerStyle{
				{Name: "Dunkel", MinTemp: -8, MaxTemp: 2},
				{Name: "Weissbier", MinTemp: -1, MaxTemp: 3},
			},
			temperature: -2,
			wantStyle:   "Dunkel",
			wantErr:     false,
		},
		{
			name: "tie break by alphabetical order",
			styles: []domain.BeerStyle{
				{Name: "Pilsens", MinTemp: -2, MaxTemp: 4}, // avg = 1
				{Name: "IPA", MinTemp: -1, MaxTemp: 3},     // avg = 1
			},
			temperature: 1,
			wantStyle:   "IPA",
			wantErr:     false,
		},
		{
			name:        "empty style list",
			styles:      []domain.BeerStyle{},
			temperature: 0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			best, err := domain.SelectBestStyle(tt.styles, tt.temperature)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tt.wantErr && best.Name != tt.wantStyle {
				t.Errorf("expected style %s, got %s", tt.wantStyle, best.Name)
			}
		})
	}
}
