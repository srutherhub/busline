package gtfs

import _ "embed"

//go:embed feed_exists.sql
var QueryFeedExists string

//go:embed save_feed.sql
var QuerySaveFeed string

//go:embed save_agency.sql
var QuerySaveAgency string

//go:embed save_route.sql
var QuerySaveRoute string

//go:embed save_stop.sql
var QuerySaveStop string

//go:embed save_trip.sql
var QuerySaveTrip string

//go:embed save_stop_time.sql
var QuerySaveStopTime string

//go:embed save_calendar.sql
var QuerySaveCalendar string

//go:embed save_shape_points.sql
var QuerySaveShapePoint string

//go:embed save_stop_times_batch.sql
var QuerySaveStopTimesBatch string

//go:embed get_feedid_by_hash.sql
var QueryGetFeedIdByHash string

//go:embed get_feed_info.sql
var QueryGetFeedInfo string

//go:embed get_active_routes.sql
var QueryGetActiveRoutes string

//go:embed get_nearby_stops.sql
var QueryGetNearybyStops string

//go:embed get_trips_for_stop.sql
var QueryGetTripsForStop string

//go:embed get_shape_for_trip.sql
var QueryGetShapeForTrip string
