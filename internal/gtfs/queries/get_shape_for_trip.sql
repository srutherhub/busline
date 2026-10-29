SELECT
  sp.shape_id,
  sp.sequence,
  sp.latitude,
  sp.longitude,
  sp.distance_traveled
FROM
  shape_points sp
WHERE
  sp.feed_id = $1 AND sp.shape_id = $2;