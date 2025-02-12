package sqlite

import (
	"context"
	"os"
	"testing"

	s "github.com/devldavydov/myhealth/internal/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StorageSQLiteTestSuite struct {
	suite.Suite

	stg    *StorageSQLite
	dbFile string
}

func (r *StorageSQLiteTestSuite) TestMigrations() {
	r.Run("check last migration", func() {
		migrationID, err := r.stg.getLastMigrationID(context.Background())
		r.NoError(err)
		r.Equal(int64(3), migrationID)
	})
}

//
// Weight.
//

func (r *StorageSQLiteTestSuite) TestWeightCRUD() {
	r.Run("check empty weight list", func() {
		_, err := r.stg.GetWeightList(context.Background(), 1, 1000, 2000)
		r.ErrorIs(err, s.ErrEmptyResult)
	})

	r.Run("set invalid weight", func() {
		r.ErrorIs(r.stg.SetWeight(context.Background(), 1, &s.Weight{Value: -1}), s.ErrWeightInvalid)
	})

	r.Run("set initial data", func() {
		r.NoError(r.stg.SetWeight(context.Background(), 1, &s.Weight{Timestamp: 1000, Value: 94.3}))
		r.NoError(r.stg.SetWeight(context.Background(), 1, &s.Weight{Timestamp: 2000, Value: 94}))
		r.NoError(r.stg.SetWeight(context.Background(), 1, &s.Weight{Timestamp: 3000, Value: 96}))
		r.NoError(r.stg.SetWeight(context.Background(), 2, &s.Weight{Timestamp: 1000, Value: 87}))
	})

	r.Run("get weight for user 1", func() {
		res, err := r.stg.GetWeightList(context.Background(), 1, 1000, 4000)
		r.NoError(err)
		r.Equal([]s.Weight{
			{Timestamp: 1000, Value: 94.3},
			{Timestamp: 2000, Value: 94},
			{Timestamp: 3000, Value: 96},
		}, res)
	})

	r.Run("get weight for user 2", func() {
		res, err := r.stg.GetWeightList(context.Background(), 2, 1000, 4000)
		r.NoError(err)
		r.Equal([]s.Weight{
			{Timestamp: 1000, Value: 87},
		}, res)
	})

	r.Run("update weight for user1", func() {
		r.NoError(r.stg.SetWeight(context.Background(), 1, &s.Weight{Timestamp: 1000, Value: 104.3}))
	})

	r.Run("check update", func() {
		res, err := r.stg.GetWeightList(context.Background(), 1, 1000, 1000)
		r.NoError(err)
		r.Equal([]s.Weight{
			{Timestamp: 1000, Value: 104.3},
		}, res)
	})

	r.Run("delete weight for user 2", func() {
		r.NoError(r.stg.DeleteWeight(context.Background(), 2, 1000))
	})

	r.Run("get weight empty list for user 2", func() {
		_, err := r.stg.GetWeightList(context.Background(), 2, 1000, 4000)
		r.ErrorIs(err, s.ErrEmptyResult)
	})
}

//
// Sport.
//

