package gtfs

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

/* struct that implements the general transit feed specification
 */
type GTFS struct {
	Config     GTFSConfig
	Agency     *Agency
	Routes     []Route
	FeedInfo   *FeedInfo
	Stops      []Stop
	Shapes     []Shape
	Calendars  []Calendar
	Trips      []Trip
	StaticHash string
}

/* config options for the gtfs struct
 */
type GTFSConfig struct {
	StaticUrl string
	City      City
}

/* initializes a new gtfs struct
* @returns
 */
func NewGTFS(config GTFSConfig) (*GTFS, error) {
	if config.StaticUrl == "" {
		return nil, errors.New("gtfs: static url cannot be null")
	}

	if config.City == "" {
		return nil, fmt.Errorf("gtfs: city cannot be null")
	}

	return &GTFS{Config: config}, nil
}

/* downloads the static gtfs files from the static url.
 */
func (gtfs *GTFS) DownloadStatic(ctx context.Context) error {
	url := gtfs.Config.StaticUrl

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("gtfs: failed to generate request %w", err)
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("gtfs: failed to download static files %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gtfs: unexpected status code when downloading static files %s", resp.Status)
	}

	GTFS_DATA_DIR, err := gtfs.GetLocalDataDir()

	if err != nil {
		return fmt.Errorf("gtfs: failed to get local data directory %w", err)
	}

	if err := os.MkdirAll(GTFS_DATA_DIR, 0o755); err != nil {
		return fmt.Errorf("gtfs: failed to create data dir: %w", err)
	}

	tmp, err := os.CreateTemp(GTFS_DATA_DIR, gtfs.GetStaticFileName()+".tmp")

	if err != nil {
		return fmt.Errorf("gtfs: failed to create tmp file %w", err)
	}

	defer os.Remove(tmp.Name())

	hasher := sha256.New()

	if _, err := io.Copy(io.MultiWriter(tmp, hasher), resp.Body); err != nil {
		tmp.Close()
		return fmt.Errorf("gtfs: failed to write zip: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("gtfs: failed to close tmp file: %w", err)
	}

	zr, err := zip.OpenReader(tmp.Name())
	if err != nil {
		return fmt.Errorf("gtfs: downloaded file is not a valid zip: %w", err)
	}
	zr.Close()

	finalPath := filepath.Join(GTFS_DATA_DIR, gtfs.GetStaticFileName())
	if err := os.Rename(tmp.Name(), finalPath); err != nil {
		return fmt.Errorf("gtfs: failed to replace zip: %w", err)
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	gtfs.StaticHash = hash

	return nil
}

/* gets the name of the static gtfs zip file. formatted cityname-gtfs.zip
 */
func (gtfs *GTFS) GetStaticFileName() string {
	return string(gtfs.Config.City) + "-gtfs.zip"
}

/* gets the directory where application data is saved. windows: %LOCALAPPDATA%/busline/gtfs, linux /.cache/busline/gtfs
 */
func (gtfs *GTFS) GetLocalDataDir() (string, error) {
	baseDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("gtfs: failed to get base app data directory %w", err)
	}
	return filepath.Join(baseDir, "busline", "gtfs"), nil
}

/* parse the static gtfs downloaded from DownloadStatic and populates the corresponding structs on GTFS
 */
