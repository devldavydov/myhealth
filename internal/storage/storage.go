package storage

import (
	"context"
	"time"
)

const StorageOperationTimeout = 15 * time.Second

type Storage interface {
	// Weight
	GetWeightList(ctx context.Context, userID int64, from, to Timestamp) ([]Weight, error)
	SetWeight(ctx context.Context, userID int64, weight *Weight) error
	DeleteWeight(ctx context.Context, userID int64, timestamp Timestamp) error

	// Sport
	GetSport(ctx context.Context, userID int64, key string) (*Sport, error)
	GetSportList(ctx context.Context, userID int64) ([]Sport, error)
	SetSport(ctx context.Context, userID int64, sp *Sport) error
	DeleteSport(ctx context.Context, userID int64, key string) error

	// SportActivity
	SetSportActivity(ctx context.Context, userID int64, sa *SportActivity) error
	DeleteSportActivity(ctx context.Context, userID int64, timestamp Timestamp, sport_key string) error
	GetSportActivityReport(ctx context.Context, userID int64, from, to Timestamp) ([]SportActivityReport, error)

	// Backup/restore
	Backup(ctx context.Context) (*Backup, error)
	Restore(ctx context.Context, backup *Backup) error

	Close() error
}
