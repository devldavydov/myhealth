package sqlite

import (
	"context"
	"database/sql"
	"errors"

	s "github.com/devldavydov/myhealth/internal/storage"
	gsql "github.com/mattn/go-sqlite3"
)

//
// Medicine.
//

func (r *StorageSQLite) GetMedicine(ctx context.Context, userID int64, key string) (*s.Medicine, error) {
	var m s.Medicine
	err := r.db.
		QueryRowContext(ctx, _sqlGetMedicine, userID, key).
		Scan(&m.Key, &m.Name, &m.Comment, &m.Unit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, s.ErrMedicineNotFound
		}
		return nil, err
	}

	return &m, nil
}

func (r *StorageSQLite) GetMedicineList(ctx context.Context, userID int64) ([]s.Medicine, error) {
	rows, err := r.db.QueryContext(ctx, _sqlGetMedicineList, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []s.Medicine{}
	for rows.Next() {
		var m s.Medicine
		err = rows.Scan(&m.Key, &m.Name, &m.Comment, &m.Unit)
		if err != nil {
			return nil, err
		}

		list = append(list, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, s.ErrEmptyResult
	}

	return list, nil
}

func (r *StorageSQLite) SetMedicine(ctx context.Context, userID int64, m *s.Medicine) error {
	if !m.Validate() {
		return s.ErrMedicineInvalid
	}

	_, err := r.db.ExecContext(ctx, _sqlSetMedicine, userID, m.Key, m.Name, m.Comment, m.Unit)
	return err
}

func (r *StorageSQLite) DeleteMedicine(ctx context.Context, userID int64, key string) error {
	_, err := r.db.ExecContext(ctx, _sqlDeleteMedicine, userID, key)
	if err != nil {
		var errSql gsql.Error
		if errors.As(err, &errSql) && errSql.Error() == _errForeignKey {
			return s.ErrMedicineIsUsed
		}
		return err
	}

	return nil
}

//
// MedicineIndicator.
//

func (r *StorageSQLite) SetMedicineIndicator(ctx context.Context, userID int64, mi *s.MedicineIndicator) error {
	if !mi.Validate() {
		return s.ErrMedicineIndicatorInvalid
	}

	_, err := r.db.ExecContext(ctx, _sqlSetMedicineIndicator, userID, mi.Timestamp, mi.MedicineKey, mi.Value)
	if err != nil {
		var errSql gsql.Error
		if errors.As(err, &errSql) && errSql.Error() == _errForeignKey {
			return s.ErrMedicineNotFound
		}
		return err
	}

	return nil
}

func (r *StorageSQLite) DeleteMedicineIndicator(ctx context.Context, userID int64, timestamp s.Timestamp, medicine_key string) error {
	_, err := r.db.ExecContext(ctx, _sqlDeleteMedicineIndicator, userID, timestamp, medicine_key)
	return err
}

func (r *StorageSQLite) GetMedicineIndicatorReport(ctx context.Context, userID int64, from, to s.Timestamp) ([]s.MedicineIndicatorReport, error) {
	rows, err := r.db.QueryContext(ctx, _sqlGetMedicineIndicatorReport, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []s.MedicineIndicatorReport{}
	for rows.Next() {
		var mr s.MedicineIndicatorReport

		err = rows.Scan(&mr.Timestamp, &mr.MedicineName, &mr.Value)
		if err != nil {
			return nil, err
		}

		list = append(list, mr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, s.ErrEmptyResult
	}

	return list, nil
}
