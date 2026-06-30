CREATE OR REPLACE FUNCTION get_date(input_str TEXT)
RETURNS DATE AS $$
DECLARE
    cleaned_str TEXT := trim(input_str);
BEGIN
    IF cleaned_str = '' OR cleaned_str IS NULL THEN
        RETURN CURRENT_DATE;
    END IF;

    IF cleaned_str ~ '^[+-]?\d+$' THEN
        RETURN CURRENT_DATE + CAST(cleaned_str AS INTEGER);
    END IF;

    BEGIN
        RETURN CAST(cleaned_str AS DATE);
    EXCEPTION WHEN others THEN
        RAISE EXCEPTION 'Неверный формат входных данных: %', input_str;
    END;
END;
$$ LANGUAGE plpgsql;