INSERT INTO trip_updates (
  feed_timestamp, trip_id, route_id, trip_schedule_status,
  vehicle_id, stop_id, stop_sequence, arrival_time, departure_time,
  stop_schedule_status, prediction_timestamp
)
SELECT
  $1::timestamptz,
  t.trip_id,
  NULLIF(t.route_id, ''),
  t.trip_schedule_status,
  NULLIF(t.vehicle_id, ''),
  t.stop_id,
  t.stop_sequence,
  to_timestamp(NULLIF(t.arrival_time, 0)),
  to_timestamp(NULLIF(t.departure_time, 0)),
  t.stop_schedule_status,
  to_timestamp(NULLIF(t.prediction_timestamp, 0))
FROM unnest(
  $2::text[], $3::text[], $4::text[], $5::text[], $6::text[],
  $7::int[], $8::bigint[], $9::bigint[], $10::text[], $11::bigint[]
) AS t(trip_id, route_id, trip_schedule_status, vehicle_id, stop_id,
       stop_sequence, arrival_time, departure_time, stop_schedule_status,
       prediction_timestamp)
ON CONFLICT (feed_timestamp, trip_id, stop_sequence) DO NOTHING