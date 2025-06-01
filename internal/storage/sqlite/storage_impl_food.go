package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	s "github.com/devldavydov/myhealth/internal/storage"
	gsql "github.com/mattn/go-sqlite3"
)

func (r *StorageSQLite) GetFood(ctx context.Context, userID int64, key string) (*s.Food, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	food, err := getFood(ctx, tx, userID, key)
	if err != nil {
		return nil, err
	}

	return food, tx.Commit()
}

func getFood(ctx context.Context, tx *sql.Tx, userID int64, key string) (*s.Food, error) {
	var f s.Food
	err := tx.
		QueryRowContext(ctx, _sqlGetFood, userID, key).
		Scan(&f.Key, &f.Name, &f.Brand, &f.Cal100, &f.Prot100, &f.Fat100, &f.Carb100, &f.Comment)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, s.ErrFoodNotFound
		}
		return nil, err
	}

	return &f, nil
}

func (r *StorageSQLite) GetFoodList(ctx context.Context, userID int64) ([]s.Food, error) {
	rows, err := r.db.QueryContext(ctx, _sqlGetFoodList, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []s.Food{}
	for rows.Next() {
		var f s.Food
		err = rows.Scan(&f.Key, &f.Name, &f.Brand, &f.Cal100, &f.Prot100, &f.Fat100, &f.Carb100, &f.Comment)
		if err != nil {
			return nil, err
		}

		list = append(list, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, s.ErrEmptyResult
	}

	return list, nil
}

func (r *StorageSQLite) FindFood(ctx context.Context, userID int64, pattern string) ([]s.Food, error) {
	rows, err := r.db.QueryContext(ctx, _sqlFindFood, userID, strings.ToUpper(pattern))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []s.Food{}
	for rows.Next() {
		var f s.Food
		err = rows.Scan(&f.Key, &f.Name, &f.Brand, &f.Cal100, &f.Prot100, &f.Fat100, &f.Carb100, &f.Comment)
		if err != nil {
			return nil, err
		}

		list = append(list, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, s.ErrEmptyResult
	}

	return list, nil
}

func (r *StorageSQLite) SetFood(ctx context.Context, userID int64, food *s.Food) error {
	if !food.Validate() {
		return s.ErrFoodInvalid
	}

	_, err := r.db.ExecContext(ctx,
		_sqlSetFood,
		userID,
		food.Key,
		food.Name,
		food.Brand,
		food.Cal100,
		food.Prot100,
		food.Fat100,
		food.Carb100,
		food.Comment,
	)
	return err
}

func (r *StorageSQLite) DeleteFood(ctx context.Context, userID int64, key string) error {
	bndlList, err := r.GetBundleList(ctx, userID)
	if err != nil && !errors.Is(err, s.ErrEmptyResult) {
		return err
	}

	for _, bndl := range bndlList {
		for k, v := range bndl.Data {
			if v != 0 && k == key {
				return s.ErrFoodIsUsed
			}
		}
	}

	_, err = r.db.ExecContext(ctx, _sqlDeleteFood, userID, key)
	if err != nil {
		var errSql gsql.Error
		if errors.As(err, &errSql) && errSql.Error() == _errForeignKey {
			return s.ErrFoodIsUsed
		}
		return err
	}

	return nil
}
