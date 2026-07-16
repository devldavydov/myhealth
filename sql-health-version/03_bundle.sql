CREATE TABLE IF NOT EXISTS bundle (
    key     TEXT NOT NULL,
    foodkey TEXT NOT NULL,
    weight  REAL NOT NULL,
    PRIMARY KEY (key, foodkey),
    FOREIGN KEY (foodkey) REFERENCES food(key) ON DELETE RESTRICT
);

--
--

CREATE OR REPLACE PROCEDURE set_bundle(
    p_bundle_key TEXT,
    p_food_array TEXT[] -- Формат: ARRAY['food_key:150', 'bundle_key']
) 
LANGUAGE plpgsql AS $$
DECLARE
    v_key   TEXT;
    v_value TEXT;
    item    TEXT;
BEGIN
    DELETE FROM bundle 
    WHERE key = p_bundle_key;

    FOREACH item IN ARRAY p_food_array
    LOOP
        IF item LIKE '%:%' THEN
            -- Добавление еды
            v_key = split_part(item, ':', 1);
            v_value = split_part(item, ':', 2);

            IF NOT EXISTS (SELECT 1 FROM food WHERE key = v_key) THEN
                RAISE EXCEPTION 'Еда с ключом "%" не найдена в базе данных.', v_key;
            END IF;

            INSERT INTO bundle (key, foodkey, weight)
            SELECT
                p_bundle_key,
                v_key,
                v_value::REAL
            ON CONFLICT (key, foodkey)
            DO UPDATE SET weight = EXCLUDED.weight;
        ELSE
            -- Добавление бандла
            v_key = item;

            IF NOT EXISTS (SELECT 1 FROM bundle WHERE key = v_key) THEN
                RAISE EXCEPTION 'Бандл с ключом "%" не найден в базе данных.', v_key;
            END IF;

            INSERT INTO bundle (key, foodkey, weight)
            SELECT 
                p_bundle_key,
                b.foodkey,
                b.weight
            FROM bundle b
            WHERE b.key = v_key
            ON CONFLICT (key, foodkey)
            DO UPDATE SET weight = EXCLUDED.weight;
        END IF;
    END LOOP;
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