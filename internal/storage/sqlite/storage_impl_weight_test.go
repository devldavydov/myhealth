package sqlite

import (
	"context"

	s "github.com/devldavydov/myhealth/internal/storage"
)

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