func (r *StorageSQLiteTestSuite) TestSportCRUD() {
	r.Run("check empty sport list for user 1", func() {
		_, err := r.stg.GetSportList(context.Background(), 1)
		r.ErrorIs(err, s.ErrEmptyResult)
	})

	r.Run("set invalid sport", func() {
		r.ErrorIs(r.stg.SetSport(context.Background(), 1, &s.Sport{}), s.ErrSportInvalid)
		r.ErrorIs(r.stg.SetSport(context.Background(), 1, &s.Sport{Key: "key"}), s.ErrSportInvalid)
	})

	r.Run("set sport", func() {
		r.NoError(r.stg.SetSport(context.Background(), 1, &s.Sport{
			Key:     "sport1 key",
			Name:    "sport1 name",
			Comment: "sport1 comment",
		}))
		r.NoError(r.stg.SetSport(context.Background(), 1, &s.Sport{
			Key:     "sport2 key",
			Name:    "sport2 name",
			Comment: "sport2 comment",
		}))
		r.NoError(r.stg.SetSport(context.Background(), 2, &s.Sport{
			Key:     "sport1 key",
			Name:    "sport1 name",
			Comment: "sport1 comment",
		}))
	})

	r.Run("get sport list for user 1", func() {
		res, err := r.stg.GetSportList(context.Background(), 1)
		r.NoError(err)
		r.Equal([]s.Sport{
			{Key: "sport1 key", Name: "sport1 name", Comment: "sport1 comment"},
			{Key: "sport2 key", Name: "sport2 name", Comment: "sport2 comment"},
		}, res)
	})

	r.Run("get sport list for user 2", func() {
		res, err := r.stg.GetSportList(context.Background(), 2)
		r.NoError(err)
		r.Equal([]s.Sport{
			{Key: "sport1 key", Name: "sport1 name", Comment: "sport1 comment"},
		}, res)
	})

	r.Run("delete sport for user 2", func() {
		r.NoError(r.stg.DeleteSport(context.Background(), 2, "sport1 key"))
	})

	r.Run("check empty sport list for user 2", func() {
		_, err := r.stg.GetSportList(context.Background(), 2)
		r.ErrorIs(err, s.ErrEmptyResult)
	})
}

//
// Backup/restore.
//

func (r *StorageSQLiteTestSuite) TestBackupRestore() {
	backup := &s.Backup{
		Timestamp: 1000,
		Weight: []s.WeightBackup{
			{UserID: 1, Timestamp: 1000, Value: 90.1},
			{UserID: 1, Timestamp: 2000, Value: 92.1},
			{UserID: 2, Timestamp: 1000, Value: 87.8},
		},
		Sport: []s.SportBackup{
			{UserID: 1, Key: "sport1 key", Name: "sport1 name", Comment: "sport1 comment"},
			{UserID: 1, Key: "sport2 key", Name: "sport2 name", Comment: "sport2 comment"},
			{UserID: 2, Key: "sport1 key", Name: "sport1 name", Comment: "sport1 comment"},
		},
	}

	r.Run("restore backup", func() {
		r.NoError(r.stg.Restore(context.Background(), backup))
	})

	r.Run("check db after restore", func() {
		// Weight.
		{
			res, err := r.stg.GetWeightList(context.Background(), 1, 1000, 3000)
			r.NoError(err)
			r.Equal([]s.Weight{
				{Timestamp: 1000, Value: 90.1},
				{Timestamp: 2000, Value: 92.1},
			}, res)

			res, err = r.stg.GetWeightList(context.Background(), 2, 1000, 3000)
			r.NoError(err)
			r.Equal([]s.Weight{
				{Timestamp: 1000, Value: 87.8},
			}, res)
		}

		// Sport2.
		{
			res, err := r.stg.GetSportList(context.Background(), 1)
			r.NoError(err)
			r.Equal([]s.Sport{
				{Key: "sport1 key", Name: "sport1 name", Comment: "sport1 comment"},
				{Key: "sport2 key", Name: "sport2 name", Comment: "sport2 comment"},
			}, res)

			res, err = r.stg.GetSportList(context.Background(), 2)
			r.NoError(err)
			r.Equal([]s.Sport{
				{Key: "sport1 key", Name: "sport1 name", Comment: "sport1 comment"},
			}, res)
		}

	})

	r.Run("do backup and check with initial", func() {
		backup2, err := r.stg.Backup(context.Background())
		r.NoError(err)

		r.Equal(backup.Weight, backup2.Weight)
		r.Equal(backup.Sport, backup2.Sport)
	})
}

//
// Suite setup.
//

func (r *StorageSQLiteTestSuite) SetupTest() {
	var err error

	f, err := os.CreateTemp("", "testdb")
	require.NoError(r.T(), err)
	r.dbFile = f.Name()
	f.Close()

	r.stg, err = NewStorageSQLite(r.dbFile, nil)
	require.NoError(r.T(), err)
}

func (r *StorageSQLiteTestSuite) TearDownTest() {
	r.stg.Close()
	require.NoError(r.T(), os.Remove(r.dbFile))
}

func TestStorageSQLite(t *testing.T) {
	suite.Run(t, new(StorageSQLiteTestSuite))
}
