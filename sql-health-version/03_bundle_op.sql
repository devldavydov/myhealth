CALL set_bundle(
    'dinner_pack', 
    ARRAY[
        ['beef_steak', '250'], 
        ['broccoli', '100']
    ]
);

CALL delete_bundle('super_lunch');

SELECT * FROM get_all_bundles();