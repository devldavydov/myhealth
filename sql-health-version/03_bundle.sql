CREATE TABLE IF NOT EXISTS bundle (
    key     TEXT NOT NULL,
    foodkey TEXT NOT NULL,
    weight  REAL NOT NULL,
    FOREIGN KEY (foodkey) REFERENCES food(key) ON DELETE RESTRICT
);

--
--

CREATE OR REPLACE PROCEDURE set_bundle(
    p_bundle_key TEXT,
    p_food_array TEXT[][] -- Формат: ARRAY[['food_key_1', '150'], ['food_key_2', '50']]
) 
LANGUAGE plpgsql AS $$
BEGIN
    IF p_food_array IS NOT NULL AND array_length(p_food_array, 2) <> 2 THEN
        RAISE EXCEPTION 'Массив должен быть двумерным и содержать ровно две колонки (ключ еды и вес).';
    END IF;

    DELETE FROM bundle 
    WHERE key = p_bundle_key;

    INSERT INTO bundle (key, foodkey, weight)
    SELECT 
        p_bundle_key,
        p_food_array[i][1] AS foodkey,
        p_food_array[i][2]::REAL AS weight
    FROM generate_subscripts(p_food_array, 1) AS i;
END;
$$;

CREATE OR REPLACE PROCEDURE del_bundle(
    p_bundle_key TEXT
) 
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM bundle 
    WHERE key = p_bundle_key;
END;
$$;

CREATE OR REPLACE FUNCTION get_all_bundles()
RETURNS TABLE (
    bundle_key TEXT,
    components TEXT
) 
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT 
        b.key AS bundle_key,
        string_agg(b.foodkey || ':' || b.weight::TEXT, ',') AS components
    FROM bundle b
    GROUP BY b.key
    ORDER BY b.key;
END;
$$;