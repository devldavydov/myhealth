package sqlite

import (
	"context"

	s "github.com/devldavydov/myhealth/internal/storage"
)

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

		foodStat, err := r.stg.GetJournalFoodStat(context.TODO(), 1, "food_b")
		r.NoError(err)
		r.Equal(&s.JournalFoodStat{
			FirstTimestamp: 1,
			LastTimestamp:  2,
			TotalWeight:    400,
			AvgWeight:      200,
			TotalCount:     2,
		}, foodStat)
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
