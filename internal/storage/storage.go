package storage

import (
	"context"
	"time"
)

const (
	StorageOperationTimeout = 15 * time.Second
	StorageRestoreTimeout   = 1 * time.Minute
)

type Storage interface {
	// Weight
	GetWeightList(ctx context.Context, userID int64, from, to Timestamp) ([]Weight, error)
	SetWeight(ctx context.Context, userID int64, weight *Weight) error
	DeleteWeight(ctx context.Context, userID int64, timestamp Timestamp) error

	// Food
	GetFood(ctx context.Context, userID int64, key string) (*Food, error)
	SetFood(ctx context.Context, userID int64, food *Food) error
	GetFoodList(ctx context.Context, userID int64) ([]Food, error)
	FindFood(ctx context.Context, userID int64, pattern string) ([]Food, error)
	DeleteFood(ctx context.Context, userID int64, key string) error

	// Bundle
	SetBundle(ctx context.Context, userID int64, bndl *Bundle, checkDeps bool) error
	GetBundle(ctx context.Context, userID int64, key string) (*Bundle, error)
	GetBundleList(ctx context.Context, userID int64) ([]Bundle, error)
	DeleteBundle(ctx context.Context, userID int64, key string) error

	// Journal
	SetJournal(ctx context.Context, userID int64, journal *Journal) error
	SetJournalBundle(ctx context.Context, userID int64, timestamp Timestamp, meal Meal, bndlKey string) error
	DeleteJournal(ctx context.Context, userID int64, timestamp Timestamp, meal Meal, foodkey string) error
	DeleteJournalMeal(ctx context.Context, userID int64, timestamp Timestamp, meal Meal) error
	DelJournalBundle(ctx context.Context, userID int64, timestamp Timestamp, meal Meal, bndlKey string) error
	GetJournalReport(ctx context.Context, userID int64, from, to Timestamp) ([]JournalReport, error)
	CopyJournal(ctx context.Context, userID int64, from Timestamp, mealFrom Meal, to Timestamp, mealTo Meal) (int, error)
	GetJournalFoodStat(ctx context.Context, userID int64, foodkey string) (*JournalFoodStat, error)

	// Sport
	GetSport(ctx context.Context, userID int64, key string) (*Sport, error)
	GetSportList(ctx context.Context, userID int64) ([]Sport, error)
	SetSport(ctx context.Context, userID int64, sp *Sport) error
	DeleteSport(ctx context.Context, userID int64, key string) error

	// SportActivity
	SetSportActivity(ctx context.Context, userID int64, sa *SportActivity) error
	DeleteSportActivity(ctx context.Context, userID int64, timestamp Timestamp, sport_key string) error
	GetSportActivityReport(ctx context.Context, userID int64, from, to Timestamp) ([]SportActivityReport, error)

	// Medicine
	GetMedicine(ctx context.Context, userID int64, key string) (*Medicine, error)
	GetMedicineList(ctx context.Context, userID int64) ([]Medicine, error)
	SetMedicine(ctx context.Context, userID int64, m *Medicine) error
	DeleteMedicine(ctx context.Context, userID int64, key string) error

	// MedicineIndicator
	SetMedicineIndicator(ctx context.Context, userID int64, mi *MedicineIndicator) error
	DeleteMedicineIndicator(ctx context.Context, userID int64, timestamp Timestamp, medicine_key string) error
	GetMedicineIndicatorReport(ctx context.Context, userID int64, from, to Timestamp) ([]MedicineIndicatorReport, error)

	// UserSettings
	GetUserSettings(ctx context.Context, userID int64) (*UserSettings, error)
	SetUserSettings(ctx context.Context, userID int64, us *UserSettings) error

	// TotalBurnedCal
	GetTotalBurnedCal(ctx context.Context, userID int64, timestamp Timestamp) (float64, error)
	SetTotalBurnedCal(ctx context.Context, userID int64, timestamp Timestamp, totalCal float64) error
	DeleteTotalBurnedCal(ctx context.Context, userID int64, timestamp Timestamp) error

	// Backup/restore
	Backup(ctx context.Context) (*Backup, error)
	Restore(ctx context.Context, backup *Backup) error

	Close() error
}