func (gtfs *GTFS) ParseStatic() error {
	baseDir, err := gtfs.GetLocalDataDir()
	if err != nil {
		return fmt.Errorf("gtfs: failed to get local app data dir %w", err)
	}

	filepath := filepath.Join(baseDir, gtfs.GetStaticFileName())

	zr, err := zip.OpenReader(filepath)

	if err != nil {
		return fmt.Errorf("gtfs: failed to read static zip file %w", err)
	}

	defer zr.Close()

	for _, f := range zr.File {
		switch f.Name {
		case string(AgencyFile):
			rows, err := gtfs.readRecordsFromFile(f)
			if err != nil {
				return fmt.Errorf("gtfs: failed to parse agency.txt %w", err)
			}
			agency, err := gtfs.parseAgencyFromRecords(f, rows)
			if err != nil {
				return fmt.Errorf("gtfs: failed to parse agency.txt %w", err)
			}
			gtfs.Agency = agency

		case string(RoutesFile):
			rows, err := gtfs.readRecordsFromFile(f)
			if err != nil {
				return fmt.Errorf("gtfs: failed to parse routes.txt %w", err)
			}
			routes, err := gtfs.parseRoutesFromRecords(f, rows)
			if err != nil {
				return fmt.Errorf("gtfs: failed to parse routes.txt %w", err)
			}
			gtfs.Routes = routes

		case string(FeedInfoFile):
			rows, err := gtfs.readRecordsFromFile(f)
			if err != nil {
				return fmt.Errorf("gtfs: failed to parse feed_info.txt %w", err)
			}
			feedInfo, err := gtfs.parseFeedInfoFromRecords(f, rows)
			if err != nil {
				return fmt.Errorf("gtfs: failed to parse feed_info.txt %w", err)
			}
			gtfs.FeedInfo = feedInfo

		case string(StopsFile):
			rows, err := gtfs.readRecordsFromFile(f)
			if err != nil {
				return fmt.Errorf("gtfs: failed to parse stops.txt %w", err)
			}
			stops, err := gtfs.parseStopsFromRecords(f, rows)
			if err != nil {
				return fmt.Errorf("gtfs: failed to parse stops.txt %w", err)
			}
			gtfs.Stops = stops

		case string(ShapesFile):
			rows, err := gtfs.readRecordsFromFile(f)
			if err != nil {
				return fmt.Errorf("gtfs: failed to parse stop_times.txt %w", err)
			}
			shapes, err := gtfs.parseShapesFromRecords(f, rows)
			if err != nil {
				return fmt.Errorf("gtfs: failed to parse stop_times.txt %w", err)
			}
			gtfs.Shapes = shapes
		case string(CalendarFile):
			rows, err := gtfs.readRecordsFromFile(f)
			if err != nil {
				return fmt.Errorf("gtfs: failed to parse calendar.txt %w", err)
			}
			calendars, err := gtfs.parseCalendarFromRecords(f, rows)
			if err != nil {
				return fmt.Errorf("gtfs: failed to parse calendar.txt %w", err)
			}
			gtfs.Calendars = calendars

		case string(TripsFile):
			rows, err := gtfs.readRecordsFromFile(f)
			if err != nil {
				fmt.Errorf("gtfs: failed to parse from trips.txt %w", err)
			}
			trips, err := gtfs.parseTripsFromRecords(f, rows)
			if err != nil {
				fmt.Errorf("gtfs: failed to parse from trips.txt %w", err)
			}
			gtfs.Trips = trips
		}

	}

	return nil
}

/* parse the agency.txt file
 */
func (gtfs *GTFS) parseAgencyFromRecords(file *zip.File, records []map[string]string) (*Agency, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("gtfs: no valid records found in %s", file.Name)
	}

	var lastAgency *Agency

	for i, rawRow := range records {
		a := &Agency{
			Id:       rawRow["agency_id"],
			Name:     strings.TrimSpace(rawRow["agency_name"]),
			Url:      strings.TrimSpace(rawRow["agency_url"]),
			Timezone: strings.TrimSpace(rawRow["agency_timezone"]),
			Language: strings.TrimSpace(rawRow["agency_lang"]),
			Phone:    strings.TrimSpace(rawRow["agency_phone"]),
			FareUrl:  strings.TrimSpace(rawRow["agency_fare_url"]),
			Email:    strings.TrimSpace(rawRow["agency_email"]),
		}

		if a.Name == "" || a.Url == "" || a.Timezone == "" {
			return nil, fmt.Errorf("gtfs: %s: record index %d missing required field (Name, Url, or Timezone). Raw row: %+v", file.Name, i, rawRow)
		}

		lastAgency = a
	}

	return lastAgency, nil
}

/* parse the feed_info.txt file
 */
func (gtfs *GTFS) parseFeedInfoFromRecords(file *zip.File, records []map[string]string) (*FeedInfo, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("gtfs: no valid records found in %s", file.Name)
	}

	var feedInfo *FeedInfo

	for _, rawRow := range records {
		startDate, err := gtfs.parseDate(rawRow["feed_start_date"])
		if err != nil {
			return nil, fmt.Errorf("gtfs: failed to parse start date%w", err)
		}

		endDate, err := gtfs.parseDate(rawRow["feed_end_date"])
		if err != nil {
			return nil, fmt.Errorf("gtfs: failed to parse end date %w", err)
		}

		a := &FeedInfo{
			PublisherName: rawRow["feed_publisher_name"],
			Url:           rawRow["feed_publisher_url"],
			Language:      rawRow["feed_lang"],
			StartDate:     startDate,
			EndDate:       endDate,
			Version:       rawRow["feed_version"],
			Email:         rawRow["feed_contact_email"],
			ContactUrl:    rawRow["feed_contact_url"],
		}
		feedInfo = a
	}

	return feedInfo, nil
}

