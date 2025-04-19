package storage

type Backup struct {
	Timestamp     Timestamp             `json:"timestamp"`
	Weight        []WeightBackup        `json:"weight"`
	Sport         []SportBackup         `json:"sport"`
	SportActivity []SportActivityBackup `json:"sport_activity"`
	UserSettings  []UserSettingsBackup  `json:"user_settings"`
	Food          []FoodBackup          `json:"food"`
	Bundle        []BundleBackup        `json:"bundle"`
	Journal       []JournalBackup       `json:"journal"`
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
	Sets      []float64 `json:"sets"`
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

type JournalBackup struct {
	UserID     int64     `json:"user_id"`
	Timestamp  Timestamp `json:"timestamp"`
	Meal       Meal      `json:"meal"`
	FoodKey    string    `json:"food_key"`
	FoodWeight float64   `json:"food_weight"`
}
