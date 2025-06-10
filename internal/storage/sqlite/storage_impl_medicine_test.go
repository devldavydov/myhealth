package sqlite

import (
	"context"

	s "github.com/devldavydov/myhealth/internal/storage"
)

func (r *StorageSQLiteTestSuite) TestMedicineCRUD() {
	r.Run("check empty medicine list for user 1", func() {
		_, err := r.stg.GetMedicineList(context.Background(), 1)
		r.ErrorIs(err, s.ErrEmptyResult)
	})

	r.Run("set invalid medicine", func() {
		r.ErrorIs(r.stg.SetMedicine(context.Background(), 1, &s.Medicine{}), s.ErrMedicineInvalid)
		r.ErrorIs(r.stg.SetMedicine(context.Background(), 1, &s.Medicine{Key: "key"}), s.ErrMedicineInvalid)
		r.ErrorIs(r.stg.SetMedicine(context.Background(), 1, &s.Medicine{Key: "key", Name: "name"}), s.ErrMedicineInvalid)
	})

	r.Run("set medicine", func() {
		r.NoError(r.stg.SetMedicine(context.Background(), 1, &s.Medicine{
			Key:     "med1 key",
			Name:    "med1 name",
			Unit:    "med1 unit",
			Comment: "med1 comment",
		}))
		r.NoError(r.stg.SetMedicine(context.Background(), 1, &s.Medicine{
			Key:     "med2 key",
			Name:    "med2 name",
			Unit:    "med2 unit",
			Comment: "med2 comment",
		}))
		r.NoError(r.stg.SetMedicine(context.Background(), 2, &s.Medicine{
			Key:     "med1 key",
			Name:    "med1 name",
			Unit:    "med1 unit",
			Comment: "med1 comment",
		}))
	})

	r.Run("get medicine for user 1", func() {
		res, err := r.stg.GetMedicine(context.Background(), 1, "med1 key")
		r.NoError(err)
		r.Equal(&s.Medicine{
			Key:     "med1 key",
			Name:    "med1 name",
			Unit:    "med1 unit",
			Comment: "med1 comment",
		}, res)
	})

	r.Run("get not exists medicine", func() {
		_, err := r.stg.GetMedicine(context.Background(), 3, "med1 key")
		r.ErrorIs(err, s.ErrMedicineNotFound)
	})

	r.Run("get medicine list for user 1", func() {
		res, err := r.stg.GetMedicineList(context.Background(), 1)
		r.NoError(err)
		r.Equal([]s.Medicine{
			{Key: "med1 key", Name: "med1 name", Unit: "med1 unit", Comment: "med1 comment"},
			{Key: "med2 key", Name: "med2 name", Unit: "med2 unit", Comment: "med2 comment"},
		}, res)
	})

	r.Run("get medicine list for user 2", func() {
		res, err := r.stg.GetMedicineList(context.Background(), 2)
		r.NoError(err)
		r.Equal([]s.Medicine{
			{Key: "med1 key", Name: "med1 name", Unit: "med1 unit", Comment: "med1 comment"},
		}, res)
	})

	r.Run("update medicine for user 1", func() {
		r.NoError(r.stg.SetMedicine(context.Background(), 1, &s.Medicine{
			Key:     "med1 key",
			Name:    "med1 name new",
			Unit:    "med1 unit new",
			Comment: "med1 comment new",
		}))
	})

	r.Run("get medicine list for user 1", func() {
		res, err := r.stg.GetMedicineList(context.Background(), 1)
		r.NoError(err)
		r.Equal([]s.Medicine{
			{Key: "med1 key", Name: "med1 name new", Unit: "med1 unit new", Comment: "med1 comment new"},
			{Key: "med2 key", Name: "med2 name", Unit: "med2 unit", Comment: "med2 comment"},
		}, res)
	})

	r.Run("delete medicine for user 2", func() {
		r.NoError(r.stg.DeleteMedicine(context.Background(), 2, "med1 key"))
	})

	r.Run("check empty medicine list for user 2", func() {
		_, err := r.stg.GetMedicineList(context.Background(), 2)
		r.ErrorIs(err, s.ErrEmptyResult)
	})
}

