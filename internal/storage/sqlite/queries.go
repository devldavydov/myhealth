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
    ) STRICT;
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

	_sqlSetWeight = `
	INSERT INTO weight(user_id, timestamp, value)
	VALUES ($1, $2, $3)
	ON CONFLICT (uuser_id, timestamp) DO
	UPDATE SET value = $3
	`

	_sqlDeleteWeight = `
	DELETE
	FROM weight
	WHERE user_id = $1 AND timestamp = $2
	`
)
