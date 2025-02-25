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
		r.Equal(int64(6), migrationID)
	})
}

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

	r.Run("get sport for user 1", func() {
		res, err := r.stg.GetSport(context.Background(), 1, "sport1 key")
		r.NoError(err)
		r.Equal(&s.Sport{
			Key:     "sport1 key",
			Name:    "sport1 name",
			Comment: "sport1 comment",
		}, res)
	})

	r.Run("get not exists sport", func() {
		_, err := r.stg.GetSport(context.Background(), 3, "sport1 key")
		r.ErrorIs(err, s.ErrSportNotFound)
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

	r.Run("update sport for user 1", func() {
		r.NoError(r.stg.SetSport(context.Background(), 1, &s.Sport{
			Key:     "sport1 key",
			Name:    "sport1 name new",
			Comment: "sport1 comment new",
		}))
	})

	r.Run("get sport list for user 1", func() {
		res, err := r.stg.GetSportList(context.Background(), 1)
		r.NoError(err)
		r.Equal([]s.Sport{
			{Key: "sport1 key", Name: "sport1 name new", Comment: "sport1 comment new"},
			{Key: "sport2 key", Name: "sport2 name", Comment: "sport2 comment"},
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

func (r *StorageSQLiteTestSuite) TestSportActivityCRUD() {
	r.Run("get empty sport activity report", func() {
		_, err := r.stg.GetSportActivityReport(context.Background(), 1, 1, 3)
		r.ErrorIs(err, s.ErrEmptyResult)
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
			Key:     "sport2 key",
			Name:    "sport2 name",
			Comment: "sport2 comment",
		}))
	})

	r.Run("set invalid sport activity", func() {
		r.ErrorIs(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{}), s.ErrSportActivityInvalid)
		r.ErrorIs(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{SportKey: "sport1"}), s.ErrSportActivityInvalid)
		r.ErrorIs(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{
			SportKey: "sport1",
			Sets:     []int64{0},
		}), s.ErrSportActivityInvalid)
		r.ErrorIs(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{
			SportKey: "sport1",
			Sets:     []int64{-1},
		}), s.ErrSportActivityInvalid)
	})

	r.Run("set sport activity for not found sport", func() {
		r.ErrorIs(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{
			SportKey:  "sport1",
			Timestamp: 1,
			Sets:      []int64{1, 2, 3},
		}), s.ErrSportNotFound)
	})

	r.Run("set sport activity", func() {
		r.NoError(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{
			SportKey:  "sport1 key",
			Timestamp: 1,
			Sets:      []int64{1, 2, 3},
		}))
		r.NoError(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{
			SportKey:  "sport2 key",
			Timestamp: 2,
			Sets:      []int64{4, 5, 6},
		}))
	})

	r.Run("get sport activity report", func() {
		res, err := r.stg.GetSportActivityReport(context.Background(), 1, 1, 3)
		r.NoError(err)
		r.Equal([]s.SportActivityReport{
			{SportName: "sport1 name", Timestamp: 1, Sets: []int64{1, 2, 3}},
			{SportName: "sport2 name", Timestamp: 2, Sets: []int64{4, 5, 6}},
		}, res)
	})

	r.Run("delete sport activity", func() {
		r.NoError(r.stg.DeleteSportActivity(context.Background(), 1, 1, "sport1 key"))
	})

	r.Run("get sport activity report", func() {
		res, err := r.stg.GetSportActivityReport(context.Background(), 1, 1, 3)
		r.NoError(err)
		r.Equal([]s.SportActivityReport{
			{SportName: "sport2 name", Timestamp: 2, Sets: []int64{4, 5, 6}},
		}, res)
	})

	r.Run("update sport activity", func() {
		r.NoError(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{
			SportKey:  "sport2 key",
			Timestamp: 2,
			Sets:      []int64{4, 5, 6, 7, 8, 9},
		}))
	})

	r.Run("get sport activity report", func() {
		res, err := r.stg.GetSportActivityReport(context.Background(), 1, 1, 3)
		r.NoError(err)
		r.Equal([]s.SportActivityReport{
			{SportName: "sport2 name", Timestamp: 2, Sets: []int64{4, 5, 6, 7, 8, 9}},
		}, res)
	})

	r.Run("success delete sport for user 2 and error for user 1", func() {
		r.NoError(r.stg.DeleteSport(context.Background(), 2, "sport2 key"))
		r.ErrorIs(r.stg.DeleteSport(context.Background(), 1, "sport2 key"), s.ErrSportIsUsed)
	})
}

