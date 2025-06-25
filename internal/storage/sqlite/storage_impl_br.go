package sqlite

import (
	"context"
	"encoding/json"
	"time"

	s "github.com/devldavydov/myhealth/internal/storage"
)

func (r *StorageSQLite) Backup(ctx context.Context) (*s.Backup, error) {
	backup := &s.Backup{
		Timestamp: s.Timestamp(time.Now().UnixMilli()),
	}

	// Weight
	{
		rows, err := r.db.QueryContext(ctx, _sqlWeightBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.Weight = []s.WeightBackup{}
		for rows.Next() {
			var w s.WeightBackup
			err = rows.Scan(&w.UserID, &w.Timestamp, &w.Value)
			if err != nil {
				return nil, err
			}

			backup.Weight = append(backup.Weight, w)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// Sport
	{
		rows, err := r.db.QueryContext(ctx, _sqlSportBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.Sport = []s.SportBackup{}
		for rows.Next() {
			var sp s.SportBackup
			err = rows.Scan(&sp.UserID, &sp.Key, &sp.Name, &sp.Comment, &sp.Unit)
			if err != nil {
				return nil, err
			}

			backup.Sport = append(backup.Sport, sp)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// SportActivity
	{
		rows, err := r.db.QueryContext(ctx, _sqlSportActivityBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.SportActivity = []s.SportActivityBackup{}
		for rows.Next() {
			var sa s.SportActivityBackup
			var saSets string
			err = rows.Scan(&sa.UserID, &sa.Timestamp, &sa.SportKey, &saSets)
			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal([]byte(saSets), &sa.Sets); err != nil {
				return nil, err
			}

			backup.SportActivity = append(backup.SportActivity, sa)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// Medicine
	{
		rows, err := r.db.QueryContext(ctx, _sqlMedicineBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.Medicine = []s.MedicineBackup{}
		for rows.Next() {
			var m s.MedicineBackup
			err = rows.Scan(&m.UserID, &m.Key, &m.Name, &m.Comment, &m.Unit)
			if err != nil {
				return nil, err
			}

			backup.Medicine = append(backup.Medicine, m)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// MedicineIndicator
	{
		rows, err := r.db.QueryContext(ctx, _sqlMedicineIndicatorBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.MedicineIndicator = []s.MedicineIndicatorBackup{}
		for rows.Next() {
			var mi s.MedicineIndicatorBackup
			err = rows.Scan(&mi.UserID, &mi.Timestamp, &mi.MedicineKey, &mi.Value)
			if err != nil {
				return nil, err
			}

			backup.MedicineIndicator = append(backup.MedicineIndicator, mi)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// UserSettings
	{
		rows, err := r.db.QueryContext(ctx, _sqlUserSettingsBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.UserSettings = []s.UserSettingsBackup{}
		for rows.Next() {
			var us s.UserSettingsBackup
			err = rows.Scan(&us.UserID, &us.CalLimit)
			if err != nil {
				return nil, err
			}

			backup.UserSettings = append(backup.UserSettings, us)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// Food
	{
		rows, err := r.db.QueryContext(ctx, _sqlFoodBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.Food = []s.FoodBackup{}
		for rows.Next() {
			var f s.FoodBackup
			err = rows.Scan(
				&f.UserID,
				&f.Key,
				&f.Name,
				&f.Brand,
				&f.Cal100,
				&f.Prot100,
				&f.Fat100,
				&f.Carb100,
				&f.Comment,
			)
			if err != nil {
				return nil, err
			}

			backup.Food = append(backup.Food, f)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// Bundle
	{
		rows, err := r.db.QueryContext(ctx, _sqlBundleBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.Bundle = []s.BundleBackup{}
		for rows.Next() {
			var b s.BundleBackup
			var bData string

			err = rows.Scan(
				&b.UserID,
				&b.Key,
				&bData,
			)
			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal([]byte(bData), &b.Data); err != nil {
				return nil, err
			}

			backup.Bundle = append(backup.Bundle, b)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// Journal
	{
		rows, err := r.db.QueryContext(ctx, _sqlJournalBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.Journal = []s.JournalBackup{}
		for rows.Next() {
			var j s.JournalBackup

			err = rows.Scan(
				&j.UserID,
				&j.Timestamp,
				&j.Meal,
				&j.FoodKey,
				&j.FoodWeight,
			)
			if err != nil {
				return nil, err
			}

			backup.Journal = append(backup.Journal, j)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// DayTotalCal
	{
		rows, err := r.db.QueryContext(ctx, _sqlDayTotalBackup)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		backup.DayTotalCal = []s.DayTotalCalBackup{}
		for rows.Next() {
			var t s.DayTotalCalBackup

			err = rows.Scan(
				&t.UserID,
				&t.Timestamp,
				&t.TotalCal,
			)
			if err != nil {
				return nil, err
			}

			backup.DayTotalCal = append(backup.DayTotalCal, t)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	// Result
	return backup, nil
}

func (r *StorageSQLite) Restore(ctx context.Context, backup *s.Backup) error {
	for _, w := range backup.Weight {
		if err := r.SetWeight(
			ctx,
			w.UserID,
			&s.Weight{Timestamp: w.Timestamp, Value: w.Value},
		); err != nil {
			return err
		}
	}

	for _, sp := range backup.Sport {
		if err := r.SetSport(
			ctx,
			sp.UserID,
			&s.Sport{Key: sp.Key, Name: sp.Name, Unit: sp.Unit, Comment: sp.Comment},
		); err != nil {
			return err
		}
	}

	for _, sa := range backup.SportActivity {
		if err := r.SetSportActivity(ctx, sa.UserID, &s.SportActivity{
			SportKey:  sa.SportKey,
			Timestamp: sa.Timestamp,
			Sets:      sa.Sets,
		}); err != nil {
			return err
		}
	}

	for _, m := range backup.Medicine {
		if err := r.SetMedicine(
			ctx,
			m.UserID,
			&s.Medicine{Key: m.Key, Name: m.Name, Unit: m.Unit, Comment: m.Comment},
		); err != nil {
			return err
		}
	}

	for _, m := range backup.MedicineIndicator {
		if err := r.SetMedicineIndicator(
			ctx,
			m.UserID,
			&s.MedicineIndicator{
				MedicineKey: m.MedicineKey,
				Timestamp:   m.Timestamp,
				Value:       m.Value,
			},
		); err != nil {
			return err
		}
	}

	for _, us := range backup.UserSettings {
		if err := r.SetUserSettings(
			ctx,
			us.UserID,
			&s.UserSettings{CalLimit: us.CalLimit},
		); err != nil {
			return err
		}
	}

	for _, f := range backup.Food {
		if err := r.SetFood(
			ctx,
			f.UserID,
			&s.Food{
				Key:     f.Key,
				Name:    f.Name,
				Brand:   f.Brand,
				Cal100:  f.Cal100,
				Prot100: f.Prot100,
				Fat100:  f.Fat100,
				Carb100: f.Carb100,
				Comment: f.Comment,
			},
		); err != nil {
			return err
		}
	}

	for _, b := range backup.Bundle {
		if err := r.SetBundle(ctx, b.UserID, &s.Bundle{
			Key:  b.Key,
			Data: b.Data,
		}, false); err != nil {
			return err
		}
	}

	for _, j := range backup.Journal {
		if err := r.SetJournal(ctx, j.UserID, &s.Journal{
			Timestamp:  j.Timestamp,
			Meal:       j.Meal,
			FoodKey:    j.FoodKey,
			FoodWeight: j.FoodWeight,
		}); err != nil {
			return err
		}
	}

	for _, t := range backup.DayTotalCal {
		if err := r.SetDayTotalCal(ctx, t.UserID, t.Timestamp, t.TotalCal); err != nil {
			return err
		}
	}

	return nil
}
