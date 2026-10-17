INSERT INTO
  stop_times (
    feed_id,
    trip_id,
    stop_sequence,
    stop_id,
    arrival_time,
    departure_time,
    stop_headsign,
    pickup_type,
    drop_off_type,
    shape_dist_traveled,
    timepoint
  )
SELECT
  t.feed_id,
  t.id,
  $4,
  s.id,
  $5,
  $6,
  $7,
  $8,
  $9,
  $10,
  $11
FROM
  trips t
  JOIN stop s ON s.feed_id = t.feed_id
  AND s.stop_id = $3
WHERE
  t.feed_id = $1
  AND t.trip_id = $2;