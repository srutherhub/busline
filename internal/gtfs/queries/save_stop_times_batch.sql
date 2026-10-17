INSERT INTO stop_times (
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
    $1::bigint,    
    t.id,                
    NULLIF(i.stop_sequence, '')::integer,
    s.id,                 
    NULLIF(i.arrival_time, ''),
    NULLIF(i.departure_time, ''),
    NULLIF(i.stop_headsign, ''),
    NULLIF(i.pickup_type, '')::smallint,
    NULLIF(i.drop_off_type, '')::smallint,
    NULLIF(i.shape_dist_traveled, '')::double precision,
    NULLIF(i.timepoint, '')::smallint
FROM unnest(      
    $2::text[],           
    $3::text[],           
    $4::text[],           
    $5::text[],           
    $6::text[],        
    $7::text[],           
    $8::text[],          
    $9::text[],          
    $10::text[],          
    $11::text[]          
) AS i(
    trip_id, stop_id, stop_sequence, arrival_time, 
    departure_time, stop_headsign, pickup_type, drop_off_type, 
    shape_dist_traveled, timepoint
)
JOIN trips t ON t.feed_id = $1::bigint AND t.trip_id = i.trip_id
JOIN stop s ON s.feed_id = $1::bigint AND s.stop_id = i.stop_id;