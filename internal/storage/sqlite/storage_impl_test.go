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
		r.Equal(int64(8), migrationID)
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

func (r *StorageSQLiteTestSuite) TestFoodCRUD() {
	r.Run("get food empty list", func() {
		_, err := r.stg.GetFoodList(context.Background(), 1)
		r.ErrorIs(err, s.ErrEmptyResult)
	})

	r.Run("set invalid food", func() {
		for _, f := range []s.Food{
			{},
			{Key: "key"},
			{Key: "key", Name: "name", Cal100: -1},
			{Key: "key", Name: "name", Cal100: 1, Prot100: -1},
			{Key: "key", Name: "name", Cal100: 1, Prot100: 1, Fat100: -1},
			{Key: "key", Name: "name", Cal100: 1, Prot100: 1, Fat100: 1, Carb100: -1},
		} {
			r.ErrorIs(r.stg.SetFood(context.Background(), 1, &f), s.ErrFoodInvalid)
		}
	})

	r.Run("set food", func() {
		r.NoError(r.stg.SetFood(context.Background(), 1, &s.Food{
			Key: "food1_key", Name: "food1_name", Brand: "food1_brand",
			Cal100: 1.1, Prot100: 2.2, Fat100: 3.3, Carb100: 4.4, Comment: "food1_comment",
		}))
		r.NoError(r.stg.SetFood(context.Background(), 1, &s.Food{
			Key: "food2_key", Name: "food2_name", Brand: "food2_brand",
			Cal100: 5.5, Prot100: 6.6, Fat100: 7.7, Carb100: 8.8, Comment: "food2_comment",
		}))
		r.NoError(r.stg.SetFood(context.Background(), 2, &s.Food{
			Key: "food1_key", Name: "food1_name", Brand: "food1_brand",
			Cal100: 1.1, Prot100: 2.2, Fat100: 3.3, Carb100: 4.4, Comment: "food1_comment",
		}))
	})

	r.Run("get food list", func() {
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
	})

	r.Run("get food", func() {
		res, err := r.stg.GetFood(context.Background(), 2, "food1_key")
		r.NoError(err)
		r.Equal(&s.Food{
			Key:     "food1_key",
			Name:    "food1_name",
			Brand:   "food1_brand",
			Cal100:  1.1,
			Prot100: 2.2,
			Fat100:  3.3,
			Carb100: 4.4,
			Comment: "food1_comment",
		}, res)

		_, err = r.stg.GetFood(context.Background(), 2, "food2_key")
		r.ErrorIs(err, s.ErrFoodNotFound)
	})

	r.Run("find food", func() {
		res, err := r.stg.FindFood(context.Background(), 1, "food2")
		r.NoError(err)
		r.Equal([]s.Food{
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

		_, err = r.stg.FindFood(context.Background(), 2, "food2")
		r.ErrorIs(err, s.ErrEmptyResult)
	})

	r.Run("update food", func() {
		r.NoError(r.stg.SetFood(context.Background(), 1, &s.Food{
			Key: "food1_key", Name: "food1_NAME", Brand: "food1_BRAND",
			Cal100: 10.10, Prot100: 11.11, Fat100: 12.12, Carb100: 13.13, Comment: "food1_COMMENT",
		}))
	})

	r.Run("get updated food", func() {
		res, err := r.stg.GetFood(context.Background(), 1, "food1_key")
		r.NoError(err)
		r.Equal(&s.Food{
			Key:     "food1_key",
			Name:    "food1_NAME",
			Brand:   "food1_BRAND",
			Cal100:  10.10,
			Prot100: 11.11,
			Fat100:  12.12,
			Carb100: 13.13,
			Comment: "food1_COMMENT",
		}, res)
	})

	r.Run("delete food", func() {
		r.NoError(r.stg.DeleteFood(context.Background(), 2, "food1_key"))
	})

	r.Run("check empty list after delete", func() {
		_, err := r.stg.GetFoodList(context.Background(), 2)
		r.ErrorIs(err, s.ErrEmptyResult)
	})
}

func (r *StorageSQLiteTestSuite) TestBundkeCRUD() {
	r.Run("get bundle empty list", func() {
		_, err := r.stg.GetBundleList(context.Background(), 1)
		r.ErrorIs(err, s.ErrEmptyResult)
	})

	r.Run("set invalid bundle", func() {
		for _, b := range []s.Bundle{
			{},
			{Key: "key"},
			{Key: "key", Data: map[string]float64{
				"food1": -1,
			}},
		} {
			r.ErrorIs(r.stg.SetBundle(context.Background(), 1, &b, true), s.ErrBundleInvalid)
		}
	})

	r.Run("set bundle recursive", func() {
		r.ErrorIs(r.stg.SetBundle(context.Background(), 1, &s.Bundle{
			Key: "bundle1",
			Data: map[string]float64{
				"bundle1": 0,
			},
		}, true), s.ErrBundleDepRecursive)
	})

	r.Run("set bundle dependent food not found", func() {
		r.ErrorIs(r.stg.SetBundle(context.Background(), 1, &s.Bundle{
			Key: "bundle1",
			Data: map[string]float64{
				"food1_key": 123,
			},
		}, true), s.ErrBundleDepFoodNotFound)
	})

	r.Run("set dependent food", func() {
		r.NoError(r.stg.SetFood(context.Background(), 1, &s.Food{
			Key: "food1_key", Name: "food1_name", Brand: "food1_brand",
			Cal100: 1.1, Prot100: 2.2, Fat100: 3.3, Carb100: 4.4, Comment: "food1_comment",
		}))
		r.NoError(r.stg.SetFood(context.Background(), 1, &s.Food{
			Key: "food2_key", Name: "food2_name", Brand: "food2_brand",
			Cal100: 5.5, Prot100: 6.6, Fat100: 7.7, Carb100: 8.8, Comment: "food2_comment",
		}))
		r.NoError(r.stg.SetFood(context.Background(), 1, &s.Food{
			Key: "food3_key", Name: "food3_name", Brand: "food3_brand",
			Cal100: 9.9, Prot100: 10.10, Fat100: 11.11, Carb100: 12.12, Comment: "food3_comment",
		}))
	})

	r.Run("set bundle dependent food not found", func() {
		r.ErrorIs(r.stg.SetBundle(context.Background(), 1, &s.Bundle{
			Key: "bundle1",
			Data: map[string]float64{
				"food1_key": 123,
				"bundle2":   0,
			},
		}, true), s.ErrBundleDepBundleNotFound)
	})

	r.Run("set dependent bundle", func() {
		r.NoError(r.stg.SetBundle(context.Background(), 1, &s.Bundle{
			Key: "bundle2",
			Data: map[string]float64{
				"food2_key": 456,
			},
		}, true))
	})

	r.Run("set bundle", func() {
		r.NoError(r.stg.SetBundle(context.Background(), 1, &s.Bundle{
			Key: "bundle1",
			Data: map[string]float64{
				"food1_key": 123,
				"bundle2":   0,
			},
		}, true))
	})

	r.Run("get bundle", func() {
		res, err := r.stg.GetBundle(context.Background(), 1, "bundle1")
		r.NoError(err)
		r.Equal(&s.Bundle{
			Key: "bundle1",
			Data: map[string]float64{
				"food1_key": 123,
				"bundle2":   0,
			},
		}, res)
	})

	r.Run("get bundle list", func() {
		res, err := r.stg.GetBundleList(context.Background(), 1)
		r.NoError(err)
		r.Equal([]s.Bundle{
			{
				Key: "bundle1",
				Data: map[string]float64{
					"food1_key": 123,
					"bundle2":   0,
				},
			},
			{
				Key: "bundle2",
				Data: map[string]float64{
					"food2_key": 456,
				},
			},
		}, res)
	})

	r.Run("update bundle", func() {
		r.NoError(r.stg.SetBundle(context.Background(), 1, &s.Bundle{
			Key: "bundle1",
			Data: map[string]float64{
				"food1_key": 123,
				"food3_key": 789,
				"bundle2":   0,
			},
		}, true))
	})

	r.Run("get bundle", func() {
		res, err := r.stg.GetBundle(context.Background(), 1, "bundle1")
		r.NoError(err)
		r.Equal(&s.Bundle{
			Key: "bundle1",
			Data: map[string]float64{
				"food1_key": 123,
				"food3_key": 789,
				"bundle2":   0,
			},
		}, res)
	})

	r.Run("delete dependent bundle", func() {
		r.ErrorIs(r.stg.DeleteBundle(context.Background(), 1, "bundle2"), s.ErrBundleIsUsed)
	})

	r.Run("delete bundle", func() {
		r.NoError(r.stg.DeleteBundle(context.Background(), 1, "bundle1"))
	})

	r.Run("get bundle list", func() {
		res, err := r.stg.GetBundleList(context.Background(), 1)
		r.NoError(err)
		r.Equal([]s.Bundle{
			{
				Key: "bundle2",
				Data: map[string]float64{
					"food2_key": 456,
				},
			},
		}, res)
	})

	r.Run("delete dependent food", func() {
		r.ErrorIs(r.stg.DeleteFood(context.Background(), 1, "food2_key"), s.ErrFoodIsUsed)
	})
}

func (r *StorageSQLiteTestSuite) TestJournalCRUD() {
	r.Run("set invalid journal", func() {
		r.ErrorIs(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 1, Meal: s.Meal(-1), FoodKey: "food", FoodWeight: 100,
		}), s.ErrJournalInvalid)
		r.ErrorIs(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 1, Meal: s.Meal(1), FoodKey: "", FoodWeight: 100,
		}), s.ErrJournalInvalid)
		r.ErrorIs(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 1, Meal: s.Meal(1), FoodKey: "food", FoodWeight: 0,
		}), s.ErrJournalInvalid)
	})

	r.Run("set journal with invalid food", func() {
		r.ErrorIs(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 1, Meal: s.Meal(0), FoodKey: "food", FoodWeight: 100,
		}), s.ErrFoodNotFound)
	})

	r.Run("add food", func() {
		// user 1
		r.NoError(r.stg.SetFood(context.TODO(), 1, &s.Food{
			Key: "food_a", Name: "aaa", Brand: "brand a", Cal100: 1, Prot100: 2, Fat100: 3, Carb100: 4, Comment: "Comment",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), 1, &s.Food{
			Key: "food_b", Name: "bbb", Brand: "brand b", Cal100: 5, Prot100: 6, Fat100: 7, Carb100: 8, Comment: "",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), 1, &s.Food{
			Key: "food_c", Name: "ccc", Brand: "brand c", Cal100: 1, Prot100: 1, Fat100: 1, Carb100: 1, Comment: "ccc",
		}))
		// user 2
		r.NoError(r.stg.SetFood(context.TODO(), 2, &s.Food{
			Key: "food_a", Name: "aaa", Brand: "brand a", Cal100: 1, Prot100: 2, Fat100: 3, Carb100: 4, Comment: "Comment",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), 2, &s.Food{
			Key: "food_b", Name: "bbb", Brand: "brand b", Cal100: 5, Prot100: 6, Fat100: 7, Carb100: 8, Comment: "",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), 2, &s.Food{
			Key: "food_c", Name: "ccc", Brand: "brand c", Cal100: 1, Prot100: 1, Fat100: 1, Carb100: 1, Comment: "ccc",
		}))
	})

	r.Run("set journal for different timestamps and users", func() {
		// user 1, timestamp 1
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 1, Meal: s.Meal(0), FoodKey: "food_b", FoodWeight: 100,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 1, Meal: s.Meal(1), FoodKey: "food_a", FoodWeight: 200,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 1, Meal: s.Meal(2), FoodKey: "food_c", FoodWeight: 300,
		}))

		// user 1, timestamp 2
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 2, Meal: s.Meal(0), FoodKey: "food_b", FoodWeight: 300,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 2, Meal: s.Meal(1), FoodKey: "food_a", FoodWeight: 200,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 2, Meal: s.Meal(1), FoodKey: "food_c", FoodWeight: 100,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 2, Meal: s.Meal(2), FoodKey: "food_c", FoodWeight: 400,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 2, Meal: s.Meal(2), FoodKey: "food_a", FoodWeight: 500,
		}))

		// user 2, timestamp 3
		r.NoError(r.stg.SetJournal(context.TODO(), 2, &s.Journal{
			Timestamp: 3, Meal: s.Meal(0), FoodKey: "food_b", FoodWeight: 300,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 2, &s.Journal{
			Timestamp: 3, Meal: s.Meal(1), FoodKey: "food_a", FoodWeight: 200,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 2, &s.Journal{
			Timestamp: 3, Meal: s.Meal(1), FoodKey: "food_c", FoodWeight: 100,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 2, &s.Journal{
			Timestamp: 3, Meal: s.Meal(1), FoodKey: "food_b", FoodWeight: 400,
		}))
	})

	r.Run("get empty report", func() {
		_, err := r.stg.GetJournalReport(context.TODO(), 1, 10, 20)
		r.ErrorIs(err, s.ErrEmptyResult)
	})

	r.Run("get journal reports for user 1", func() {
		rep, err := r.stg.GetJournalReport(context.TODO(), 1, 1, 2)
		r.NoError(err)
		r.Equal([]s.JournalReport{
			{Timestamp: 1, Meal: s.Meal(0), FoodKey: "food_b", FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 100, Cal: 5, Prot: 6, Fat: 7, Carb: 8},
			{Timestamp: 1, Meal: s.Meal(1), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 200, Cal: 2, Prot: 4, Fat: 6, Carb: 8},
			{Timestamp: 1, Meal: s.Meal(2), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 300, Cal: 3, Prot: 3, Fat: 3, Carb: 3},
			{Timestamp: 2, Meal: s.Meal(0), FoodKey: "food_b", FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 300, Cal: 15, Prot: 18, Fat: 21, Carb: 24},
			{Timestamp: 2, Meal: s.Meal(1), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 200, Cal: 2, Prot: 4, Fat: 6, Carb: 8},
			{Timestamp: 2, Meal: s.Meal(1), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 100, Cal: 1, Prot: 1, Fat: 1, Carb: 1},
			{Timestamp: 2, Meal: s.Meal(2), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 500, Cal: 5, Prot: 10, Fat: 15, Carb: 20},
			{Timestamp: 2, Meal: s.Meal(2), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 400, Cal: 4, Prot: 4, Fat: 4, Carb: 4},
		}, rep)

		foodAvgW, err := r.stg.GetJournalFoodAvgWeight(context.TODO(), 1, 1, 2, "food_b")
		r.NoError(err)
		r.Equal(float64(200), foodAvgW)
	})

	r.Run("check that user 2 gets his data", func() {
		rep, err := r.stg.GetJournalReport(context.TODO(), 2, 1, 3)
		r.NoError(err)
		r.Equal([]s.JournalReport{
			{Timestamp: 3, Meal: s.Meal(0), FoodKey: "food_b", FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 300, Cal: 15, Prot: 18, Fat: 21, Carb: 24},
			{Timestamp: 3, Meal: s.Meal(1), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 200, Cal: 2, Prot: 4, Fat: 6, Carb: 8},
			{Timestamp: 3, Meal: s.Meal(1), FoodKey: "food_b", FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 400, Cal: 20, Prot: 24, Fat: 28, Carb: 32},
			{Timestamp: 3, Meal: s.Meal(1), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 100, Cal: 1, Prot: 1, Fat: 1, Carb: 1},
		}, rep)
	})

	r.Run("update and delete for user 1", func() {
		r.NoError(r.stg.DeleteJournal(context.TODO(), 1, 1, s.Meal(0), "food_b"))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &s.Journal{Timestamp: 1, Meal: s.Meal(1), FoodKey: "food_a", FoodWeight: 300}))

		rep, err := r.stg.GetJournalReport(context.TODO(), 1, 1, 1)
		r.NoError(err)
		r.Equal([]s.JournalReport{
			{Timestamp: 1, Meal: s.Meal(1), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 300, Cal: 3, Prot: 6, Fat: 9, Carb: 12},
			{Timestamp: 1, Meal: s.Meal(2), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 300, Cal: 3, Prot: 3, Fat: 3, Carb: 3},
		}, rep)
	})

	r.Run("try delete used food", func() {
		r.ErrorIs(r.stg.DeleteFood(context.TODO(), 1, "food_a"), s.ErrFoodIsUsed)
	})

	r.Run("delete meal for day", func() {
		rep, err := r.stg.GetJournalReport(context.TODO(), 2, 3, 3)
		r.NoError(err)
		r.Equal(4, len(rep))

		r.NoError(r.stg.DeleteJournalMeal(context.TODO(), 2, 3, s.Meal(0)))
		r.NoError(r.stg.DeleteJournalMeal(context.TODO(), 2, 3, s.Meal(1)))

		_, err = r.stg.GetJournalReport(context.TODO(), 2, 3, 3)
		r.ErrorIs(err, s.ErrEmptyResult)
	})
}

