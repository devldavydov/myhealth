package sqlite

import (
	"context"

	s "github.com/devldavydov/myhealth/internal/storage"
)

func (r *StorageSQLiteTestSuite) TestSportCRUD() {
	r.Run("check empty sport list for user 1", func() {
		_, err := r.stg.GetSportList(context.Background(), 1)
		r.ErrorIs(err, s.ErrEmptyResult)
	})

	r.Run("set invalid sport", func() {
		r.ErrorIs(r.stg.SetSport(context.Background(), 1, &s.Sport{}), s.ErrSportInvalid)
		r.ErrorIs(r.stg.SetSport(context.Background(), 1, &s.Sport{Key: "key"}), s.ErrSportInvalid)
		r.ErrorIs(r.stg.SetSport(context.Background(), 1, &s.Sport{Key: "key", Name: "name"}), s.ErrSportInvalid)
	})

	r.Run("set sport", func() {
		r.NoError(r.stg.SetSport(context.Background(), 1, &s.Sport{
			Key:     "sport1 key",
			Name:    "sport1 name",
			Unit:    "sport1 unit",
			Comment: "sport1 comment",
		}))
		r.NoError(r.stg.SetSport(context.Background(), 1, &s.Sport{
			Key:     "sport2 key",
			Name:    "sport2 name",
			Unit:    "sport2 unit",
			Comment: "sport2 comment",
		}))
		r.NoError(r.stg.SetSport(context.Background(), 2, &s.Sport{
			Key:     "sport1 key",
			Name:    "sport1 name",
			Unit:    "sport1 unit",
			Comment: "sport1 comment",
		}))
	})

	r.Run("get sport for user 1", func() {
		res, err := r.stg.GetSport(context.Background(), 1, "sport1 key")
		r.NoError(err)
		r.Equal(&s.Sport{
			Key:     "sport1 key",
			Name:    "sport1 name",
			Unit:    "sport1 unit",
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
			{Key: "sport1 key", Name: "sport1 name", Unit: "sport1 unit", Comment: "sport1 comment"},
			{Key: "sport2 key", Name: "sport2 name", Unit: "sport2 unit", Comment: "sport2 comment"},
		}, res)
	})

	r.Run("get sport list for user 2", func() {
		res, err := r.stg.GetSportList(context.Background(), 2)
		r.NoError(err)
		r.Equal([]s.Sport{
			{Key: "sport1 key", Name: "sport1 name", Unit: "sport1 unit", Comment: "sport1 comment"},
		}, res)
	})

	r.Run("update sport for user 1", func() {
		r.NoError(r.stg.SetSport(context.Background(), 1, &s.Sport{
			Key:     "sport1 key",
			Name:    "sport1 name new",
			Unit:    "sport1 unit new",
			Comment: "sport1 comment new",
		}))
	})

	r.Run("get sport list for user 1", func() {
		res, err := r.stg.GetSportList(context.Background(), 1)
		r.NoError(err)
		r.Equal([]s.Sport{
			{Key: "sport1 key", Name: "sport1 name new", Unit: "sport1 unit new", Comment: "sport1 comment new"},
			{Key: "sport2 key", Name: "sport2 name", Unit: "sport2 unit", Comment: "sport2 comment"},
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
			Unit:    "sport1 unit",
			Comment: "sport1 comment",
		}))
		r.NoError(r.stg.SetSport(context.Background(), 1, &s.Sport{
			Key:     "sport2 key",
			Name:    "sport2 name",
			Unit:    "sport2 unit",
			Comment: "sport2 comment",
		}))
		r.NoError(r.stg.SetSport(context.Background(), 2, &s.Sport{
			Key:     "sport2 key",
			Name:    "sport2 name",
			Unit:    "sport2 unit",
			Comment: "sport2 comment",
		}))
	})

	r.Run("set invalid sport activity", func() {
		r.ErrorIs(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{}), s.ErrSportActivityInvalid)
		r.ErrorIs(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{SportKey: "sport1"}), s.ErrSportActivityInvalid)
		r.ErrorIs(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{
			SportKey: "sport1",
			Sets:     []float64{0},
		}), s.ErrSportActivityInvalid)
		r.ErrorIs(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{
			SportKey: "sport1",
			Sets:     []float64{-1},
		}), s.ErrSportActivityInvalid)
	})

	r.Run("set sport activity for not found sport", func() {
		r.ErrorIs(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{
			SportKey:  "sport1",
			Timestamp: 1,
			Sets:      []float64{1, 2, 3},
		}), s.ErrSportNotFound)
	})

	r.Run("set sport activity", func() {
		r.NoError(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{
			SportKey:  "sport1 key",
			Timestamp: 1,
			Sets:      []float64{1, 2, 3},
		}))
		r.NoError(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{
			SportKey:  "sport2 key",
			Timestamp: 2,
			Sets:      []float64{4, 5, 6},
		}))
	})

	r.Run("get sport activity report", func() {
		res, err := r.stg.GetSportActivityReport(context.Background(), 1, 1, 3)
		r.NoError(err)
		r.Equal([]s.SportActivityReport{
			{SportName: "sport1 name", Timestamp: 1, Sets: []float64{1, 2, 3}},
			{SportName: "sport2 name", Timestamp: 2, Sets: []float64{4, 5, 6}},
		}, res)
	})

	r.Run("delete sport activity", func() {
		r.NoError(r.stg.DeleteSportActivity(context.Background(), 1, 1, "sport1 key"))
	})

	r.Run("get sport activity report", func() {
		res, err := r.stg.GetSportActivityReport(context.Background(), 1, 1, 3)
		r.NoError(err)
		r.Equal([]s.SportActivityReport{
			{SportName: "sport2 name", Timestamp: 2, Sets: []float64{4, 5, 6}},
		}, res)
	})

	r.Run("update sport activity", func() {
		r.NoError(r.stg.SetSportActivity(context.Background(), 1, &s.SportActivity{
			SportKey:  "sport2 key",
			Timestamp: 2,
			Sets:      []float64{4, 5, 6, 7, 8, 9},
		}))
	})

	r.Run("get sport activity report", func() {
		res, err := r.stg.GetSportActivityReport(context.Background(), 1, 1, 3)
		r.NoError(err)
		r.Equal([]s.SportActivityReport{
			{SportName: "sport2 name", Timestamp: 2, Sets: []float64{4, 5, 6, 7, 8, 9}},
		}, res)
	})

	r.Run("success delete sport for user 2 and error for user 1", func() {
		r.NoError(r.stg.DeleteSport(context.Background(), 2, "sport2 key"))
		r.ErrorIs(r.stg.DeleteSport(context.Background(), 1, "sport2 key"), s.ErrSportIsUsed)
	})
}