func (r *StorageSQLiteTestSuite) TestUserSettingsCRUD() {
	r.Run("get empty user settings", func() {
		_, err := r.stg.GetUserSettings(context.Background(), 1)
		r.ErrorIs(err, s.ErrUserSettingsNotFound)
	})

	r.Run("set invalid settings", func() {
		r.ErrorIs(r.stg.SetUserSettings(context.Background(), 1, &s.UserSettings{}), s.ErrUserSettingsInvalid)
	})

	r.Run("set user settings", func() {
		r.NoError(r.stg.SetUserSettings(context.Background(), 1, &s.UserSettings{CalLimit: 123.123}))
	})

	r.Run("get user settings", func() {
		res, err := r.stg.GetUserSettings(context.Background(), 1)
		r.NoError(err)
		r.Equal(&s.UserSettings{CalLimit: 123.123}, res)
	})

	r.Run("update user settings", func() {
		r.NoError(r.stg.SetUserSettings(context.Background(), 1, &s.UserSettings{CalLimit: 456.456}))
	})

	r.Run("get updated user settings", func() {
		res, err := r.stg.GetUserSettings(context.Background(), 1)
		r.NoError(err)
		r.Equal(&s.UserSettings{CalLimit: 456.456}, res)
	})
}

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
		SportActivity: []s.SportActivityBackup{
			{UserID: 1, SportKey: "sport1 key", Timestamp: 1, Sets: "[1,2,3]"},
			{UserID: 1, SportKey: "sport2 key", Timestamp: 2, Sets: "[4,5,6]"},
			{UserID: 2, SportKey: "sport1 key", Timestamp: 1, Sets: "[7,8,9]"},
		},
		UserSettings: []s.UserSettingsBackup{
			{UserID: 1, CalLimit: 123.123},
			{UserID: 2, CalLimit: 456.456},
		},
		Food: []s.FoodBackup{
			{
				UserID: 1, Key: "food1_key", Name: "food1_name", Brand: "food1_brand",
				Cal100: 1.1, Prot100: 2.2, Fat100: 3.3, Carb100: 4.4, Comment: "food1_comment",
			},
			{
				UserID: 1, Key: "food2_key", Name: "food2_name", Brand: "food2_brand",
				Cal100: 5.5, Prot100: 6.6, Fat100: 7.7, Carb100: 8.8, Comment: "food2_comment",
			},
			{
				UserID: 2, Key: "food1_key", Name: "food1_name", Brand: "food1_brand",
				Cal100: 1.1, Prot100: 2.2, Fat100: 3.3, Carb100: 4.4, Comment: "food1_comment",
			},
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

		// Sport.
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

		// SportActivity.
		{
			res, err := r.stg.GetSportActivityReport(context.Background(), 1, 1, 3)
			r.NoError(err)
			r.Equal([]s.SportActivityReport{
				{SportName: "sport1 name", Timestamp: 1, Sets: []int64{1, 2, 3}},
				{SportName: "sport2 name", Timestamp: 2, Sets: []int64{4, 5, 6}},
			}, res)

			res, err = r.stg.GetSportActivityReport(context.Background(), 2, 1, 3)
			r.NoError(err)
			r.Equal([]s.SportActivityReport{
				{SportName: "sport1 name", Timestamp: 1, Sets: []int64{7, 8, 9}},
			}, res)
		}

		// UserSettings.
		{
			res, err := r.stg.GetUserSettings(context.Background(), 1)
			r.NoError(err)
			r.Equal(&s.UserSettings{CalLimit: 123.123}, res)

			res, err = r.stg.GetUserSettings(context.Background(), 2)
			r.NoError(err)
			r.Equal(&s.UserSettings{CalLimit: 456.456}, res)
		}

		// Food
		{
			res, err := r.stg.GetFoodList(context.Background(), 1)
			r.NoError(err)
			r.Equal([]s.Food{
				{
					Key:     "food1_key",
					Name:    "food1_name",
					Brand:   "food1_brand",
					Cal100:  1.1,
					Prot100: 2.2,
					Fat100:  3.3,
					Carb100: 4.4,
					Comment: "food1_comment",
				},
				{
					Key:     "food2_key",
					Name:    "food2_name",
					Brand:   "food2_brand",
					Cal100:  5.5,
					Prot100: 6.6,
					Fat100:  7.7,
					Carb100: 8.8,
					Comment: "food2_comment",
				},
			}, res)

			res, err = r.stg.GetFoodList(context.Background(), 2)
			r.NoError(err)
			r.Equal([]s.Food{
				{
					Key:     "food1_key",
					Name:    "food1_name",
					Brand:   "food1_brand",
					Cal100:  1.1,
					Prot100: 2.2,
					Fat100:  3.3,
					Carb100: 4.4,
					Comment: "food1_comment",
				},
			}, res)
		}
	})

	r.Run("do backup and check with initial", func() {
		backup2, err := r.stg.Backup(context.Background())
		r.NoError(err)

		r.Equal(backup.Weight, backup2.Weight)
		r.Equal(backup.Sport, backup2.Sport)
		r.Equal(backup.SportActivity, backup2.SportActivity)
		r.Equal(backup.UserSettings, backup2.UserSettings)
		r.Equal(backup.Food, backup2.Food)
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
