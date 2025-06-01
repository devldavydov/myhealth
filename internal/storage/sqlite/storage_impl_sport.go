package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	s "github.com/devldavydov/myhealth/internal/storage"
	gsql "github.com/mattn/go-sqlite3"
)

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
