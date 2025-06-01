package sqlite

import (
	"context"

	s "github.com/devldavydov/myhealth/internal/storage"
)

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