func (r *StorageSQLiteTestSuite) TestMedicineIndicatorCRUD() {
	r.Run("get empty medicine indicator report", func() {
		_, err := r.stg.GetMedicineIndicatorReport(context.Background(), 1, 1, 3)
		r.ErrorIs(err, s.ErrEmptyResult)
	})

	r.Run("set medicine", func() {
		r.NoError(r.stg.SetMedicine(context.Background(), 1, &s.Medicine{
			Key:     "med1 key",
			Name:    "med1 name",
			Unit:    "med1 unit",
			Comment: "med1 comment",
		}))
		r.NoError(r.stg.SetMedicine(context.Background(), 1, &s.Medicine{
			Key:     "med2 key",
			Name:    "med2 name",
			Unit:    "med2 unit",
			Comment: "med2 comment",
		}))
		r.NoError(r.stg.SetMedicine(context.Background(), 2, &s.Medicine{
			Key:     "med2 key",
			Name:    "med2 name",
			Unit:    "med2 unit",
			Comment: "med2 comment",
		}))
	})

	r.Run("set invalid medicine indicator", func() {
		r.ErrorIs(r.stg.SetMedicineIndicator(context.Background(), 1, &s.MedicineIndicator{}), s.ErrMedicineIndicatorInvalid)
		r.ErrorIs(r.stg.SetMedicineIndicator(context.Background(), 1, &s.MedicineIndicator{
			MedicineKey: "med1",
			Value:       -1,
		}), s.ErrMedicineIndicatorInvalid)
	})

	r.Run("set medicine indicator for not found medicine", func() {
		r.ErrorIs(r.stg.SetMedicineIndicator(context.Background(), 1, &s.MedicineIndicator{
			MedicineKey: "med1",
			Timestamp:   1,
			Value:       1.1,
		}), s.ErrMedicineNotFound)
	})

	r.Run("set medicine indicator", func() {
		r.NoError(r.stg.SetMedicineIndicator(context.Background(), 1, &s.MedicineIndicator{
			MedicineKey: "med1 key",
			Timestamp:   1,
			Value:       1.1,
		}))
		r.NoError(r.stg.SetMedicineIndicator(context.Background(), 1, &s.MedicineIndicator{
			MedicineKey: "med2 key",
			Timestamp:   2,
			Value:       2.2,
		}))
	})

	r.Run("get medicine indicator report", func() {
		res, err := r.stg.GetMedicineIndicatorReport(context.Background(), 1, 1, 3)
		r.NoError(err)
		r.Equal([]s.MedicineIndicatorReport{
			{MedicineName: "med1 name [med1 unit]", Timestamp: 1, Value: 1.1},
			{MedicineName: "med2 name [med2 unit]", Timestamp: 2, Value: 2.2},
		}, res)
	})

	r.Run("delete medicine indicator", func() {
		r.NoError(r.stg.DeleteMedicineIndicator(context.Background(), 1, 1, "med1 key"))
	})

	r.Run("get medicine indicator report", func() {
		res, err := r.stg.GetMedicineIndicatorReport(context.Background(), 1, 1, 3)
		r.NoError(err)
		r.Equal([]s.MedicineIndicatorReport{
			{MedicineName: "med2 name [med2 unit]", Timestamp: 2, Value: 2.2},
		}, res)
	})

	r.Run("update medicine indicator", func() {
		r.NoError(r.stg.SetMedicineIndicator(context.Background(), 1, &s.MedicineIndicator{
			MedicineKey: "med2 key",
			Timestamp:   2,
			Value:       3.3,
		}))
	})

	r.Run("get medicine indicator report", func() {
		res, err := r.stg.GetMedicineIndicatorReport(context.Background(), 1, 1, 3)
		r.NoError(err)
		r.Equal([]s.MedicineIndicatorReport{
			{MedicineName: "med2 name [med2 unit]", Timestamp: 2, Value: 3.3},
		}, res)
	})

	r.Run("success delete medicine for user 2 and error for user 1", func() {
		r.NoError(r.stg.DeleteMedicine(context.Background(), 2, "med2 key"))
		r.ErrorIs(r.stg.DeleteMedicine(context.Background(), 1, "med2 key"), s.ErrMedicineIsUsed)
	})
}
