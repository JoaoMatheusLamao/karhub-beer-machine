package beer

import (
	"math"
	"sort"
)

// BeerStyle represents a beer style and its ideal temperature range.
// This is a core domain entity and must not depend on external layers.
type BeerStyle struct {
	ID      string
	Name    string
	MinTemp float64
	MaxTemp float64
}

// NewBeerStyle creates a new BeerStyle ensuring domain invariants.
func NewBeerStyle(id, name string, minTemp, maxTemp float64) (BeerStyle, error) {
	if name == "" {
		return BeerStyle{}, ErrInvalidBeerStyle
	}

	if minTemp > maxTemp {
		return BeerStyle{}, ErrInvalidBeerStyle
	}

	return BeerStyle{
		ID:      id,
		Name:    name,
		MinTemp: minTemp,
		MaxTemp: maxTemp,
	}, nil
}

// AverageTemperature returns the average temperature of the beer style.
func (b BeerStyle) AverageTemperature() float64 {
	return (b.MinTemp + b.MaxTemp) / 2
}

// DistanceTo returns the absolute distance between the style average temperature
// and a target temperature.
func (b BeerStyle) DistanceTo(targetTemp float64) float64 {
	return math.Abs(b.AverageTemperature() - targetTemp)
}

// SelectBestStyle selects the most suitable beer style for a given temperature.
// Selection rules:
// 1. Choose the style whose average temperature is closest to the target.
// 2. In case of a tie, select the style by alphabetical order (lexicographical).
func SelectBestStyle(styles []BeerStyle, targetTemp float64) (BeerStyle, error) {
	if len(styles) == 0 {
		return BeerStyle{}, ErrEmptyBeerStyleList
	}

	type candidate struct {
		style    BeerStyle
		distance float64
	}

	candidates := make([]candidate, 0, len(styles))

	for _, style := range styles {
		candidates = append(candidates, candidate{
			style:    style,
			distance: style.DistanceTo(targetTemp),
		})
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].distance != candidates[j].distance {
			return candidates[i].distance < candidates[j].distance
		}
		return candidates[i].style.Name < candidates[j].style.Name
	})

	return candidates[0].style, nil
}
