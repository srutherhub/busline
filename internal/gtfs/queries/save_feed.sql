INSERT INTO
    feed_info (
        publisher_name,
        url,
        language,
        start_date,
        end_date,
        version,
        content_hash,
        city
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;