package gtfs

import (
	queries "busline/internal/gtfs/queries"
	"context"
	"fmt"

	"database/sql"
	_ "embed"
)

/* saves gtfs data to the database
 */
type Repository struct {
	db *sql.DB
}

/* initialize the gtfs repository
 */
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) SaveStaticGTFS(ctx context.Context, gtfs *GTFS) error {
	feedId, err := repo.SaveFeed(ctx, *gtfs.FeedInfo, gtfs.StaticHash, string(gtfs.Config.City))

	if err != nil {
		return fmt.Errorf("gtfs: %w", err)
	}

	agencyId, err := repo.SaveAgency(ctx, feedId, *gtfs.Agency)

	if err != nil {
		return fmt.Errorf("gtfs: %w", err)
	}

	for _, route := range gtfs.Routes {
		err = repo.SaveRoute(ctx, feedId, agencyId, route)
		if err != nil {
			return fmt.Errorf("gtfs: %w", err)
		}
	}

	for _, stop := range gtfs.Stops {
		err = repo.SaveStop(ctx, feedId, stop)
		if err != nil {
			return fmt.Errorf("gtfs: %w", err)
		}
	}

	for _, trip := range gtfs.Trips {
		err = repo.SaveTrip(ctx, feedId, trip)
		if err != nil {
			return fmt.Errorf("gtfs: %w", err)
		}
	}

	for _, calendar := range gtfs.Calendars {
		err = repo.SaveCalendar(ctx, feedId, calendar)
		if err != nil {
			return fmt.Errorf("gtfs: %w", err)
		}
	}

	for _, shape := range gtfs.Shapes {
		err = repo.SaveShapePoint(ctx, feedId, shape)
		if err != nil {
			return fmt.Errorf("gtfs: %w", err)
		}
	}

	return err
}

/* checks if a feed already exists in the database
 */
func (repo *Repository) DoesFeedExist(hash string) (bool, error) {
	var exists bool

	err := repo.db.QueryRow(queries.QueryFeedExists, hash).Scan(&exists)

	return exists, err
}

/* inserts feed info into the database and returns the id
 */
func (repo *Repository) SaveFeed(ctx context.Context, feed FeedInfo, contentHash string, city string) (int64, error) {
	var id int64
	err := repo.db.QueryRowContext(ctx, queries.QuerySaveFeed,
		feed.PublisherName,
		feed.Url,
		feed.Language,
		feed.StartDate,
		feed.EndDate,
		feed.Version,
		contentHash,
		city,
	).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("gtfs: save feed info: %w", err)
	}
	return id, nil
}

/* inserts agency into database
 */
func (repo *Repository) SaveAgency(ctx context.Context, feedId int64, agency Agency) (int64, error) {
	var id int64
	err := repo.db.QueryRowContext(ctx, queries.QuerySaveAgency, feedId, agency.Id, agency.Name, agency.Url, agency.Timezone, agency.Language, agency.Phone, agency.FareUrl, agency.Email, agency.BikePolicy).Scan(&id)
	return id, err
}

/* inserts route into database
 */
func (repo *Repository) SaveRoute(ctx context.Context, feedID, agencyID int64, r Route) error {
	_, err := repo.db.ExecContext(ctx, queries.QuerySaveRoute,
		feedID, agencyID, r.Id, r.NetworkId, r.ShortName, r.LongName, r.Description,
		r.Type, r.Url, r.Color, r.TextColor, r.SortOrder,
	)

	return err
}

/* insert a stop into the database
 */
func (repo *Repository) SaveStop(ctx context.Context, feedID int64, s Stop) error {
	_, err := repo.db.ExecContext(ctx, queries.QuerySaveStop,
		feedID, s.Id, s.Code, s.Name, s.Description, s.Latitude, s.Longitude,
		s.ZoneId, s.Url, s.LocationType, s.Timezone, s.WheelchairBoarding,
	)
	return err
}

/* insert a trip to the database
 */
