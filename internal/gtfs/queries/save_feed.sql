INSERT INTO
    feed_info (
        publisher_name,
        url,
        language,
        start_date,
        end_date,
        version,
        content_hash
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7) RETURNING id;