/* parse the routes.txt file
 */
func (gtfs *GTFS) parseRoutesFromRecords(file *zip.File, records []map[string]string) ([]Route, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("gtfs: no valid records found in %s", file.Name)
	}

	routes := make([]Route, 0, len(records))

	for _, row := range records {
		routeType, err := strconv.ParseInt(row["route_type"], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("gtfs: failed to parse route type %w", err)
		}

		sortOrder, err := strconv.ParseInt(row["route_sort_order"], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("gtfs: failed to parse sort sort %w", err)
		}

		route := Route{
			Id:           row["route_id"],
			AgencyId:     row["agency_id"],
			NetworkId:    row["network_id"],
			ShortName:    row["route_short_name"],
			LongName:     row["route_long_name"],
			Description:  row["route_desc"],
			Type:         RouteType(routeType),
			Url:          row["route_url"],
			Color:        row["route_color"],
			TextColor:    row["route_text_color"],
			ServiceName:  row["route_service_name"],
			SortOrder:    sortOrder,
			BikesAllowed: row["bikes_allowed"],
		}

		routes = append(routes, route)

	}

	return routes, nil
}

/* parse the stops.txt file
 */
func (gtfs *GTFS) parseStopsFromRecords(file *zip.File, records []map[string]string) ([]Stop, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("gtfs: no valid records found in %s", file.Name)
	}

	stops := make([]Stop, 0, len(records))

	for _, row := range records {
		lat, err := strconv.ParseFloat(row["stop_lat"], 64)
		if err != nil {
			return nil, fmt.Errorf("gtfs: failed to parse coordinates %w", err)
		}

		lon, err := strconv.ParseFloat(row["stop_lon"], 64)

		if err != nil {
			return nil, fmt.Errorf("gtfs: failed to parse coordinates %w", err)
		}

		locationType, err := strconv.ParseInt(row["location_type"], 10, 64)

		if err != nil {
			return nil, fmt.Errorf("gtfs: failed to parse location type %w", err)
		}

		wheelchairBoarding, err := strconv.ParseInt(row["wheelchair_boarding"], 10, 64)

		if err != nil {
			return nil, fmt.Errorf("gtfs: failed to parse wheelchair boarding %w", err)
		}

		stop := Stop{
			Id:                 row["stop_id"],
			Code:               row["stop_code"],
			Name:               row["stop_name"],
			Description:        row["stop_desc"],
			Latitude:           lat,
			Longitude:          lon,
			ZoneId:             row["zone_id"],
			Url:                row["stop_url"],
			LocationType:       LocationType(locationType),
			ParentStation:      row["parent_station"],
			Timezone:           row["stop_timezone"],
			WheelchairBoarding: wheelchairBoarding,
			SignType:           row["sign_type"],
			Jurisdiction:       row["jurisdiction_id"],
			RelativePosition:   row["relative_position"],
			CardinalDirection:  row["cardinal_direction"],
			PrimaryStreet:      row["primary_street"],
			AddressRange:       row["address_range"],
			CrossLocation:      row["cross_location"],
		}

		stops = append(stops, stop)
	}

	return stops, nil
}

/* parse the stop_times.txt file */
func (gtfs *GTFS) parseStopTimesFromRecords(file *zip.File, records []map[string]string) ([]StopTimes, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("gtfs: no valid records found in %s", file.Name)
	}

	stopTimes := make([]StopTimes, 0, len(records))

	for _, row := range records {
		stopTime := StopTimes{
			TripId:                row["trip_id"],
			ArrivalTime:           row["arrival_time"],
			DepartureTime:         row["departure_time"],
			StopId:                row["stop_id"],
			StopSequence:          row["stop_sequence"],
			StopHeadsign:          row["stop_headsign"],
			PickupType:            row["pickup_type"],
			DropoffType:           row["drop_off_type"],
			ShapeDistanceTraveled: row["shape_dist_traveled"],
			Timepoint:             row["timepoint"],
		}

		stopTimes = append(stopTimes, stopTime)
	}

	return stopTimes, nil
}

