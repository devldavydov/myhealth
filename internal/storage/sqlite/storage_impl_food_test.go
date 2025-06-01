package sqlite

import (
	"context"

	s "github.com/devldavydov/myhealth/internal/storage"
)

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
