package app

import (
	"busline/internal/gtfs"
	"os"
)

type App struct {
	GtfsRepository *gtfs.Repository
	MadisonGtfs    *gtfs.GTFS
}

func NewApp(gtfsRepo *gtfs.Repository) *App {
	staticUrl := os.Getenv("MADISONMETROGTFSSTATICURL")
	realtimeTripUrl := os.Getenv("MADISONMETROGTFSREALTIMETRIP")
	config := gtfs.GTFSConfig{StaticUrl: staticUrl, City: gtfs.Madison, RealtimeTripUrl: realtimeTripUrl}
	madisonGtfs, err := gtfs.NewGTFS(config)
	if err != nil {
		panic(err)
	}

	app := &App{GtfsRepository: gtfsRepo, MadisonGtfs: madisonGtfs}
	return app
}
