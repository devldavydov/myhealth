package sqlite

import (
	"context"

	s "github.com/devldavydov/myhealth/internal/storage"
)

func (r *StorageSQLiteTestSuite) TestBackupRestore() {
	backup := &s.Backup{
		Timestamp: 1000,
		Weight: []s.WeightBackup{
			{UserID: 1, Timestamp: 1000, Value: 90.1},
			{UserID: 1, Timestamp: 2000, Value: 92.1},
			{UserID: 2, Timestamp: 1000, Value: 87.8},
		},
		Sport: []s.SportBackup{
			{UserID: 1, Key: "sport1 key", Name: "sport1 name", Unit: "sport1 unit", Comment: "sport1 comment"},
			{UserID: 1, Key: "sport2 key", Name: "sport2 name", Unit: "sport2 unit", Comment: "sport2 comment"},
			{UserID: 2, Key: "sport1 key", Name: "sport1 name", Unit: "sport1 unit", Comment: "sport1 comment"},
		},
		SportActivity: []s.SportActivityBackup{
			{UserID: 1, SportKey: "sport1 key", Timestamp: 1, Sets: []float64{1, 2, 3}},
			{UserID: 1, SportKey: "sport2 key", Timestamp: 2, Sets: []float64{4, 5, 6}},
			{UserID: 2, SportKey: "sport1 key", Timestamp: 1, Sets: []float64{7, 8, 9}},
		},
		Medicine: []s.MedicineBackup{
			{UserID: 1, Key: "med1 key", Name: "med1 name", Comment: "med1 comment"},
			{UserID: 1, Key: "med2 key", Name: "med2 name", Comment: "med2 comment"},
			{UserID: 2, Key: "med1 key", Name: "med1 name", Comment: "med1 comment"},
		},
		MedicineIndicator: []s.MedicineIndicatorBackup{
			{UserID: 1, MedicineKey: "med1 key", Timestamp: 1, Value: 1.23},
			{UserID: 1, MedicineKey: "med2 key", Timestamp: 2, Value: 4.56},
			{UserID: 2, MedicineKey: "med1 key", Timestamp: 1, Value: 7.89},
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
		Bundle: []s.BundleBackup{
			{UserID: 1, Key: "bundle1", Data: map[string]float64{
				"bundle2":   0,
				"food1_key": 123,
			}},
			{UserID: 1, Key: "bundle2", Data: map[string]float64{"food2_key": 456}},
			{UserID: 2, Key: "bundle1", Data: map[string]float64{"food1_key": 789}},
		},
		Journal: []s.JournalBackup{
			{UserID: 1, Timestamp: 1, Meal: s.Meal(0), FoodKey: "food1_key", FoodWeight: 100},
			{UserID: 1, Timestamp: 1, Meal: s.Meal(1), FoodKey: "food2_key", FoodWeight: 200},
			{UserID: 2, Timestamp: 2, Meal: s.Meal(2), FoodKey: "food1_key", FoodWeight: 100},
		},
	}

	r.Run("restore backup", func() {
		r.NoError(r.stg.Restore(context.Background(), backup))
	})

	r.Run("check db after restore", func() {
		// Weight
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

		// Sport
		{
			res, err := r.stg.GetSportList(context.Background(), 1)
			r.NoError(err)
			r.Equal([]s.Sport{
				{Key: "sport1 key", Name: "sport1 name", Unit: "sport1 unit", Comment: "sport1 comment"},
				{Key: "sport2 key", Name: "sport2 name", Unit: "sport2 unit", Comment: "sport2 comment"},
			}, res)

			res, err = r.stg.GetSportList(context.Background(), 2)
			r.NoError(err)
			r.Equal([]s.Sport{
				{Key: "sport1 key", Name: "sport1 name", Unit: "sport1 unit", Comment: "sport1 comment"},
			}, res)
		}

		// SportActivity
		{
			res, err := r.stg.GetSportActivityReport(context.Background(), 1, 1, 3)
			r.NoError(err)
			r.Equal([]s.SportActivityReport{
				{SportName: "sport1 name", Timestamp: 1, Sets: []float64{1, 2, 3}},
				{SportName: "sport2 name", Timestamp: 2, Sets: []float64{4, 5, 6}},
			}, res)

			res, err = r.stg.GetSportActivityReport(context.Background(), 2, 1, 3)
			r.NoError(err)
			r.Equal([]s.SportActivityReport{
				{SportName: "sport1 name", Timestamp: 1, Sets: []float64{7, 8, 9}},
			}, res)
		}

		// Medicine
		{
			res, err := r.stg.GetMedicineList(context.Background(), 1)
			r.NoError(err)
			r.Equal([]s.Medicine{
				{Key: "med1 key", Name: "med1 name", Comment: "med1 comment"},
				{Key: "med2 key", Name: "med2 name", Comment: "med2 comment"},
			}, res)

			res, err = r.stg.GetMedicineList(context.Background(), 2)
			r.NoError(err)
			r.Equal([]s.Medicine{
				{Key: "med1 key", Name: "med1 name", Comment: "med1 comment"},
			}, res)
		}

		// MedicineIndicator
		{
			res, err := r.stg.GetMedicineIndicatorReport(context.Background(), 1, 1, 3)
			r.NoError(err)
			r.Equal([]s.MedicineIndicatorReport{
				{MedicineName: "med1 name", Timestamp: 1, Value: 1.23},
				{MedicineName: "med2 name", Timestamp: 2, Value: 4.56},
			}, res)

			res, err = r.stg.GetMedicineIndicatorReport(context.Background(), 2, 1, 3)
			r.NoError(err)
			r.Equal([]s.MedicineIndicatorReport{
				{MedicineName: "med1 name", Timestamp: 1, Value: 7.89},
			}, res)
		}

		// UserSettings
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

		// Bundle
		{
			res, err := r.stg.GetBundleList(context.Background(), 1)
			r.NoError(err)
			r.Equal([]s.Bundle{
				{Key: "bundle1", Data: map[string]float64{
					"bundle2":   0,
					"food1_key": 123,
				}},
				{Key: "bundle2", Data: map[string]float64{
					"food2_key": 456,
				}},
			}, res)

			res, err = r.stg.GetBundleList(context.Background(), 2)
			r.NoError(err)
			r.Equal([]s.Bundle{
				{Key: "bundle1", Data: map[string]float64{
					"food1_key": 789,
				}},
			}, res)
		}

		// Journal
		{
			rep, err := r.stg.GetJournalReport(context.Background(), 1, 1, 2)
			r.NoError(err)
			r.Equal([]s.JournalReport{
				{Timestamp: 1, Meal: s.Meal(0), FoodKey: "food1_key", FoodName: "food1_name", FoodBrand: "food1_brand",
					FoodWeight: 100, Cal: 1.1, Prot: 2.2, Fat: 3.3, Carb: 4.4},
				{Timestamp: 1, Meal: s.Meal(1), FoodKey: "food2_key", FoodName: "food2_name", FoodBrand: "food2_brand",
					FoodWeight: 200, Cal: 11, Prot: 13.2, Fat: 15.4, Carb: 17.6},
			}, rep)

			rep, err = r.stg.GetJournalReport(context.Background(), 2, 1, 2)
			r.NoError(err)
			r.Equal([]s.JournalReport{
				{Timestamp: 2, Meal: s.Meal(2), FoodKey: "food1_key", FoodName: "food1_name", FoodBrand: "food1_brand",
					FoodWeight: 100, Cal: 1.1, Prot: 2.2, Fat: 3.3, Carb: 4.4},
			}, rep)
		}
	})

	r.Run("do backup and check with initial", func() {
		backup2, err := r.stg.Backup(context.Background())
		r.NoError(err)

		r.Equal(backup.Weight, backup2.Weight)
		r.Equal(backup.Sport, backup2.Sport)
		r.Equal(backup.SportActivity, backup2.SportActivity)
		r.Equal(backup.Medicine, backup2.Medicine)
		r.Equal(backup.MedicineIndicator, backup2.MedicineIndicator)
		r.Equal(backup.UserSettings, backup2.UserSettings)
		r.Equal(backup.Food, backup2.Food)
		r.Equal(backup.Bundle, backup2.Bundle)
		r.Equal(backup.Journal, backup2.Journal)
	})
}
