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
			err = rows.Scan(&sa.UserID, &sa.Timestamp, &sa.SportKey, &sa.Sets)
			if err != nil {
				return nil, err
			}

			backup.SportActivity = append(backup.SportActivity, sa)
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
		var sets []int64
		if err := json.Unmarshal([]byte(sa.Sets), &sets); err != nil {
			return err
		}

		if err := r.SetSportActivity(ctx, sa.UserID, &s.SportActivity{
			SportKey:  sa.SportKey,
			Timestamp: sa.Timestamp,
			Sets:      sets,
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
