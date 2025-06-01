package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	s "github.com/devldavydov/myhealth/internal/storage"
)

func (r *StorageSQLite) SetBundle(ctx context.Context, userID int64, bndl *s.Bundle, checkDeps bool) error {
	if !bndl.Validate() {
		return s.ErrBundleInvalid
	}

	if checkDeps {
		for k, v := range bndl.Data {
			if v == 0 {
				if k == bndl.Key {
					return s.ErrBundleDepRecursive
				}

				_, err := r.GetBundle(ctx, userID, k)
				if err != nil {
					if errors.Is(err, s.ErrBundleNotFound) {
						return s.ErrBundleDepBundleNotFound
					}

					return err
				}
			} else {
				_, err := r.GetFood(ctx, userID, k)
				if err != nil {
					if errors.Is(err, s.ErrFoodNotFound) {
						return s.ErrBundleDepFoodNotFound
					}

					return err
				}
			}
		}
	}

	bData, err := json.Marshal(&bndl.Data)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx,
		_sqlSetBundle,
		userID,
		bndl.Key,
		string(bData),
	)
	return err
}

func (r *StorageSQLite) GetBundle(ctx context.Context, userID int64, key string) (*s.Bundle, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bndl, err := getBundle(ctx, tx, userID, key)
	if err != nil {
		return nil, err
	}

	return bndl, tx.Commit()
}

func getBundle(ctx context.Context, tx *sql.Tx, userID int64, key string) (*s.Bundle, error) {
	var b s.Bundle
	var bData string
	err := tx.
		QueryRowContext(ctx, _sqlGetBundle, userID, key).
		Scan(&b.Key, &bData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, s.ErrBundleNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(bData), &b.Data); err != nil {
		return nil, err
	}

	return &b, nil
}

func (r *StorageSQLite) GetBundleList(ctx context.Context, userID int64) ([]s.Bundle, error) {
	rows, err := r.db.QueryContext(ctx, _sqlGetBundleList, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []s.Bundle{}
	for rows.Next() {
		var b s.Bundle
		var bData string
		err = rows.Scan(&b.Key, &bData)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(bData), &b.Data); err != nil {
			return nil, err
		}

		list = append(list, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, s.ErrEmptyResult
	}

	return list, nil
}

func (r *StorageSQLite) DeleteBundle(ctx context.Context, userID int64, key string) error {
	bndlList, err := r.GetBundleList(ctx, userID)
	if err != nil {
		if errors.Is(err, s.ErrEmptyResult) {
			return nil
		}

		return err
	}

	for _, bndl := range bndlList {
		for k, v := range bndl.Data {
			if v == 0 && k == key {
				return s.ErrBundleIsUsed
			}
		}
	}

	_, err = r.db.ExecContext(ctx, _sqlDeleteBundle, userID, key)

	return err
}
