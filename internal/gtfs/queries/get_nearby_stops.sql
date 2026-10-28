SELECT
  id, stop_id, code, name, latitude, longitude, location_type, wheelchair_boarding,
  ST_Distance(
    ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)::geography,
    ST_SetSRID(ST_MakePoint($3, $2), 4326)::geography
  ) AS meters
FROM stop s
WHERE s.feed_id = $1
ORDER BY meters
LIMIT 20;