package beer

import "errors"

// Domain-level errors.
// These errors represent business rule violations and must not
// carry HTTP or infrastructure semantics.

var (
	// ErrInvalidBeerStyle is returned when a beer style violates domain invariants
	// such as empty name or invalid temperature range.
	ErrInvalidBeerStyle = errors.New("invalid beer style")

	// ErrBeerStyleNotFound is returned when a beer style cannot be found.
	ErrBeerStyleNotFound = errors.New("beer style not found")

	// ErrEmptyBeerStyleList is returned when no beer styles are available
	// to perform a selection.
	ErrEmptyBeerStyleList = errors.New("beer style list is empty")
)
