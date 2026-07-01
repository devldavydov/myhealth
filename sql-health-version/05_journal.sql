CREATE TYPE meal_type AS ENUM (
    'завтрак', 
    'до обеда', 
    'обед', 
    'полдник', 
    'до ужина', 
    'ужин'
);

CREATE TABLE IF NOT EXISTS journal (
    dt         DATE      NOT NULL,
    meal       meal_type NOT NULL,
    foodkey    TEXT      NOT NULL,
    foodweight REAL      NOT NULL,
    PRIMARY KEY (dt, meal, foodkey),
    FOREIGN KEY (foodkey) REFERENCES food(key) ON DELETE RESTRICT
);

--
--

CREATE OR REPLACE PROCEDURE set_journal(
    p_date_str TEXT,
    p_meal meal_type,
    p_food_array TEXT[][] -- Формат: ARRAY[['food_key_1', '150'], ['food_key_2', '50']]
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_date DATE := get_date(p_date_str);
BEGIN
    IF p_food_array IS NOT NULL AND array_length(p_food_array, 2) <> 2 THEN
        RAISE EXCEPTION 'Массив должен быть двумерным и содержать ровно две колонки (ключ еды и вес).';
    END IF;

    INSERT INTO journal (dt, meal, foodkey, foodweight)
    SELECT 
        v_date,
        p_meal,
        p_food_array[i][1],
        p_food_array[i][2]::REAL
    FROM generate_subscripts(p_food_array, 1) AS i
    ON CONFLICT (dt, meal, foodkey) 
    DO UPDATE SET foodweight = EXCLUDED.foodweight;
END;
$$;

CREATE OR REPLACE PROCEDURE set_journal_bundle(
    p_date_str TEXT,
    p_meal meal_type,
    p_bundle_key TEXT
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_date DATE := get_date(p_date_str);
BEGIN
    IF NOT EXISTS (SELECT 1 FROM bundle WHERE key = p_bundle_key) THEN
        RAISE EXCEPTION 'Бандл с ключом "%" не найден в базе данных.', p_bundle_key;
    END IF;

    INSERT INTO journal (dt, meal, foodkey, foodweight)
    SELECT 
        v_date,
        p_meal,
        b.foodkey,
        b.weight
    FROM bundle b
    WHERE b.key = p_bundle_key
    ON CONFLICT (dt, meal, foodkey) 
    DO UPDATE SET foodweight = EXCLUDED.foodweight;
END;
$$;

CREATE OR REPLACE PROCEDURE del_journal(
    p_date_str TEXT,
    p_meal meal_type,
    p_food_key TEXT
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_date DATE := get_date(p_date_str);
BEGIN
    DELETE FROM journal
    WHERE
        dt = v_date AND
        meal = p_meal AND
        foodkey = p_food_key;
END;
$$;

CREATE OR REPLACE PROCEDURE cp_journal(
    p_date_from_str TEXT,
    p_meal_from meal_type,
    p_date_to_str TEXT,
    p_meal_to meal_type
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_date_from DATE := get_date(p_date_from_str);
    v_date_to DATE := get_date(p_date_to_str);
BEGIN
    INSERT INTO journal
    SELECT
        v_date_to,
        p_meal_to,
        foodkey,
        foodweight
    FROM journal
    WHERE
        dt = v_date_from AND
        meal = p_meal_from
    ON CONFLICT (dt, meal, foodkey)
    DO UPDATE SET foodweight = EXCLUDED.foodweight;
END;
$$;

CREATE OR REPLACE FUNCTION get_journal(p_date_str TEXT)
RETURNS TABLE (
    meal_block TEXT,
    food_name TEXT,
    food_weight TEXT,
    calories TEXT,
    proteins TEXT,
    fats TEXT,
    carbs TEXT
) 
LANGUAGE plpgsql
AS $$
DECLARE
    v_date DATE := get_date(p_date_str);
    cal_limit REAL := get_act_cal(v_date);
BEGIN
    RETURN QUERY
    WITH detailed_data AS (
        SELECT 
            j.meal,
            CASE 
                WHEN trim(f.brand) <> '' 
                THEN  f.name || ' [' || f.key || '] - ' || f.brand
                ELSE f.name || ' [' || f.key || ']'
            END::TEXT AS f_name,
            j.foodweight,
            ROUND(((f.cal100 * j.foodweight) / 100.0)::NUMERIC, 2) AS cal,
            ROUND(((f.prot100 * j.foodweight) / 100.0)::NUMERIC, 2) AS prot,
            ROUND(((f.fat100 * j.foodweight) / 100.0)::NUMERIC, 2) AS fat,
            ROUND(((f.carb100 * j.foodweight) / 100.0)::NUMERIC, 2) AS carb
        FROM journal j
        JOIN food f ON j.foodkey = f.key
        WHERE j.dt = v_date
    ), aggregated_data AS (
        SELECT
            CASE
                WHEN GROUPING(d.meal) = 1 THEN ''
                WHEN GROUPING(d.f_name) = 1 THEN ''
                ELSE d.meal::TEXT
            END AS meal_block,
            CASE 
                WHEN GROUPING(d.meal) = 1 THEN '*** ИТОГО ЗА ДЕНЬ ***'
                WHEN GROUPING(d.f_name) = 1 THEN '*** ПОДЫТОГ (' || d.meal::TEXT || ') ***'
                ELSE d.f_name
            END::TEXT AS food_name,
            CASE
                WHEN GROUPING(d.meal) = 1 THEN ''
                WHEN GROUPING(d.f_name) = 1 THEN ''
                ELSE SUM(d.foodweight)::TEXT
            END AS food_weight,
            SUM(d.cal)::NUMERIC(10,2) AS calories,
            SUM(d.prot)::NUMERIC(10,2) AS proteins,
            SUM(d.fat)::NUMERIC(10,2) AS fats,
            SUM(d.carb)::NUMERIC(10,2) AS carbs
        FROM detailed_data d
        GROUP BY GROUPING SETS (
            (d.meal, d.f_name),
            (d.meal),
            ()
        )
        ORDER BY 
            GROUPING(d.meal) ASC,
            d.meal ASC,
            GROUPING(d.f_name) ASC,
            d.f_name ASC
    ), total_pfc AS (
        SELECT
            ROUND((cal_limit - t.calories)::NUMERIC, 2) AS diff_cal,
            t.proteins AS total_p,
            t.fats AS total_f,
            t.carbs AS total_c
        FROM
            aggregated_data t
        WHERE
            t.food_name = '*** ИТОГО ЗА ДЕНЬ ***'
    )
    SELECT
        t.meal_block,
        t.food_name,
        t.food_weight,
        t.calories::TEXT,
        t.proteins::TEXT,
        t.fats::TEXT,
        t.carbs::TEXT
    FROM
        aggregated_data t
    UNION ALL
    SELECT
        '',
        '*** ЛИМИТ ККАЛ: ' || cal_limit::TEXT || ', % БЖУ ***',
        '',
        CASE
            WHEN diff_cal < 0 THEN '- ' || diff_cal::TEXT
            WHEN diff_cal = 0 THEN diff_cal::TEXT
            ELSE '+ ' || diff_cal::TEXT
        END,
        ROUND((total_p / (total_p + total_f + total_c) * 100)::NUMERIC, 1)::TEXT || '%',
        ROUND((total_f / (total_p + total_f + total_c) * 100)::NUMERIC, 1)::TEXT || '%',
        ROUND((total_c / (total_p + total_f + total_c) * 100)::NUMERIC, 1)::TEXT || '%'
    FROM
        total_pfc
    ;
END;
$$;
