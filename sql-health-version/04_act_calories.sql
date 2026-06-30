CREATE TABLE IF NOT EXISTS act_calories (
    dt        DATE NOT NULL,
    value     REAL NOT NULL,
    PRIMARY KEY (dt)
);

--
--

CREATE OR REPLACE FUNCTION get_act_cal(target_date DATE) 
RETURNS REAL AS $$
DECLARE
    default_calories CONSTANT REAL := 2400;
BEGIN
    RETURN COALESCE(
        (SELECT value FROM act_calories WHERE dt = target_date), 
        default_calories
    );
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE set_act_cal(p_date_str TEXT, act_value REAL) 
AS $$
DECLARE
    v_date DATE := get_date(p_date_str);
BEGIN
    INSERT INTO act_calories (dt, value) 
    VALUES (v_date, act_value)
    ON CONFLICT (dt) 
    DO UPDATE SET value = EXCLUDED.value;
END;
$$ LANGUAGE plpgsql;