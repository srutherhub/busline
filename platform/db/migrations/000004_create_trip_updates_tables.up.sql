CREATE TABLE trip_updates (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  feed_timestamp timestamptz NOT NULL,
  trip_id TEXT NOT NULL,
  route_id TEXT,
  trip_schedule_status TEXT,
  vehicle_id TEXT,
  stop_id TEXT NOT NULL,
  stop_sequence INT NOT NULL,
  arrival_time timestamptz,
  departure_time timestamptz,
  stop_schedule_status TEXT,
  prediction_timestamp timestamptz,
  imported_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (feed_timestamp, trip_id, stop_sequence)
);

CREATE INDEX trip_updates_trip_idx ON trip_updates (trip_id, stop_sequence, feed_timestamp);
CREATE INDEX trip_updates_stop_idx ON trip_updates (stop_id, arrival_time);