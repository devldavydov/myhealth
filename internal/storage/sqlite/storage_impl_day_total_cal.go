package sqlite

import (
	"context"
	"database/sql"

	s "github.com/devldavydov/myhealth/internal/storage"
)

func (r *StorageSQLite) GetDayTotalCal(ctx context.Context, userID int64, timestamp s.Timestamp) (float64, error) {
	var totalCal float64
	err := r.db.
		QueryRowContext(ctx, _sqlGetDayTotalCal, userID, timestamp).
		Scan(&totalCal)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, s.ErrDayTotalCalNotFound
		}
		return 0, err
	}

	return totalCal, nil
}

func (r *StorageSQLite) SetDayTotalCal(ctx context.Context, userID int64, timestamp s.Timestamp, totalCal float64) error {
	if totalCal <= 0 {
		return s.ErrDayTotalCalInvalid
	}

	_, err := r.db.ExecContext(ctx,
		_sqlSetDayTotalCal,
		userID,
		timestamp,
		totalCal,
	)

	return err
}

func (r *StorageSQLite) DeleteDayTotalCal(ctx context.Context, userID int64, timestamp s.Timestamp) error {
	_, err := r.db.ExecContext(ctx,
		_sqlDeleteDayTotalCal,
		userID,
		timestamp,
	)

	return err
}