/* parse the shapes.txt file */
func (gtfs *GTFS) parseShapesFromRecords(file *zip.File, records []map[string]string) ([]Shape, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("gtfs: no valid records found in %s", file.Name)
	}

	shapes := make([]Shape, 0, len(records))

	for _, row := range records {
		shape := Shape{
			Id:               row["shape_id"],
			PointLatitude:    row["shape_pt_lat"],
			PointLongitude:   row["shape_pt_lon"],
			Sequence:         row["shape_pt_sequence"],
			DistanceTraveled: row["shape_dist_traveled"],
		}

		shapes = append(shapes, shape)
	}

	return shapes, nil
}

/* parse the calendar.txt file */
func (gtfs *GTFS) parseCalendarFromRecords(file *zip.File, records []map[string]string) ([]Calendar, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("gtfs: no valid records found in %s", file.Name)
	}

	calendars := make([]Calendar, 0, len(records))

	for _, row := range records {
		calendar := Calendar{
			ServiceId: row["service_id"],
			Monday:    row["monday"],
			Tuesday:   row["tuesday"],
			Wednesday: row["wednesday"],
			Thursday:  row["thursday"],
			Friday:    row["friday"],
			Saturday:  row["saturday"],
			Sunday:    row["sunday"],
			StartDate: row["start_date"],
			EndDate:   row["end_date"],
		}

		calendars = append(calendars, calendar)
	}

	return calendars, nil
}

/* parse the trips.txt file */
func (gtfs *GTFS) parseTripsFromRecords(file *zip.File, records []map[string]string) ([]Trip, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("gtfs: no valid records found in %s", file.Name)
	}

	trips := make([]Trip, 0, len(records))

	for _, row := range records {
		trip := Trip{
			RouteId:              row["route_id"],
			ServiceId:            row["service_id"],
			TripId:               row["trip_id"],
			TripHeadsign:         row["trip_headsign"],
			TripShortName:        row["trip_short_name"],
			DirectionId:          row["direction_id"],
			BlockId:              row["block_id"],
			ShapeId:              row["shape_id"],
			WheelchairAccessible: row["wheelchair_accessible"],
			BikesAllowed:         row["bikes_allowed"],
		}

		trips = append(trips, trip)
	}

	return trips, nil
}

/* reads csv file and returns an array of rows
 */
func (gtfs *GTFS) readRecordsFromFile(file *zip.File) ([]map[string]string, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("gtfs: failed to open %s: %w", file.Name, err)
	}
	defer rc.Close()

	r := csv.NewReader(rc)
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("gtfs: failed to read %s header: %w", file.Name, err)
	}

	cleanedHeader := make([]string, len(header))
	for i, h := range header {
		cleanedHeader[i] = strings.TrimSpace(strings.TrimPrefix(h, "\ufeff"))
	}

	var records []map[string]string

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("gtfs: failed to read record from %s: %w", file.Name, err)
		}

		row := make(map[string]string, len(cleanedHeader))
		for i, colName := range cleanedHeader {
			if i < len(rec) {
				row[colName] = rec[i]
			} else {
				row[colName] = ""
			}
		}
		records = append(records, row)
	}

	return records, nil
}

/* streams records from a specific zip file line-by-line without loading them all into RAM
 */
func (gtfs *GTFS) streamStopTimes(file *zip.File, fn func(row map[string]string) error) error {
	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("gtfs: failed to open %s: %w", file.Name, err)
	}
	defer rc.Close()

	r := csv.NewReader(rc)
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	header, err := r.Read()
	if err != nil {
		return fmt.Errorf("gtfs: failed to read %s header: %w", file.Name, err)
	}

	cleanedHeader := make([]string, len(header))
	for i, h := range header {
		cleanedHeader[i] = strings.TrimSpace(strings.TrimPrefix(h, "\ufeff"))
	}

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("gtfs: failed to read record from %s: %w", file.Name, err)
		}

		row := make(map[string]string, len(cleanedHeader))
		for i, colName := range cleanedHeader {
			if i < len(rec) {
				row[colName] = rec[i]
			} else {
				row[colName] = ""
			}
		}

		if err := fn(row); err != nil {
			return err
		}
	}

	return nil
}

