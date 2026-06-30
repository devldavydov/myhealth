CALL set_journal(
    '2026-06-30', 
    'завтрак', 
    ARRAY[['apple', '150.0'], ['oatmeal', '60.5']]
);

CALL set_journal_bundle(
    '2026-06-30', 
    'обед', 
    'два_бут_кс'
);

SELECT * FROM get_journal('2026-06-30');
SELECT * FROM get_journal('');