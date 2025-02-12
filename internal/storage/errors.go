package storage

import "errors"

var (
	// Meal
	ErrMealWrong = errors.New("wrong meal")

	// Weight
	ErrWeightInvalid = errors.New("invalid weight")

	// Sport
	ErrSportInvalid = errors.New("invalid sport")

	// Common
	ErrEmptyResult = errors.New("empty result")
)
