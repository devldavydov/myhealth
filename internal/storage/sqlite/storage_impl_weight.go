package sqlite

import (
	"context"

	s "github.com/devldavydov/myhealth/internal/storage"
)

func (r *StorageSQLite) GetWeightList(ctx context.Context, userID int64, from, to s.Timestamp) ([]s.Weight, error) {
	rows, err := r.db.QueryContext(ctx, _sqlGetWeightList, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []s.Weight{}
	for rows.Next() {
		var w s.Weight
		err = rows.Scan(&w.Timestamp, &w.Value)
		if err != nil {
			return nil, err
		}

		list = append(list, w)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, s.ErrEmptyResult
	}

	return list, nil
}

func (r *StorageSQLite) SetWeight(ctx context.Context, userID int64, weight *s.Weight) error {
	if !weight.Validate() {
		return s.ErrWeightInvalid
	}

	_, err := r.db.ExecContext(ctx, _sqlSetWeight, userID, weight.Timestamp, weight.Value)
	return err
}

func (r *StorageSQLite) DeleteWeight(ctx context.Context, userID int64, timestamp s.Timestamp) error {
	_, err := r.db.ExecContext(ctx, _sqlDeleteWeight, userID, timestamp)
	return err
}
