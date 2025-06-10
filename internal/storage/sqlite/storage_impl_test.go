package sqlite

import (
	"context"
	"os"
	"testing"

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
		r.Equal(int64(13), migrationID)
	})
}

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
