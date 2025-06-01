package sqlite

import (
	"context"
	"database/sql"
)

type migration struct {
	mID int64
	mF  func(ctx context.Context, tx *sql.Tx) error
}

func (r *StorageSQLite) doMigrations() error {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseInitTimeout)
	defer cancel()

	// Create system table
	if _, err := r.db.ExecContext(ctx, _sqlCreateTableSystem); err != nil {
		return err
	}

	// Apply migrations after last
	lastMigrationID, err := r.getLastMigrationID(ctx)
	if err != nil {
		return err
	}

	for _, m := range getAllMigrations() {
		if m.mID <= lastMigrationID {
			continue
		}

		tx, err := r.db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		if err = m.mF(ctx, tx); err != nil {
			return err
		}

		if err = r.updateLastMigrationID(ctx, tx, m.mID); err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func (r *StorageSQLite) getLastMigrationID(ctx context.Context) (int64, error) {
	var migration_id int64
	err := r.db.QueryRowContext(ctx, _sqlGetLastMigrationID).Scan(&migration_id)
	switch {
	case err == sql.ErrNoRows:
		return 0, nil
	case err != nil:
		return 0, err
	}

	return migration_id, nil
}

func (r *StorageSQLite) updateLastMigrationID(ctx context.Context, tx *sql.Tx, migrationID int64) error {
	_, err := tx.ExecContext(ctx, _sqlUpdateLastMigrationID, migrationID)
	return err
}

func getAllMigrations() []migration {
	return []migration{
		{1, createInitialMigration},
		{2, createTableWeight},
		{3, createTableSport},
		{4, createTableSportActivity},
		{5, createTableUserSettings},
		{6, createTableFood},
		{7, createTableBundle},
		{8, createTableJournal},
		{9, createTableJournalIndexUserIDFoodKey},
		{10, createTableMedicine},
		{11, createTableMedicineIndicator},
	}
}

func createInitialMigration(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, _sqlInsertInitialMigrationID)
	return err
}

func createTableWeight(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, _sqlCreateTableWeight)
	return err
}

func createTableSport(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, _sqlCreateTableSport)
	return err
}

func createTableSportActivity(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, _sqlCreateTableSportActivity)
	return err
}

func createTableUserSettings(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, _sqlCreateTableUserSettings)
	return err
}

func createTableFood(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, _sqlCreateTableFood)
	return err
}

func createTableBundle(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, _sqlCreateTableBundle)
	return err
}

func createTableJournal(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, _sqlCreateTableJournal)
	return err
}

func createTableJournalIndexUserIDFoodKey(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, _sqlCreateTableJournalIndexUserIDFoodKey)
	return err
}

func createTableMedicine(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, _sqlCreateTableMedicine)
	return err
}

func createTableMedicineIndicator(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, _sqlCreateTableMedicineIndicator)
	return err
}
