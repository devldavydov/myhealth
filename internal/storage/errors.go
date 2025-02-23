package storage

import "errors"

var (
	// Meal
	ErrMealWrong = errors.New("wrong meal")

	// Weight
	ErrWeightInvalid = errors.New("invalid weight")

	// Sport
	ErrSportInvalid  = errors.New("invalid sport")
	ErrSportNotFound = errors.New("sport not found")
	ErrSportIsUsed   = errors.New("sport is used")

	// SportActivity
	ErrSportActivityInvalid = errors.New("invalid sport activity")

	// UserSettings
	ErrUserSettingsNotFound = errors.New("user settings not found")
	ErrUserSettingsInvalid  = errors.New("invalid user settings")

	// Common
	ErrEmptyResult = errors.New("empty result")
)
