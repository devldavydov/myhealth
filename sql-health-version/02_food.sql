CREATE TABLE IF NOT EXISTS food (
    key     TEXT NOT NULL,
    name    TEXT NOT NULL,
    brand   TEXT NOT NULL,
    cal100  REAL NOT NULL,
    prot100 REAL NOT NULL, 
    fat100  REAL NOT NULL,
    carb100 REAL NOT NULL,
    comment TEXT NOT NULL,
    PRIMARY KEY (key)
);

--
--

CREATE OR REPLACE PROCEDURE set_food(
    p_key     TEXT,
    p_name    TEXT,
    p_brand   TEXT,
    p_cal100  REAL,
    p_prot100 REAL,
    p_fat100  REAL,
    p_carb100 REAL,
    p_comment TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO food (key, name, brand, cal100, prot100, fat100, carb100, comment)
    VALUES (p_key, p_name, p_brand, p_cal100, p_prot100, p_fat100, p_carb100, p_comment)
    ON CONFLICT (key) 
    DO UPDATE SET
        name    = EXCLUDED.name,
        brand   = EXCLUDED.brand,
        cal100  = EXCLUDED.cal100,
        prot100 = EXCLUDED.prot100,
        fat100  = EXCLUDED.fat100,
        carb100 = EXCLUDED.carb100,
        comment = EXCLUDED.comment;
END;
$$;

CREATE OR REPLACE PROCEDURE set_food_by_weight(
    p_key     TEXT,
    p_name    TEXT,
    p_brand   TEXT,
    p_weight  REAL, -- Вес в граммах, для которого указаны КБЖУ
    p_cal     REAL, -- Калории на указанный вес
    p_prot    REAL, -- Белки на указанный вес
    p_fat     REAL, -- Жиры на указанный вес
    p_carb    REAL, -- Углеводы на указанный вес
    p_comment TEXT
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_cal100  REAL;
    v_prot100 REAL;
    v_fat100  REAL;
    v_carb100 REAL;
BEGIN
    IF p_weight IS NULL OR p_weight <= 0 THEN
        RAISE EXCEPTION 'Вес должен быть больше 0. Передано: %', p_weight;
    END IF;

    v_cal100  := ROUND(((p_cal / p_weight) * 100)::numeric, 2)::real;
    v_prot100 := ROUND(((p_prot / p_weight) * 100)::numeric, 2)::real;
    v_fat100  := ROUND(((p_fat / p_weight) * 100)::numeric, 2)::real;
    v_carb100 := ROUND(((p_carb / p_weight) * 100)::numeric, 2)::real;

    INSERT INTO food (key, name, brand, cal100, prot100, fat100, carb100, comment)
    VALUES (p_key, p_name, p_brand, v_cal100, v_prot100, v_fat100, v_carb100, p_comment)
    ON CONFLICT (key) 
    DO UPDATE SET
        name    = EXCLUDED.name,
        brand   = EXCLUDED.brand,
        cal100  = EXCLUDED.cal100,
        prot100 = EXCLUDED.prot100,
        fat100  = EXCLUDED.fat100,
        carb100 = EXCLUDED.carb100,
        comment = EXCLUDED.comment;
END;
$$;

CREATE OR REPLACE PROCEDURE del_food(p_key TEXT)
AS $$
BEGIN
    DELETE FROM food
    WHERE key = p_key;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION search_food(
    p_pattern TEXT
)
RETURNS SETOF food
LANGUAGE plpgsql
AS $$
DECLARE
    v_search TEXT;
BEGIN
    v_search := '%' || COALESCE(TRIM(p_pattern), '') || '%';

    RETURN QUERY
    SELECT key, name, brand, cal100, prot100, fat100, carb100, comment
    FROM food
    WHERE key     ILIKE v_search
       OR name    ILIKE v_search
       OR brand   ILIKE v_search
       OR comment ILIKE v_search
    ORDER BY name ASC;
END;
$$;