func (repo *Repository) SaveTrip(ctx context.Context, feedID int64, t Trip) error {
	res, err := repo.db.ExecContext(ctx, queries.QuerySaveTrip,
		feedID, t.RouteId, t.ServiceId, t.TripId, t.TripHeadsign, t.TripShortName,
		t.DirectionId, t.BlockId, t.ShapeId, t.WheelchairAccessible, t.BikesAllowed,
	)
	if err != nil {
		return fmt.Errorf("gtfs: save trip %q: %w", t.TripId, err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("gtfs: trip %q references unknown route %q", t.TripId, t.RouteId)
	}
	return nil
}

/* inserts a stop time into the database
 */
func (repo *Repository) SaveStopTime(ctx context.Context, feedID int64, st StopTimes) error {
	res, err := repo.db.ExecContext(ctx, queries.QuerySaveStopTime,
		feedID,
		st.TripId,
		st.StopId,
		st.StopSequence,
		st.ArrivalTime,
		st.DepartureTime,
		st.StopHeadsign,
		st.PickupType,
		st.DropoffType,
		st.ShapeDistanceTraveled,
		st.Timepoint,
	)
	if err != nil {
		return fmt.Errorf("gtfs: save stop time (trip %q, seq %d): %w", st.TripId, st.StopSequence, err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("gtfs: stop time references unknown trip %q or stop %q", st.TripId, st.StopId)
	}
	return nil
}

func (repo *Repository) SaveStopTimesBatch(ctx context.Context, feedId int64, stopTimes []StopTimes) error {
	if len(stopTimes) == 0 {
		return nil
	}

	n := len(stopTimes)

	tripIDs := make([]string, n)
	stopIDs := make([]string, n)
	stopSequences := make([]any, n)
	arrivalTimes := make([]any, n)
	departureTimes := make([]any, n)
	stopHeadsigns := make([]any, n)
	pickupTypes := make([]any, n)
	dropoffTypes := make([]any, n)
	shapeDistances := make([]any, n)
	timepoints := make([]any, n)

	for i, st := range stopTimes {
		tripIDs[i] = st.TripId
		arrivalTimes[i] = st.ArrivalTime
		departureTimes[i] = st.DepartureTime
		stopIDs[i] = st.StopId
		stopSequences[i] = st.StopSequence
		stopHeadsigns[i] = st.StopHeadsign
		pickupTypes[i] = st.PickupType
		dropoffTypes[i] = st.DropoffType
		shapeDistances[i] = st.ShapeDistanceTraveled
		timepoints[i] = st.Timepoint
	}

	_, err := repo.db.ExecContext(ctx, queries.QuerySaveStopTimesBatch,
		feedId,
		tripIDs,
		stopIDs,
		stopSequences,
		arrivalTimes,
		departureTimes,
		stopHeadsigns,
		pickupTypes,
		dropoffTypes,
		shapeDistances,
		timepoints,
	)
	if err != nil {
		return fmt.Errorf("gtfs: failed to batch save stop times: %w", err)
	}

	return nil
}

/* inserts a calendar into the database
 */
func (repo *Repository) SaveCalendar(ctx context.Context, feedID int64, c Calendar) error {
	_, err := repo.db.ExecContext(ctx, queries.QuerySaveCalendar,
		feedID,
		c.ServiceId,
		c.Monday,
		c.Tuesday,
		c.Wednesday,
		c.Thursday,
		c.Friday,
		c.Saturday,
		c.Sunday,
		c.StartDate,
		c.EndDate,
	)
	if err != nil {
		return fmt.Errorf("failed to save calendar %w", err)
	}
	return nil
}

/* inserts a shape point into the database
 */
func (repo *Repository) SaveShapePoint(ctx context.Context, feedID int64, s Shape) error {
	_, err := repo.db.ExecContext(ctx, queries.QuerySaveShapePoint,
		feedID,
		s.Id,
		s.Sequence,
		s.PointLatitude,
		s.PointLongitude,
		s.DistanceTraveled,
	)
	if err != nil {
		return fmt.Errorf("gtfs: save shape point %q seq %d: %w", s.Id, s.Sequence, err)
	}
	return nil
}

/* retrieves the internal feed id by content hash
 */
func (repo *Repository) GetFeedIdByHash(ctx context.Context, hash string) (int64, error) {
	var feedId int64

	err := repo.db.QueryRowContext(ctx, queries.QueryGetFeedIdByHash, hash).Scan(&feedId)

	if err != nil {
		return 0, fmt.Errorf("failed to retrieve feed id %w", err)
	}

	return feedId, nil
}
