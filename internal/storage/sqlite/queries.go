package sqlite

const (
	//
	// System.
	//

	_sqlCreateTableSystem = `
	CREATE TABLE IF NOT EXISTS system (
		migration_id INTEGER
	) STRICT;
	`
	_sqlGetLastMigrationID = `
	SELECT migration_id FROM system
	`
	_sqlInsertInitialMigrationID = `
	INSERT INTO system (migration_id) VALUES(0)
	`
	_sqlUpdateLastMigrationID = `
	UPDATE system
	SET migration_id = $1
	`

	//
	// Weight.
	//

	_sqlCreateTableWeight = `
	CREATE TABLE weight (
        user_id   INTEGER NOT NULL,
        timestamp INTEGER NOT NULL,
        value     REAL    NOT NULL,
        PRIMARY KEY (user_id, timestamp)
    ) STRICT
	`

	_sqlGetWeightList = `
	SELECT timestamp, value
    FROM weight
    WHERE
        user_id = $1 AND
        timestamp >= $2 AND
        timestamp <= $3
    ORDER BY
        timestamp
	`

	_sqlWeightBackup = `
	SELECT user_id, timestamp, value
    FROM weight
	ORDER BY user_id, timestamp
	`

	_sqlSetWeight = `
	INSERT INTO weight (user_id, timestamp, value)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id, timestamp) DO
	UPDATE SET value = $3
	`

	_sqlDeleteWeight = `
	DELETE
	FROM weight
	WHERE user_id = $1 AND timestamp = $2
	`

	//
	// Sport.
	//

	_sqlCreateTableSport = `
	CREATE TABLE sport (
	 	user_id  INTEGER NOT NULL,
        key      TEXT NOT NULL,
        name     TEXT NOT NULL,
        comment  TEXT NULL,
		PRIMARY KEY (user_id, key)
    ) STRICT
	`

	_sqlGetSport = `
	SELECT key, name, comment
    FROM sport
    WHERE user_id = $1 AND key = $2
	`

	_sqlGetSportList = `
	SELECT key, name, comment
    FROM sport
    WHERE user_id = $1
	ORDER BY name
	`

	_sqlSetSport = `
	INSERT INTO sport (user_id, key, name, comment)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (user_id, key) DO
	UPDATE SET name = $3, comment = $4
	`

	_sqlDeleteSport = `
	DELETE
	FROM sport
    WHERE user_id = $1 AND key = $2
	`

	_sqlSportBackup = `
	SELECT user_id, key, name, comment
    FROM sport
	ORDER BY user_id, key
	`

	//
	// SportActivity.
	//

	_sqlCreateTableSportActivity = `
	CREATE TABLE sport_activity (
        user_id   INTEGER NOT NULL,
        timestamp INTEGER NOT NULL,
        sport_key TEXT NOT NULL,
        sets      TEXT NOT NULL,
        PRIMARY KEY (user_id, timestamp, sport_key),
        FOREIGN KEY (user_id, sport_key) REFERENCES sport(user_id, key) ON DELETE RESTRICT
    ) STRICT
	`

	_sqlSetSportActivity = `
    INSERT INTO sport_activity (
        user_id, timestamp, sport_key, sets
    )
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (user_id, timestamp, sport_key) DO
    UPDATE SET
        sets = $4	
	`

	_sqlDeleteSportActivity = `
	DELETE FROM sport_activity
    WHERE
        user_id = $1 AND
        timestamp = $2 AND
        sport_key = $3
	`

	_sqlGetSportActivityReport = `
    SELECT sa.timestamp, s.name as sport_name, sa.sets
    FROM
        sport_activity sa,
        sport s
    WHERE
        sa.sport_key = s.key AND	
        s.user_id = $1 AND
		sa.user_id = $1 AND
        sa.timestamp >= $2 AND
        sa.timestamp <= $3
    ORDER BY
        sa.timestamp,
        s.name	
	`

	_sqlSportActivityBackup = `
	SELECT user_id, timestamp, sport_key, sets
	FROM sport_activity
	ORDER BY user_id, timestamp, sport_key
	`

	//
	// UserSettings.
	//

	_sqlCreateTableUserSettings = `
	CREATE TABLE user_settings (
        user_id   INTEGER NOT NULL PRIMARY KEY,
        cal_limit REAL    NOT NULL
    ) STRICT
	`

	_sqlGetUserSettings = `
	SELECT cal_limit
    FROM user_settings
    WHERE user_id = $1
	`

	_sqlSetUserSettings = `
	INSERT INTO user_settings (
        user_id, cal_limit
    )
    VALUES ($1, $2)
    ON CONFLICT (user_id) DO
    UPDATE SET
        cal_limit = $2
	`

	_sqlUserSettingsBackup = `
	SELECT user_id, cal_limit
    FROM user_settings
    ORDER BY user_id
	`

	//
	// Food.
	//

	_sqlCreateTableFood = `
	CREATE TABLE food (
		user_id INTEGER NOT NULL,
        key     TEXT NOT NULL,
        name    TEXT NOT NULL,
        brand   TEXT NULL,
        cal100  REAL NOT NULL,
        prot100 REAL NOT NULL, 
        fat100  REAL NOT NULL,
        carb100 REAL NOT NULL,
        comment TEXT NULL,
		PRIMARY KEY (user_id, key)
    ) STRICT
	`

	_sqlGetFood = `
	SELECT 
        key, name, brand, cal100,
        prot100, fat100, carb100, comment
    FROM food
    WHERE user_id = $1 AND key = $2
	`

	_sqlGetFoodList = `
	SELECT 
        key, name, brand, cal100,
        prot100, fat100, carb100, comment
    FROM food
    WHERE user_id = $1
	ORDER BY name, key
	`

	_sqlFindFood = `
	SELECT 
        key, name, brand, cal100,
        prot100, fat100, carb100, comment
    FROM food
    WHERE
		user_id = $1 AND
		(
        	go_upper(key)     LIKE '%' || $2 || '%' OR
        	go_upper(name)    LIKE '%' || $2 || '%' OR
        	go_upper(brand)   LIKE '%' || $2 || '%' OR
        	go_upper(comment) LIKE '%' || $2 || '%'
		)
    ORDER BY name, key
	`

	_sqlSetFood = `
	INSERT INTO food (
        user_id, key, name, brand, cal100,
        prot100, fat100, carb100, comment
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    ON CONFLICT (user_id, key) DO
    UPDATE SET
        name = $3, brand = $4, cal100 = $5,
        prot100 = $6, fat100 = $7, carb100 = $8,
        comment = $9
	`

	_sqlDeleteFood = `
	DELETE FROM food
    WHERE user_id = $1 AND key = $2
	`

	_sqlFoodBackup = `
	SELECT 
        user_id, key, name, brand, cal100,
        prot100, fat100, carb100, comment
    FROM food
	ORDER BY user_id, key
	`

	//
	// Bundle.
	//

	_sqlCreateTableBundle = `
	    CREATE TABLE bundle (
        user_id    INTEGER NOT NULL,
        key        TEXT    NOT NULL, 
        data       TEXT    NOT NULL,  
        PRIMARY KEY (user_id, key)
    ) STRICT
	`

	_sqlGetBundle = `
	SELECT key, data
    FROM bundle
    WHERE user_id = $1 AND key = $2
	`

	_sqlGetBundleList = `
	SELECT key, data
    FROM bundle
    WHERE user_id = $1
	ORDER BY key
	`

	_sqlSetBundle = `
	INSERT INTO bundle (
        user_id, key, data
    )
    VALUES ($1, $2, $3)
    ON CONFLICT (user_id, key) DO
    UPDATE SET
        data = $3
	`

	_sqlDeleteBundle = `
	DELETE FROM bundle
    WHERE user_id = $1 AND key = $2
	`

	_sqlBundleBackup = `
	SELECT user_id, key, data
	FROM bundle
	ORDER BY user_id, key
	`

	//
	// Journal.
	//

	_sqlCreateTableJournal = `
	CREATE TABLE journal (
        user_id    INTEGER NOT NULL,
        timestamp  INTEGER NOT NULL,
        meal       INTEGER NOT NULL,
        foodkey    TEXT NOT NULL,
        foodweight REAL NOT NULL,
        PRIMARY KEY (user_id, timestamp, meal, foodkey),
        FOREIGN KEY (user_id, foodkey) REFERENCES food(user_id, key) ON DELETE RESTRICT
    ) STRICT
	`

	_sqlCreateTableJournalIndexUserIDFoodKey = `
	CREATE INDEX journal_userid_foodkey ON journal(user_id, foodkey);
	`

	_sqlSetJournal = `
	INSERT INTO journal (
        user_id, timestamp, meal, foodkey, foodweight
    )
    VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT (user_id, timestamp, meal, foodkey) DO
    UPDATE SET
        foodweight = $5
	`

	_sqlDeleteJournal = `
    DELETE FROM journal
    WHERE user_id = $1 AND
          timestamp = $2 AND
          meal = $3 AND
          foodkey = $4
	`

	_sqlDeleteJournalMeal = `
	DELETE FROM journal
	WHERE user_id = $1 AND
		timestamp = $2 AND
		meal = $3
	`

	_sqlJournalFoodStat = `
	SELECT
		min(timestamp) AS first_timestamp,
		max(timestamp) AS last_timestamp,
		sum(foodweight) AS total_weight,
		avg(foodweight) AS avg_weight,
		count(*) AS total_cnt
	FROM journal j
	WHERE
		j.user_id = $1 AND
		j.foodkey = $2
	`

	_sqlGetJournalReport = `
    SELECT
        j.timestamp,
        j.meal,
        j.foodkey,
        f.name AS foodname,
        f.brand AS foodbrand,
        j.foodweight,
        j.foodweight / 100 * f.cal100 AS cal,
        j.foodweight / 100 * f.prot100 AS prot,
        j.foodweight / 100 * f.fat100 AS fat,
        j.foodweight / 100 * f.carb100 AS carb
    FROM journal j, food f
    WHERE
        j.foodkey = f.key AND
		f.user_id = $1 AND
        j.user_id = $1 AND
        j.timestamp >= $2 AND
        j.timestamp <= $3
    ORDER BY
        j.timestamp,
        j.meal,
        f.name
	`

	_sqlGetJournalListForCopy = `
	SELECT foodkey, foodweight
	FROM journal
	WHERE user_id = $1 AND
		timestamp = $2 AND
		meal = $3
	`

	_sqlJournalBackup = `
	SELECT user_id, timestamp, meal, foodkey, foodweight
	FROM journal
	ORDER BY user_id, timestamp, meal, foodkey
	`
)
