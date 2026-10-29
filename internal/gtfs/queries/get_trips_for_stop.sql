WITH active_services AS (
  SELECT c.service_id
  FROM calendar c
  WHERE c.feed_id = $1
    AND c.start_date <= $2
    AND c.end_date >= $2
    AND CASE $3::int
          WHEN 0 THEN c.sunday
          WHEN 1 THEN c.monday
          WHEN 2 THEN c.tuesday
          WHEN 3 THEN c.wednesday
          WHEN 4 THEN c.thursday
          WHEN 5 THEN c.friday
          WHEN 6 THEN c.saturday
        END
)
SELECT
  t.id,
  t.trip_id,
  t.service_id,
  t.headsign,
  t.direction_id,
  t.shape_id,
  t.wheelchair_accessible,
  r.route_id,
  r.short_name,
  r.long_name,
  r.color,
  r.text_color,
  r.type,
  st.arrival_time,
  st.departure_time,
  st.stop_sequence,
  st.shape_dist_traveled
FROM stop_times st
JOIN trips t  ON t.id = st.trip_id
JOIN route r ON r.id = t.route_id
JOIN active_services a ON a.service_id = t.service_id
WHERE st.feed_id = $1
  AND t.feed_id = $1
  AND st.stop_id = $4
ORDER BY st.departure_time;