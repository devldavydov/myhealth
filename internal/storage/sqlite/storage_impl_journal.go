package sqlite

import (
	"context"
	"database/sql"
	"errors"

	s "github.com/devldavydov/myhealth/internal/storage"
	gsql "github.com/mattn/go-sqlite3"
)

func (r *StorageSQLite) SetJournal(ctx context.Context, userID int64, journal *s.Journal) error {
	if !journal.Validate() {
		return s.ErrJournalInvalid
	}

	_, err := r.db.ExecContext(ctx,
		_sqlSetJournal,
		userID,
		journal.Timestamp,
		journal.Meal,
		journal.FoodKey,
		journal.FoodWeight,
	)
	if err != nil {
		var errSql gsql.Error
		if errors.As(err, &errSql) && errSql.Error() == _errForeignKey {
			return s.ErrFoodNotFound
		}
		return err
	}

	return nil
}

func (r *StorageSQLite) SetJournalBundle(ctx context.Context, userID int64, timestamp s.Timestamp, meal s.Meal, bndlKey string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	foodItems, err := getBundleFoodItems(ctx, tx, userID, bndlKey)
	if err != nil {
		return err
	}

	for _, item := range foodItems {
		if _, err := tx.ExecContext(ctx,
			_sqlSetJournal,
			userID,
			timestamp,
			meal,
			item.foodKey,
			item.foodWeight,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

type bundleFoodItem struct {
	foodKey    string
	foodWeight float64
}

func getBundleFoodItems(ctx context.Context, tx *sql.Tx, userID int64, bndlKey string) ([]bundleFoodItem, error) {
	foodItems := []bundleFoodItem{}
	bndlList := []string{bndlKey}
	i := 0

	for i < len(bndlList) {
		nextBndlKey := bndlList[i]
		bndl, err := getBundle(ctx, tx, userID, nextBndlKey)
		if err != nil {
			return nil, err
		}

		for k, v := range bndl.Data {
			if v == 0 {
				bndlList = append(bndlList, k)
				continue
			}

			_, err := getFood(ctx, tx, userID, k)
			if err != nil {
				return nil, err
			}

			foodItems = append(foodItems, bundleFoodItem{foodKey: k, foodWeight: v})
		}

		i++
	}

	return foodItems, nil
}

func (r *StorageSQLite) DeleteJournal(ctx context.Context, userID int64, timestamp s.Timestamp, meal s.Meal, foodkey string) error {
	_, err := r.db.ExecContext(ctx, _sqlDeleteJournal, userID, timestamp, meal, foodkey)
	return err
}

func (r *StorageSQLite) DeleteJournalMeal(ctx context.Context, userID int64, timestamp s.Timestamp, meal s.Meal) error {
	_, err := r.db.ExecContext(ctx, _sqlDeleteJournalMeal, userID, timestamp, meal)
	return err
}

func (r *StorageSQLite) DelJournalBundle(ctx context.Context, userID int64, timestamp s.Timestamp, meal s.Meal, bndlKey string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	foodItems, err := getBundleFoodItems(ctx, tx, userID, bndlKey)
	if err != nil {
		return err
	}

	for _, item := range foodItems {
		if _, err := tx.ExecContext(ctx,
			_sqlDeleteJournal,
			userID,
			timestamp,
			meal,
			item.foodKey,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *StorageSQLite) GetJournalReport(ctx context.Context, userID int64, from, to s.Timestamp) ([]s.JournalReport, error) {
	rows, err := r.db.QueryContext(ctx, _sqlGetJournalReport, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []s.JournalReport{}
	for rows.Next() {
		var jr s.JournalReport
		err = rows.Scan(
			&jr.Timestamp,
			&jr.Meal,
			&jr.FoodKey,
			&jr.FoodName,
			&jr.FoodBrand,
			&jr.FoodWeight,
			&jr.Cal,
			&jr.Prot,
			&jr.Fat,
			&jr.Carb,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, jr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, s.ErrEmptyResult
	}

	return list, nil
}

func (r *StorageSQLite) CopyJournal(ctx context.Context, userID int64, from s.Timestamp, mealFrom s.Meal, to s.Timestamp, mealTo s.Meal) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Get journal data for mealFrom
	type jData struct {
		foodKey    string
		foodWeight string
	}

	rows, err := r.db.QueryContext(ctx, _sqlGetJournalListForCopy, userID, from, mealFrom)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	list := []jData{}
	for rows.Next() {
		var jd jData
		err = rows.Scan(
			&jd.foodKey,
			&jd.foodWeight,
		)
		if err != nil {
			return 0, err
		}

		list = append(list, jd)
	}

	if err = rows.Err(); err != nil {
		return 0, err
	}

	if len(list) == 0 {
		return 0, nil
	}

	// Save to new meal
	for _, item := range list {
		if _, err := tx.ExecContext(ctx,
			_sqlSetJournal,
			userID,
			to,
			mealTo,
			item.foodKey,
			item.foodWeight,
		); err != nil {
			return 0, err
		}
	}

	return len(list), tx.Commit()
}

func (r *StorageSQLite) GetJournalFoodStat(ctx context.Context, userID int64, foodkey string) (*s.JournalFoodStat, error) {
	var fs s.JournalFoodStat
	err := r.db.
		QueryRowContext(ctx, _sqlJournalFoodStat, userID, foodkey).
		Scan(&fs.FirstTimestamp, &fs.LastTimestamp, &fs.TotalWeight, &fs.AvgWeight, &fs.TotalCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, s.ErrEmptyResult
		}
		return nil, err
	}

	return &fs, nil
}
