package sqlite

import (
	"context"
	"database/sql"

	s "github.com/devldavydov/myhealth/internal/storage"
)

func (r *StorageSQLite) GetTotalBurnedCal(ctx context.Context, userID int64, timestamp s.Timestamp) (float64, error) {
	var totalCal float64
	err := r.db.
		QueryRowContext(ctx, _sqlGetTotalBurnedCal, userID, timestamp).
		Scan(&totalCal)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, s.ErrTotalBurnedCalNotFound
		}
		return 0, err
	}

	return totalCal, nil
}

func (r *StorageSQLite) SetTotalBurnedCal(ctx context.Context, userID int64, timestamp s.Timestamp, totalCal float64) error {
	if totalCal <= 0 {
		return s.ErrDayTotalCalInvalid
	}

	_, err := r.db.ExecContext(ctx,
		_sqlSetTotalBurnedCal,
		userID,
		timestamp,
		totalCal,
	)

	return err
}

func (r *StorageSQLite) DeleteTotalBurnedCal(ctx context.Context, userID int64, timestamp s.Timestamp) error {
	_, err := r.db.ExecContext(ctx,
		_sqlDeleteTotalBurnedCal,
		userID,
		timestamp,
	)

	return err
}
