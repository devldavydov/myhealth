package storage

import "errors"

var (
	// Meal
	ErrMealWrong = errors.New("wrong meal")

	// Weight
	ErrWeightInvalid = errors.New("invalid weight")

	// Food
	ErrFoodNotFound = errors.New("food not found")
	ErrFoodInvalid  = errors.New("invalid food")
	ErrFoodIsUsed   = errors.New("food is used")

	// Bundle
	ErrBundleNotFound          = errors.New("bundle not found")
	ErrBundleInvalid           = errors.New("invalid bundle")
	ErrBundleDepFoodNotFound   = errors.New("dependent food not found")
	ErrBundleDepBundleNotFound = errors.New("dependent bundle not found")
	ErrBundleDepRecursive      = errors.New("dependent recursive bundle not allowed")
	ErrBundleIsUsed            = errors.New("bundle is used")

	// Journal
	ErrJournalInvalid = errors.New("journal invalid")

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
