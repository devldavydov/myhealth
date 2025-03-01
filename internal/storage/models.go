package storage

import (
	"strings"
	"time"
)

// Unix time in milliseconds
type Timestamp int64

func NewTimestamp(t time.Time) Timestamp {
	return Timestamp(t.UnixMilli())
}

func (r Timestamp) ToTime(tz *time.Location) time.Time {
	return time.UnixMilli(int64(r)).In(tz)
}

type Food struct {
	Key     string
	Name    string
	Brand   string
	Cal100  float64
	Prot100 float64
	Fat100  float64
	Carb100 float64
	Comment string
}

func (r *Food) Validate() bool {
	return r.Key != "" &&
		r.Name != "" &&
		r.Cal100 >= 0 &&
		r.Prot100 >= 0 &&
		r.Fat100 >= 0 &&
		r.Carb100 >= 0
}

type Meal int64

func NewMealFromString(m string) (Meal, error) {
	switch strings.ToUpper(m) {
	case "ЗАВТРАК":
		return Meal(0), nil
	case "ДО ОБЕДА":
		return Meal(1), nil
	case "ОБЕД":
		return Meal(2), nil
	case "ПОЛДНИК":
		return Meal(3), nil
	case "ДО УЖИНА":
		return Meal(4), nil
	case "УЖИН":
		return Meal(5), nil
	}
	return Meal(-1), ErrMealWrong
}

func (r Meal) MustToString() string {
	switch r {
	case 0:
		return "Завтрак"
	case 1:
		return "До обеда"
	case 2:
		return "Обед"
	case 3:
		return "Полдник"
	case 4:
		return "До ужина"
	case 5:
		return "Ужин"
	}

	panic(ErrMealWrong)
}

type Journal struct {
	Timestamp  Timestamp
	Meal       Meal
	FoodKey    string
	FoodWeight float64
}

func (r *Journal) Validate() bool {
	return r.Meal >= 0 &&
		r.FoodKey != "" &&
		r.FoodWeight > 0
}

type Weight struct {
	Timestamp Timestamp
	Value     float64
}

func (r *Weight) Validate() bool {
	return r.Value > 0
}

type UserSettings struct {
	CalLimit float64
}

func (r *UserSettings) Validate() bool {
	return r.CalLimit > 0
}

type Bundle struct {
	Key string
	// Map of bundle data
	// Variants:
	// if food: food_key -> weight > 0
	// if bundle: bundle_key -> 0
	Data map[string]float64
}

func (r *Bundle) Validate() bool {
	if r.Key == "" || len(r.Data) == 0 {
		return false
	}

	for _, v := range r.Data {
		if v < 0 {
			return false
		}
	}

	return true
}

type Sport struct {
	Key     string
	Name    string
	Comment string
}

func (r *Sport) Validate() bool {
	return r.Key != "" &&
		r.Name != ""
}

type SportActivity struct {
	SportKey  string
	Timestamp Timestamp
	Sets      []int64
}

func (r *SportActivity) Validate() bool {
	allPositive := true
	for _, s := range r.Sets {
		if s <= 0 {
			allPositive = false
			break
		}
	}
	return r.SportKey != "" && len(r.Sets) != 0 && allPositive
}

type SportActivityReport struct {
	SportName string
	Timestamp Timestamp
	Sets      []int64
}

type Backup struct {
	Timestamp     Timestamp             `json:"timestamp"`
	Weight        []WeightBackup        `json:"weight"`
	Sport         []SportBackup         `json:"sport"`
	SportActivity []SportActivityBackup `json:"sport_activity"`
	UserSettings  []UserSettingsBackup  `json:"user_settings"`
	Food          []FoodBackup          `json:"food"`
	Bundle        []BundleBackup        `json:"bundle"`
}

type WeightBackup struct {
	UserID    int64     `json:"user_id"`
	Timestamp Timestamp `json:"timestamp"`
	Value     float64   `json:"value"`
}

type SportBackup struct {
	UserID  int64  `json:"user_id"`
	Key     string `json:"key"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type SportActivityBackup struct {
	UserID    int64     `json:"user_id"`
	SportKey  string    `json:"sport_key"`
	Timestamp Timestamp `json:"timestamp"`
	Sets      []int64   `json:"sets"`
}

type UserSettingsBackup struct {
	UserID   int64   `json:"user_id"`
	CalLimit float64 `json:"cal_limit"`
}

type FoodBackup struct {
	UserID  int64   `json:"user_id"`
	Key     string  `json:"key"`
	Name    string  `json:"name"`
	Brand   string  `json:"brand"`
	Cal100  float64 `json:"cal100"`
	Prot100 float64 `json:"prot100"`
	Fat100  float64 `json:"fat100"`
	Carb100 float64 `json:"carb100"`
	Comment string  `json:"comment"`
}

type BundleBackup struct {
	UserID int64              `json:"user_id"`
	Key    string             `json:"key"`
	Data   map[string]float64 `json:"data"`
}