func (r *StorageSQLiteTestSuite) TestJournalCopy() {
	r.Run("add food", func() {
		r.NoError(r.stg.SetFood(context.TODO(), 1, &s.Food{
			Key: "food_a", Name: "aaa", Brand: "brand a", Cal100: 1, Prot100: 2, Fat100: 3, Carb100: 4, Comment: "Comment",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), 1, &s.Food{
			Key: "food_b", Name: "bbb", Brand: "brand b", Cal100: 5, Prot100: 6, Fat100: 7, Carb100: 8, Comment: "",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), 1, &s.Food{
			Key: "food_c", Name: "ccc", Brand: "brand c", Cal100: 1, Prot100: 1, Fat100: 1, Carb100: 1, Comment: "ccc",
		}))
	})

	r.Run("set initial journal", func() {
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 1, Meal: s.Meal(0), FoodKey: "food_b", FoodWeight: 100,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 1, Meal: s.Meal(0), FoodKey: "food_a", FoodWeight: 200,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 2, Meal: s.Meal(0), FoodKey: "food_c", FoodWeight: 300,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &s.Journal{
			Timestamp: 2, Meal: s.Meal(1), FoodKey: "food_a", FoodWeight: 100,
		}))
	})

	r.Run("copy success", func() {
		cnt, err := r.stg.CopyJournal(context.TODO(), 1, 1, s.Meal(0), 2, s.Meal(0))
		r.NoError(err)
		r.Equal(2, cnt)

		cnt, err = r.stg.CopyJournal(context.TODO(), 1, 1, s.Meal(0), 2, s.Meal(1))
		r.NoError(err)
		r.Equal(2, cnt)

		rep, err := r.stg.GetJournalReport(context.TODO(), 1, 2, 2)
		r.NoError(err)
		r.Equal([]s.JournalReport{
			{Timestamp: 2, Meal: s.Meal(0), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 200, Cal: 2, Prot: 4, Fat: 6, Carb: 8},
			{Timestamp: 2, Meal: s.Meal(0), FoodKey: "food_b", FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 100, Cal: 5, Prot: 6, Fat: 7, Carb: 8},
			{Timestamp: 2, Meal: s.Meal(0), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 300, Cal: 3, Prot: 3, Fat: 3, Carb: 3},
			{Timestamp: 2, Meal: s.Meal(1), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 200, Cal: 2, Prot: 4, Fat: 6, Carb: 8},
			{Timestamp: 2, Meal: s.Meal(1), FoodKey: "food_b", FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 100, Cal: 5, Prot: 6, Fat: 7, Carb: 8},
		}, rep)
	})

	r.Run("copy zero", func() {
		cnt, err := r.stg.CopyJournal(context.TODO(), 1, 10, s.Meal(0), 20, s.Meal(1))
		r.NoError(err)
		r.Equal(0, cnt)
	})
}

