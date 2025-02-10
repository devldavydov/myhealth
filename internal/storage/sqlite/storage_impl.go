package sqlite

import (
	"context"
	"database/sql"
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
// Weight
//

func (r *StorageSQLite) GetWeightList(ctx context.Context, userID int64, from time.Time, to time.Time) ([]s.Weight, error) {
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

	if len(list) == 0 {
		return nil, s.ErrEmptyResult
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *StorageSQLite) SetWeight(ctx context.Context, userID int64, weight *s.Weight) error {
	if !weight.Validate() {
		return s.ErrWeightInvalid
	}

	_, err := r.db.ExecContext(ctx, _sqlSetWeight, userID, weight.Timestamp, weight.Value)
	if err != nil {
		return err
	}

	return nil
}

func (r *StorageSQLite) DeleteWeight(ctx context.Context, userID int64, timestamp time.Time) error {
	_, err := r.db.ExecContext(ctx, _sqlDeleteWeight, userID, timestamp)
	return err
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
