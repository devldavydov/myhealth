package sqlite

import (
	"context"

	s "github.com/devldavydov/myhealth/internal/storage"
)

func (r *StorageSQLiteTestSuite) TestBundleCRUD() {
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