/* streams saves stop times directly to database, to save on memory
 */
func (gtfs *GTFS) StreamAndSaveStopTimes(ctx context.Context, repo *Repository) error {
	baseDir, err := gtfs.GetLocalDataDir()
	if err != nil {
		return err
	}

	zipPath := filepath.Join(baseDir, gtfs.GetStaticFileName())
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zr.Close()

	feedId, err := repo.GetFeedIdByHash(ctx, gtfs.StaticHash)

	if err != nil {
		return err
	}

	for _, f := range zr.File {
		if f.Name == string(StopTimesFile) {
			batch := make([]StopTimes, 0, 1000)

			err := gtfs.streamStopTimes(f, func(row map[string]string) error {
				batch = append(batch, StopTimes{
					TripId:                row["trip_id"],
					ArrivalTime:           row["arrival_time"],
					DepartureTime:         row["departure_time"],
					StopId:                row["stop_id"],
					StopSequence:          row["stop_sequence"],
					StopHeadsign:          row["stop_headsign"],
					PickupType:            row["pickup_type"],
					DropoffType:           row["drop_off_type"],
					ShapeDistanceTraveled: row["shape_dist_traveled"],
					Timepoint:             row["timepoint"],
				})

				if len(batch) >= 1000 {
					if err := repo.SaveStopTimesBatch(ctx, feedId, batch); err != nil {
						return err
					}
					batch = make([]StopTimes, 0, 1000)
				}
				return nil
			})

			if err != nil {
				return err
			}

			if len(batch) > 0 {
				if err := repo.SaveStopTimesBatch(ctx, feedId, batch); err != nil {
					return err
				}
				batch = make([]StopTimes, 0, 1000)
			}
		}
	}

	return nil
}

/* parses a date formatted 20261011 YYYYMMDD into go's Time struct
 */
func (gtfs *GTFS) parseDate(date string) (time.Time, error) {
	layout := "20060102"

	parsedDate, err := time.Parse(layout, date)

	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date %w", err)
	}

	return parsedDate, nil
}

/* cities constant
 */
type City string

const (
	Madison City = "madison"
)

/* location types constant
 */
type LocationType int64

const (
	LocationStop           LocationType = 0
	LocationStation        LocationType = 1
	LocationEntranceOrExit LocationType = 2
	LocationGenericNode    LocationType = 3
	LocationBoardingArea   LocationType = 4
)

/* route types constant
 */

type RouteType int64

const (
	StreetCar  RouteType = 0
	Subway     RouteType = 1
	Rail       RouteType = 2
	Bus        RouteType = 3
	Ferry      RouteType = 4
	CableTram  RouteType = 5
	AerialList RouteType = 6
	Funicular  RouteType = 7
	TrolleyBus RouteType = 11
	Monorail   RouteType = 12
)

/* static file contant
 */
type GTFSStaticFile string

const (
	AgencyFile                 GTFSStaticFile = "agency.txt"
	CalendarDatesFile          GTFSStaticFile = "calendar_dates.txt"
	CalendarFile               GTFSStaticFile = "calendar.txt"
	ExtendedDataDictionaryFile GTFSStaticFile = "extended_data_dictionary.txt"
	FareAttributesFile         GTFSStaticFile = "fare_attributes.txt"
	FareCappingFile            GTFSStaticFile = "fare_capping.txt"
	FareLegRulesFile           GTFSStaticFile = "fare_leg_rules.txt"
	FareMediaFile              GTFSStaticFile = "fare_media.txt"
	FareProductsFile           GTFSStaticFile = "fare_products.txt"
	FareRulesFile              GTFSStaticFile = "fare_rules.txt"
	FareTransferRulesFile      GTFSStaticFile = "fare_transfer_rules.txt"
	FeedInfoFile               GTFSStaticFile = "feed_info.txt"
	FrequenciesFile            GTFSStaticFile = "frequencies.txt"
	LinkedDatasetsFile         GTFSStaticFile = "linked_datasets.txt"
	RiderCategoriesFile        GTFSStaticFile = "rider_categories.txt"
	RoutesFile                 GTFSStaticFile = "routes.txt"
	ShapesFile                 GTFSStaticFile = "shapes.txt"
	StopFeaturesFile           GTFSStaticFile = "stop_features.txt"
	StopTimesFile              GTFSStaticFile = "stop_times.txt"
	StopsFile                  GTFSStaticFile = "stops.txt"
	TermsOfUseFile             GTFSStaticFile = "terms_of_use.txt"
	TransfersFile              GTFSStaticFile = "transfers.txt"
	TranslationsFile           GTFSStaticFile = "translations.txt"
	TripsFile                  GTFSStaticFile = "trips.txt"
)

