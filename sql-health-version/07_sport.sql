CREATE TABLE IF NOT EXISTS sport (
    key     TEXT NOT NULL,
    name    TEXT NOT NULL,
    unit    TEXT NOT NULL,
    comment TEXT NOT NULL,
    PRIMARY KEY (key)
);

CREATE TABLE IF NOT EXISTS sport_activity (
    dt        DATE NOT NULL,
    sport_key TEXT NOT NULL,
    sets      REAL[] NOT NULL,
    PRIMARY KEY (dt, sport_key),
    FOREIGN KEY (sport_key) REFERENCES sport(key) ON DELETE RESTRICT
);

--
--

CREATE OR REPLACE PROCEDURE set_sport(
    p_key TEXT,
    p_name TEXT,
    p_unit TEXT,
    p_comment TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO sport(key, name, unit, comment)
    VALUES (p_key, p_name, p_unit, p_comment)
    ON CONFLICT (key) 
    DO UPDATE SET
        name    = EXCLUDED.name,
        unit    = EXCLUDED.unit,
        comment = EXCLUDED.comment;
END;
$$;

CREATE OR REPLACE PROCEDURE del_sport(
    p_key TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM sport
    WHERE key = p_key;
END;
$$;

CREATE OR REPLACE PROCEDURE set_sport_act(
    p_date_str TEXT,
    p_sport_key TEXT,
    p_sets REAL[]
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_date DATE := get_date(p_date_str);
BEGIN
    INSERT INTO sport_activity (dt, sport_key, sets)
    VALUES (v_date, p_sport_key, p_sets)
    ON CONFLICT (dt, sport_key) 
    DO UPDATE SET sets = EXCLUDED.sets;
END;
$$;

CREATE OR REPLACE PROCEDURE del_sport_act(
    p_date_str TEXT,
    p_sport_key TEXT
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_date DATE := get_date(p_date_str);
BEGIN
    DELETE FROM sport_activity
    WHERE dt = v_date AND sport_key = p_sport_key;
END;
$$;