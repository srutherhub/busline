package app

import "busline/internal/gtfs"

type App struct {
	GtfsRepository *gtfs.Repository
}

func NewApp(gtfsRepo *gtfs.Repository) *App {
	app := &App{GtfsRepository: gtfsRepo}
	return app
}
