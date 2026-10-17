INSERT INTO
  trips (
    feed_id,
    route_id,
    service_id,
    trip_id,
    headsign,
    short_name,
    direction_id,
    block_id,
    shape_id,
    wheelchair_accessible,
    bikes_allowed
  )
SELECT
  route.feed_id,
  route.id,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $9,
  $10,
  $11
FROM
  route
WHERE
  route.feed_id = $1
  AND route.route_id = $2;