/* represents an agency from agency.txt
 */
type Agency struct {
	Id         string
	Name       string
	Url        string
	Timezone   string
	Language   string
	Phone      string
	FareUrl    string
	Email      string
	BikePolicy string
}

/* represents a route from routes.txt
 */
type Route struct {
	Id           string    `json:"id"`
	RouteId      string    `json:"route_id"`
	AgencyId     string    `json:"agency_id,omitempty"`
	NetworkId    string    `json:"network_id,omitempty"`
	ShortName    string    `json:"short_name"`
	LongName     string    `json:"long_name"`
	Description  string    `json:"description,omitempty"`
	Type         RouteType `json:"type"`
	Url          string    `json:"url"`
	Color        string    `json:"color"`
	TextColor    string    `json:"text_color"`
	ServiceName  string    `json:"service_name,omitempty"`
	SortOrder    int64     `json:"sort_order,omitempty"`
	BikesAllowed string    `json:"bikes_allowed,omitempty"`
}

/* represents feed information from feed_info.txt
 */
type FeedInfo struct {
	Id            string    `json:"id"`
	PublisherName string    `json:"publisher_name"`
	Url           string    `json:"url,omitempty"`
	Language      string    `json:"language"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	ImportedAt    time.Time `json:"imported_at"`
	Version       string    `json:"version"`
	Email         string    `json:"email,omitempty"`
	ContactUrl    string    `json:"contact_url,omitempty"`
	City          City      `json:"city,omitempty"`
}

/* represents a stop from stops.txt
 */
type Stop struct {
	Id                 string       `json:"id"`
	StopId             string       `json:"stop_id"`
	Code               string       `json:"code,omitempty"`
	Name               string       `json:"name"`
	Description        string       `json:"description,omitempty"`
	Latitude           float64      `json:"latitude"`
	Longitude          float64      `json:"longitude"`
	ZoneId             string       `json:"zone_id,omitempty"`
	Url                string       `json:"url,omitempty"`
	LocationType       LocationType `json:"location_type,omitempty"`
	ParentStation      string       `json:"parent_station,omitempty"`
	Timezone           string       `json:"timezone,omitempty"`
	WheelchairBoarding int64        `json:"wheelchair_boarding,omitempty"`
	SignType           string       `json:"sign_type,omitempty"`
	Jurisdiction       string       `json:"jurisdiction,omitempty"`
	RelativePosition   string       `json:"relative_position,omitempty"`
	CardinalDirection  string       `json:"cardinal_direction,omitempty"`
	PrimaryStreet      string       `json:"primary_street,omitempty"`
	AddressRange       string       `json:"address_range,omitempty"`
	CrossLocation      string       `json:"cross_location,omitempty"`
	DistanceTo         float64      `json:"distance_to,omitempty"`
}

/* represents stop_times from stops.txt
 */
type StopTimes struct {
	TripId                string
	ArrivalTime           string
	DepartureTime         string
	StopId                string
	StopSequence          string
	StopHeadsign          string
	PickupType            string
	DropoffType           string
	ShapeDistanceTraveled string
	Timepoint             string
}

/* represents shapes from shapes.txt
 */
type Shape struct {
	Id               string
	PointLatitude    string
	PointLongitude   string
	Sequence         string
	DistanceTraveled string
}

/* represents a calendar from calendar.txt
 */
type Calendar struct {
	ServiceId string
	Monday    string
	Tuesday   string
	Wednesday string
	Thursday  string
	Friday    string
	Saturday  string
	Sunday    string
	StartDate string
	EndDate   string
}

/* represents a trip from trips.txt
 */
type Trip struct {
	RouteId              string
	ServiceId            string
	TripId               string
	TripHeadsign         string
	TripShortName        string
	DirectionId          string
	BlockId              string
	ShapeId              string
	WheelchairAccessible string
	BikesAllowed         string
}
