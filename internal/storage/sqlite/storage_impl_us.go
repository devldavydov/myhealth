package sqlite

import (
	"context"
	"database/sql"

	s "github.com/devldavydov/myhealth/internal/storage"
)

func (r *StorageSQLite) GetUserSettings(ctx context.Context, userID int64) (*s.UserSettings, error) {
	var us s.UserSettings
	err := r.db.
		QueryRowContext(ctx, _sqlGetUserSettings, userID).
		Scan(&us.CalLimit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, s.ErrUserSettingsNotFound
		}
		return nil, err
	}

	return &us, nil
}

func (r *StorageSQLite) SetUserSettings(ctx context.Context, userID int64, us *s.UserSettings) error {
	if !us.Validate() {
		return s.ErrUserSettingsInvalid
	}

	_, err := r.db.ExecContext(ctx, _sqlSetUserSettings, userID, us.CalLimit)
	return err
}
