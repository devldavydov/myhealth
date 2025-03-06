package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	s "github.com/devldavydov/myhealth/internal/storage"
	gsql "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

const (
	_databaseInitTimeout = 30 * time.Second

	_customDriverName = "sqlite3_custom"
	_errForeignKey    = "FOREIGN KEY constraint failed"
)

type StorageSQLite struct {
	db     *sql.DB
	logger *zap.Logger
}

var _ s.Storage = (*StorageSQLite)(nil)

func go_upper(str string) string {
	return strings.ToUpper(str)
}

func NewStorageSQLite(dbFilePath string, logger *zap.Logger) (*StorageSQLite, error) {
	//
	// Driver register (check registration twice).
	//

	if !isDriverRegistered(_customDriverName) {
		sql.Register(_customDriverName, &gsql.SQLiteDriver{
			ConnectHook: func(conn *gsql.SQLiteConn) error {
				if err := conn.RegisterFunc("go_upper", go_upper, false); err != nil {
					return err
				}
				return nil
			},
		})
	}

	//
	// Open DB.
	//

	db, err := sql.Open(
		_customDriverName,
		fmt.Sprintf("file:%s?mode=rwc&_timeout=5000&_fk=1&_sync=1&_journal=wal", dbFilePath),
	)
	if err != nil {
		return nil, err
	}

	stg := &StorageSQLite{db: db, logger: logger}
	if err := stg.doMigrations(); err != nil {
		return nil, err
	}

	return stg, nil
}

//
// Weight.
//

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

//
// Food.
//

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

//
// Bundle.
//

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

//
// Journal.
//

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

func (r *StorageSQLite) GetJournalFoodAvgWeight(ctx context.Context, userID int64, from, to s.Timestamp, foodkey string) (float64, error) {
	var foodAvgWeight float64

	if err := r.db.
		QueryRowContext(ctx, _sqlJournalFoodAvgWeight, userID, foodkey, from, to).
		Scan(&foodAvgWeight); err != nil {
		return 0, err
	}

	return foodAvgWeight, nil
}

//
// Sport.
//

func (r *StorageSQLite) GetSport(ctx context.Context, userID int64, key string) (*s.Sport, error) {
	var sp s.Sport
	err := r.db.
		QueryRowContext(ctx, _sqlGetSport, userID, key).
		Scan(&sp.Key, &sp.Name, &sp.Comment)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, s.ErrSportNotFound
		}
		return nil, err
	}

	return &sp, nil
}

