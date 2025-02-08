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

func (s *StorageSQLite) GetWeightList(ctx context.Context, userID int64, from time.Time, to time.Time) ([]s.Weight, error) {
	panic("unimplemented")
}

func (s *StorageSQLite) SetWeight(ctx context.Context, userID int64, weight *s.Weight) error {
	panic("unimplemented")
}

func (s *StorageSQLite) DeleteWeight(ctx context.Context, userID int64, timestamp time.Time) error {
	panic("unimplemented")
}

//
//
//

func (r *StorageSQLite) Close() {
	if r.db == nil {
		return
	}

	if err := r.db.Close(); err != nil {
		r.logger.Error("db close error", zap.Error(err))
	}
}
