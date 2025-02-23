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
)