func (r *StorageSQLite) GetSportList(ctx context.Context, userID int64) ([]s.Sport, error) {
	rows, err := r.db.QueryContext(ctx, _sqlGetSportList, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []s.Sport{}
	for rows.Next() {
		var sp s.Sport
		err = rows.Scan(&sp.Key, &sp.Name, &sp.Comment)
		if err != nil {
			return nil, err
		}

		list = append(list, sp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, s.ErrEmptyResult
	}

	return list, nil
}

func (r *StorageSQLite) SetSport(ctx context.Context, userID int64, sp *s.Sport) error {
	if !sp.Validate() {
		return s.ErrSportInvalid
	}

	_, err := r.db.ExecContext(ctx, _sqlSetSport, userID, sp.Key, sp.Name, sp.Comment)
	return err
}

func (r *StorageSQLite) DeleteSport(ctx context.Context, userID int64, key string) error {
	_, err := r.db.ExecContext(ctx, _sqlDeleteSport, userID, key)
	if err != nil {
		var errSql gsql.Error
		if errors.As(err, &errSql) && errSql.Error() == _errForeignKey {
			return s.ErrSportIsUsed
		}
		return err
	}

	return nil
}

//
// SportActivity.
//

func (r *StorageSQLite) SetSportActivity(ctx context.Context, userID int64, sa *s.SportActivity) error {
	if !sa.Validate() {
		return s.ErrSportActivityInvalid
	}

	bSets, err := json.Marshal(sa.Sets)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, _sqlSetSportActivity, userID, sa.Timestamp, sa.SportKey, string(bSets))
	if err != nil {
		var errSql gsql.Error
		if errors.As(err, &errSql) && errSql.Error() == _errForeignKey {
			return s.ErrSportNotFound
		}
		return err
	}

	return nil
}

func (r *StorageSQLite) DeleteSportActivity(ctx context.Context, userID int64, timestamp s.Timestamp, sport_key string) error {
	_, err := r.db.ExecContext(ctx, _sqlDeleteSportActivity, userID, timestamp, sport_key)
	return err
}

func (r *StorageSQLite) GetSportActivityReport(ctx context.Context, userID int64, from, to s.Timestamp) ([]s.SportActivityReport, error) {
	rows, err := r.db.QueryContext(ctx, _sqlGetSportActivityReport, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []s.SportActivityReport{}
	for rows.Next() {
		var sr s.SportActivityReport
		var sSets string

		err = rows.Scan(&sr.Timestamp, &sr.SportName, &sSets)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal([]byte(sSets), &sr.Sets); err != nil {
			return nil, err
		}

		list = append(list, sr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, s.ErrEmptyResult
	}

	return list, nil
}

//
// UserSettings.
//

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

//
// Backup/restore.
//

func (r *StorageSQLite) Backup(ctx context.Context) (*s.Backup, error) {
	backup := &s.Backup{
		Timestamp: s.Timestamp(time.Now().UnixMilli()),
	}

	// Weight
	{
		rows, err := r.db.QueryContext(ctx, _sqlWeightBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.Weight = []s.WeightBackup{}
		for rows.Next() {
			var w s.WeightBackup
			err = rows.Scan(&w.UserID, &w.Timestamp, &w.Value)
			if err != nil {
				return nil, err
			}

			backup.Weight = append(backup.Weight, w)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// Sport
	{
		rows, err := r.db.QueryContext(ctx, _sqlSportBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.Sport = []s.SportBackup{}
		for rows.Next() {
			var sp s.SportBackup
			err = rows.Scan(&sp.UserID, &sp.Key, &sp.Name, &sp.Comment)
			if err != nil {
				return nil, err
			}

			backup.Sport = append(backup.Sport, sp)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// SportActivity
	{
		rows, err := r.db.QueryContext(ctx, _sqlSportActivityBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.SportActivity = []s.SportActivityBackup{}
		for rows.Next() {
			var sa s.SportActivityBackup
			var saSets string
			err = rows.Scan(&sa.UserID, &sa.Timestamp, &sa.SportKey, &saSets)
			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal([]byte(saSets), &sa.Sets); err != nil {
				return nil, err
			}

			backup.SportActivity = append(backup.SportActivity, sa)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// UserSettings
	{
		rows, err := r.db.QueryContext(ctx, _sqlUserSettingsBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.UserSettings = []s.UserSettingsBackup{}
		for rows.Next() {
			var us s.UserSettingsBackup
			err = rows.Scan(&us.UserID, &us.CalLimit)
			if err != nil {
				return nil, err
			}

			backup.UserSettings = append(backup.UserSettings, us)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// Food
	{
		rows, err := r.db.QueryContext(ctx, _sqlFoodBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.Food = []s.FoodBackup{}
		for rows.Next() {
			var f s.FoodBackup
			err = rows.Scan(
				&f.UserID,
				&f.Key,
				&f.Name,
				&f.Brand,
				&f.Cal100,
				&f.Prot100,
				&f.Fat100,
				&f.Carb100,
				&f.Comment,
			)
			if err != nil {
				return nil, err
			}

			backup.Food = append(backup.Food, f)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// Bundle
	{
		rows, err := r.db.QueryContext(ctx, _sqlBundleBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.Bundle = []s.BundleBackup{}
		for rows.Next() {
			var b s.BundleBackup
			var bData string

			err = rows.Scan(
				&b.UserID,
				&b.Key,
				&bData,
			)
			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal([]byte(bData), &b.Data); err != nil {
				return nil, err
			}

			backup.Bundle = append(backup.Bundle, b)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// Journal
	{
		rows, err := r.db.QueryContext(ctx, _sqlJournalBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.Journal = []s.JournalBackup{}
		for rows.Next() {
			var j s.JournalBackup

			err = rows.Scan(
				&j.UserID,
				&j.Timestamp,
				&j.Meal,
				&j.FoodKey,
				&j.FoodWeight,
			)
			if err != nil {
				return nil, err
			}

			backup.Journal = append(backup.Journal, j)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// Result
	return backup, nil
}

func (r *StorageSQLite) Restore(ctx context.Context, backup *s.Backup) error {
	for _, w := range backup.Weight {
		if err := r.SetWeight(
			ctx,
			w.UserID,
			&s.Weight{Timestamp: w.Timestamp, Value: w.Value},
		); err != nil {
			return err
		}
	}

	for _, sp := range backup.Sport {
		if err := r.SetSport(
			ctx,
			sp.UserID,
			&s.Sport{Key: sp.Key, Name: sp.Name, Comment: sp.Comment},
		); err != nil {
			return err
		}
	}

	for _, sa := range backup.SportActivity {
		if err := r.SetSportActivity(ctx, sa.UserID, &s.SportActivity{
			SportKey:  sa.SportKey,
			Timestamp: sa.Timestamp,
			Sets:      sa.Sets,
		}); err != nil {
			return err
		}
	}

	for _, us := range backup.UserSettings {
		if err := r.SetUserSettings(
			ctx,
			us.UserID,
			&s.UserSettings{CalLimit: us.CalLimit},
		); err != nil {
			return err
		}
	}

	for _, f := range backup.Food {
		if err := r.SetFood(
			ctx,
			f.UserID,
			&s.Food{
				Key:     f.Key,
				Name:    f.Name,
				Brand:   f.Brand,
				Cal100:  f.Cal100,
				Prot100: f.Prot100,
				Fat100:  f.Fat100,
				Carb100: f.Carb100,
				Comment: f.Comment,
			},
		); err != nil {
			return err
		}
	}

	for _, b := range backup.Bundle {
		if err := r.SetBundle(ctx, b.UserID, &s.Bundle{
			Key:  b.Key,
			Data: b.Data,
		}, false); err != nil {
			return err
		}
	}

	for _, j := range backup.Journal {
		if err := r.SetJournal(ctx, j.UserID, &s.Journal{
			Timestamp:  j.Timestamp,
			Meal:       j.Meal,
			FoodKey:    j.FoodKey,
			FoodWeight: j.FoodWeight,
		}); err != nil {
			return err
		}
	}

	return nil
}

//
//
//

func (r *StorageSQLite) Close() error {
	if r.db == nil {
		return nil
	}

	return r.db.Close()
}
