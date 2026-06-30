CREATE TABLE IF NOT EXISTS weight (
    dt        DATE NOT NULL,
    value     REAL NOT NULL,
    PRIMARY KEY (dt)
);

--
--

CREATE OR REPLACE PROCEDURE set_weight(p_dt TEXT, p_value REAL)
AS $$
BEGIN
    INSERT INTO weight (dt, value)
    VALUES (p_dt::DATE, p_value)
    ON CONFLICT (dt) 
    DO UPDATE SET value = EXCLUDED.value;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE delete_weight(p_dt TEXT)
AS $$
BEGIN
    DELETE FROM weight
    WHERE dt = p_dt::DATE;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_weight(
    p_start_date TEXT,
    p_end_date TEXT
)
RETURNS RETURNS SETOF weight
AS $$
DECLARE
    v_start DATE := get_date(p_start_date);
    v_end DATE := get_date(p_end_date);
BEGIN
    RETURN QUERY
    SELECT w.dt, w.value
    FROM weight w
    WHERE w.dt BETWEEN v_start AND v_end
    ORDER BY w.dt DESC;
END;
$$ LANGUAGE plpgsql;

