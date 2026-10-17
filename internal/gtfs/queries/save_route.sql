INSERT INTO route (
    feed_id, agency_id, route_id, network_id, short_name, long_name, description,
    type, url, color, text_color, sort_order
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
);