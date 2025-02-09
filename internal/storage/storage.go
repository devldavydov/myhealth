package storage

import (
	"context"
	"time"
)

type Storage interface {
	// Weight
	GetWeightList(ctx context.Context, userID int64, from, to time.Time) ([]Weight, error)
	SetWeight(ctx context.Context, userID int64, weight *Weight) error
	DeleteWeight(ctx context.Context, userID int64, timestamp time.Time) error

	Close() error
}
