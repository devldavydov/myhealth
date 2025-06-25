package sqlite

import (
	"context"

	s "github.com/devldavydov/myhealth/internal/storage"
)

func (r *StorageSQLiteTestSuite) TestTotalBurnedCalCRUD() {
	r.Run("set invalid day total cal", func() {
		r.ErrorIs(r.stg.SetTotalBurnedCal(context.Background(), 1, 1, 0), s.ErrDayTotalCalInvalid)
		r.ErrorIs(r.stg.SetTotalBurnedCal(context.Background(), 1, 1, -1), s.ErrDayTotalCalInvalid)
	})

	r.Run("get not exists", func() {
		_, err := r.stg.GetTotalBurnedCal(context.Background(), 1, 1)
		r.ErrorIs(err, s.ErrTotalBurnedCalNotFound)
	})

	r.Run("set and get day total cal", func() {
		r.NoError(r.stg.SetTotalBurnedCal(context.Background(), 1, 1, 100))

		val, err := r.stg.GetTotalBurnedCal(context.Background(), 1, 1)
		r.NoError(err)
		r.Equal(float64(100), val)
	})

	r.Run("update and get day total cal", func() {
		r.NoError(r.stg.SetTotalBurnedCal(context.Background(), 1, 1, 200))

		val, err := r.stg.GetTotalBurnedCal(context.Background(), 1, 1)
		r.NoError(err)
		r.Equal(float64(200), val)
	})

	r.Run("delete and check not exists", func() {
		r.NoError(r.stg.DeleteTotalBurnedCal(context.Background(), 1, 1))

		_, err := r.stg.GetTotalBurnedCal(context.Background(), 1, 1)
		r.ErrorIs(err, s.ErrTotalBurnedCalNotFound)
	})
}
