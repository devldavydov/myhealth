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
)
