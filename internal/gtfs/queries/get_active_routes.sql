SELECT 
  r.id, r.route_id, r.short_name, r.long_name, r.url, r.color, r.text_color
FROM 
  route r
WHERE 
  r.feed_id = $1
  AND EXISTS (
    SELECT 1
    FROM trips t 
    WHERE t.feed_id = $1 
    AND t.route_id = r.id
    AND t.service_id IN
  (
    SELECT 
      c.service_id
    FROM 
      calendar c
    WHERE c.feed_id = $1
      AND c.start_date <= $2
      AND c.end_Date >= $2
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
    )
ORDER BY r.id;