INSERT INTO
  agency (
    feed_id,
    agency_id,
    name,
    url,
    timezone,
    language,
    phone,
    fare_url,
    email,
    bike_policy
  )
VALUES
  (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
  ) RETURNING id;