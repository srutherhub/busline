CREATE TABLE feed_info (
    id             bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    publisher_name text,
    url            text,
    language       text,
    start_date     date,
    end_date       date,
    version        text,
    imported_at    timestamptz NOT NULL DEFAULT now(),
    content_hash   text NOT NULL UNIQUE
);

CREATE TABLE agency (
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    feed_id     bigint NOT NULL REFERENCES feed_info(id) ON DELETE CASCADE,
    agency_id   text,
    name        text NOT NULL,
    url         text NOT NULL,
    timezone    text NOT NULL,
    language    text,
    phone       text,
    fare_url    text,
    email       text,
    bike_policy text,
    UNIQUE (feed_id, agency_id),
    UNIQUE (feed_id, id)
);

CREATE UNIQUE INDEX agency_one_null_per_feed_idx
    ON agency (feed_id) WHERE agency_id IS NULL;

CREATE TABLE route (
    id                  bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    feed_id             bigint NOT NULL REFERENCES feed_info(id) ON DELETE CASCADE,
    agency_id           bigint NOT NULL,
    route_id            text NOT NULL,
    network_id          text,
    short_name          text,
    long_name           text,
    description         text,
    type                integer NOT NULL,
    url                 text,
    color               text,
    text_color          text ,
    sort_order          integer,
    continuous_pickup   smallint,
    continuous_drop_off smallint,
    service_name        text,
    bikes_allowed       text,
    UNIQUE (feed_id, route_id),
    UNIQUE (feed_id, id),
    FOREIGN KEY (feed_id, agency_id) REFERENCES agency (feed_id, id) ON DELETE CASCADE,
    CHECK (short_name IS NOT NULL OR long_name IS NOT NULL)
);

CREATE INDEX route_agency_idx ON route (feed_id, agency_id);

CREATE TABLE stop (
    id                  bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    feed_id             bigint NOT NULL REFERENCES feed_info(id) ON DELETE CASCADE,
    stop_id             text NOT NULL,
    code                text,
    name                text,
    description         text,
    latitude            double precision CHECK (latitude  BETWEEN -90  AND 90),
    longitude           double precision CHECK (longitude BETWEEN -180 AND 180),
    zone_id             text,
    url                 text,
    location_type       smallint NOT NULL DEFAULT 0 CHECK (location_type BETWEEN 0 AND 4),
    parent_station_id   bigint,
    timezone            text,
    wheelchair_boarding smallint NOT NULL DEFAULT 0 CHECK (wheelchair_boarding BETWEEN 0 AND 2),
    sign_type           text,
    jurisdiction        text,
    relative_position   text,
    cardinal_direction  text,
    primary_street      text,
    address_range       text,
    cross_location      text,
    UNIQUE (feed_id, stop_id),
    UNIQUE (feed_id, id),
    FOREIGN KEY (feed_id, parent_station_id) REFERENCES stop (feed_id, id),
    CHECK (location_type > 2 OR (name IS NOT NULL AND latitude IS NOT NULL AND longitude IS NOT NULL))
);

CREATE INDEX stop_parent_idx ON stop (feed_id, parent_station_id);

CREATE TABLE trips (
    id                    bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    feed_id               bigint NOT NULL,
    route_id              bigint NOT NULL,
    service_id            text NOT NULL,     -- GTFS service_id (calendar link comes later)
    trip_id               text NOT NULL,     -- GTFS trip_id
    headsign              text,
    short_name            text,
    direction_id          smallint CHECK (direction_id IN (0, 1)),
    block_id              text,
    shape_id              text,
    wheelchair_accessible smallint NOT NULL DEFAULT 0 CHECK (wheelchair_accessible BETWEEN 0 AND 2),
    bikes_allowed         smallint NOT NULL DEFAULT 0 CHECK (bikes_allowed BETWEEN 0 AND 2),

    UNIQUE (feed_id, trip_id),
    UNIQUE (feed_id, id),
    FOREIGN KEY (feed_id, route_id) REFERENCES route (feed_id, id) ON DELETE CASCADE
);

CREATE INDEX trips_route_idx ON trips (feed_id, route_id);
CREATE INDEX trips_service_idx ON trips (feed_id, service_id);

CREATE TABLE stop_times (
    feed_id         bigint   NOT NULL,
    trip_id         bigint   NOT NULL,
    stop_sequence   integer  NOT NULL CHECK (stop_sequence >= 0),
    stop_id         bigint   NOT NULL,

    arrival_time    text ,
    departure_time  text,
    stop_headsign   text,
    pickup_type     smallint NOT NULL DEFAULT 0 CHECK (pickup_type   BETWEEN 0 AND 3),
    drop_off_type   smallint NOT NULL DEFAULT 0 CHECK (drop_off_type BETWEEN 0 AND 3),
    shape_dist_traveled double precision CHECK (shape_dist_traveled >= 0),
    timepoint       smallint NOT NULL DEFAULT 1 CHECK (timepoint IN (0, 1)),

    PRIMARY KEY (feed_id, trip_id, stop_sequence),
    FOREIGN KEY (feed_id, trip_id) REFERENCES trips (feed_id, id) ON DELETE CASCADE,
    FOREIGN KEY (feed_id, stop_id) REFERENCES stop (feed_id, id)
);

CREATE INDEX stop_times_stop_idx ON stop_times (feed_id, stop_id);

CREATE TABLE calendar (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    feed_id    bigint NOT NULL REFERENCES feed_info(id) ON DELETE CASCADE,
    service_id text   NOT NULL,          -- GTFS service_id
    monday     boolean NOT NULL,
    tuesday    boolean NOT NULL,
    wednesday  boolean NOT NULL,
    thursday   boolean NOT NULL,
    friday     boolean NOT NULL,
    saturday   boolean NOT NULL,
    sunday     boolean NOT NULL,
    start_date date NOT NULL,
    end_date   date NOT NULL,

    UNIQUE (feed_id, service_id),
    CHECK (end_date >= start_date)
);

CREATE TABLE shape_points (
    feed_id             bigint NOT NULL REFERENCES feed_info(id) ON DELETE CASCADE,
    shape_id            text   NOT NULL,   -- GTFS shape_id
    sequence            integer NOT NULL CHECK (sequence >= 0),
    latitude            double precision NOT NULL CHECK (latitude  BETWEEN -90  AND 90),
    longitude           double precision NOT NULL CHECK (longitude BETWEEN -180 AND 180),
    distance_traveled   double precision CHECK (distance_traveled >= 0),

    PRIMARY KEY (feed_id, shape_id, sequence)
);