INSERT INTO stop (
    feed_id, stop_id, code, name, description, latitude, longitude,
    zone_id, url, location_type, timezone, wheelchair_boarding
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
);