func (r *StorageSQLiteTestSuite) TestSetJournalBundle() {
	r.Run("add food", func() {
		r.NoError(r.stg.SetFood(context.TODO(), 1, &s.Food{
			Key: "food_a", Name: "aaa", Brand: "brand a", Cal100: 1, Prot100: 2, Fat100: 3, Carb100: 4, Comment: "Comment",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), 1, &s.Food{
			Key: "food_b", Name: "bbb", Brand: "brand b", Cal100: 5, Prot100: 6, Fat100: 7, Carb100: 8, Comment: "",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), 1, &s.Food{
			Key: "food_c", Name: "ccc", Brand: "brand c", Cal100: 9, Prot100: 10, Fat100: 11, Carb100: 12, Comment: "ccc",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), 1, &s.Food{
			Key: "food_d", Name: "ddd", Brand: "brand d", Cal100: 13, Prot100: 14, Fat100: 15, Carb100: 16, Comment: "ccc",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), 1, &s.Food{
			Key: "food_e", Name: "eee", Brand: "brand e", Cal100: 17, Prot100: 18, Fat100: 19, Carb100: 20, Comment: "ccc",
		}))
	})

	r.Run("create bundles", func() {
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &s.Bundle{
			Key: "bndl1",
			Data: map[string]float64{
				"food_a": 100,
			},
		}, true))
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &s.Bundle{
			Key: "bndl2",
			Data: map[string]float64{
				"food_b": 200,
				"bndl1":  0,
			},
		}, true))
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &s.Bundle{
			Key: "bndl3",
			Data: map[string]float64{
				"food_c": 300,
				"bndl2":  0,
			},
		}, true))
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &s.Bundle{
			Key: "bndl4",
			Data: map[string]float64{
				"food_d": 400,
				"bndl3":  0,
			},
		}, true))
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &s.Bundle{
			Key: "bndl5",
			Data: map[string]float64{
				"food_e": 500,
				"bndl4":  0,
			},
		}, true))
		//
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &s.Bundle{
			Key: "bndlA",
			Data: map[string]float64{
				"food_a": 100,
				"food_b": 200,
			},
		}, true))
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &s.Bundle{
			Key: "bndlB",
			Data: map[string]float64{
				"food_c": 300,
				"food_d": 400,
			},
		}, true))
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &s.Bundle{
			Key: "bndlC",
			Data: map[string]float64{
				"food_e": 500,
				"bndlA":  0,
				"bndlB":  0,
			},
		}, true))
	})

	r.Run("set journal bundle", func() {
		r.NoError(r.stg.SetJournalBundle(context.TODO(), 1, 1, s.Meal(0), "bndl5"))
		r.NoError(r.stg.SetJournalBundle(context.TODO(), 1, 1, s.Meal(1), "bndlC"))
	})

	r.Run("check journal", func() {
		rep, err := r.stg.GetJournalReport(context.TODO(), 1, 1, 2)
		r.NoError(err)
		r.Equal([]s.JournalReport{
			{Timestamp: 1, Meal: s.Meal(0), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 100, Cal: 1, Prot: 2, Fat: 3, Carb: 4},
			{Timestamp: 1, Meal: s.Meal(0), FoodKey: "food_b", FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 200, Cal: 10, Prot: 12, Fat: 14, Carb: 16},
			{Timestamp: 1, Meal: s.Meal(0), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 300, Cal: 27, Prot: 30, Fat: 33, Carb: 36},
			{Timestamp: 1, Meal: s.Meal(0), FoodKey: "food_d", FoodName: "ddd", FoodBrand: "brand d",
				FoodWeight: 400, Cal: 52, Prot: 56, Fat: 60, Carb: 64},
			{Timestamp: 1, Meal: s.Meal(0), FoodKey: "food_e", FoodName: "eee", FoodBrand: "brand e",
				FoodWeight: 500, Cal: 85, Prot: 90, Fat: 95, Carb: 100},
			//
			{Timestamp: 1, Meal: s.Meal(1), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 100, Cal: 1, Prot: 2, Fat: 3, Carb: 4},
			{Timestamp: 1, Meal: s.Meal(1), FoodKey: "food_b", FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 200, Cal: 10, Prot: 12, Fat: 14, Carb: 16},
			{Timestamp: 1, Meal: s.Meal(1), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 300, Cal: 27, Prot: 30, Fat: 33, Carb: 36},
			{Timestamp: 1, Meal: s.Meal(1), FoodKey: "food_d", FoodName: "ddd", FoodBrand: "brand d",
				FoodWeight: 400, Cal: 52, Prot: 56, Fat: 60, Carb: 64},
			{Timestamp: 1, Meal: s.Meal(1), FoodKey: "food_e", FoodName: "eee", FoodBrand: "brand e",
				FoodWeight: 500, Cal: 85, Prot: 90, Fat: 95, Carb: 100},
		}, rep)
	})

	r.Run("delete journal bundle", func() {
		r.NoError(r.stg.DelJournalBundle(context.TODO(), 1, 1, s.Meal(0), "bndl5"))
		r.NoError(r.stg.DelJournalBundle(context.TODO(), 1, 1, s.Meal(1), "bndlC"))

		_, err := r.stg.GetJournalReport(context.TODO(), 1, 1, 2)
		r.ErrorIs(err, s.ErrEmptyResult)
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
			{UserID: 1, SportKey: "sport1 key", Timestamp: 1, Sets: []int64{1, 2, 3}},
			{UserID: 1, SportKey: "sport2 key", Timestamp: 2, Sets: []int64{4, 5, 6}},
			{UserID: 2, SportKey: "sport1 key", Timestamp: 1, Sets: []int64{7, 8, 9}},
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
				{Key: "sport1 key", Name: "sport1 name", Comment: "sport1 comment"},
				{Key: "sport2 key", Name: "sport2 name", Comment: "sport2 comment"},
			}, res)

			res, err = r.stg.GetSportList(context.Background(), 2)
			r.NoError(err)
			r.Equal([]s.Sport{
				{Key: "sport1 key", Name: "sport1 name", Comment: "sport1 comment"},
			}, res)
		}

		// SportActivity
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
		r.Equal(backup.UserSettings, backup2.UserSettings)
		r.Equal(backup.Food, backup2.Food)
		r.Equal(backup.Bundle, backup2.Bundle)
		r.Equal(backup.Journal, backup2.Journal)
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
