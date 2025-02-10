package storage

import (
	"context"
)

type Storage interface {
	// Weight
	GetWeightList(ctx context.Context, userID int64, from, to Timestamp) ([]Weight, error)
	SetWeight(ctx context.Context, userID int64, weight *Weight) error
	DeleteWeight(ctx context.Context, userID int64, timestamp Timestamp) error

	// Backup/restore
	Backup(ctx context.Context) (*Backup, error)
	Restore(ctx context.Context, backup *Backup) error

	Close() error
}
