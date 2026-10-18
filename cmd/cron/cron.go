package cron

import (
	"busline/internal/app"
	"busline/internal/gtfs"
	"context"
	"log/slog"
	"os"
	"time"
)

func StartCronJobs(app *app.App) {
	ctx := context.Background()

	//CONSIDER REMOVING THIS IN PRODUCTION
	go getMadisonGTFS(ctx, app)

	startDailyJobs(ctx, app, 2, 0)
}

func startDailyJobs(ctx context.Context, app *app.App, targetHour, targetMinute int) {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), targetHour, targetMinute, 0, 0, now.Location())
			time.Sleep(time.Until(next))

			getMadisonGTFS(ctx, app)
		}
	}()
}

func getMadisonGTFS(ctx context.Context, app *app.App) {
	staticUrl := os.Getenv("MADISONMETROGTFSSTATICURL")
	config := gtfs.GTFSConfig{StaticUrl: staticUrl, City: gtfs.Madison}

	madisonGtfs, err := gtfs.NewGTFS(config)

	if err != nil {
		slog.Error("failed to initialize madison static gtfs data", "msg", err)
		return
	}
	err = madisonGtfs.DownloadStatic(ctx)
	if err != nil {
		slog.Error("failed to download madison static gtfs files", "msg", err)
		return
	}

	err = madisonGtfs.ParseStatic()
	if err != nil {
		slog.Error("failed to parse madison static gtfs files", "msg", err)
		return
	}

	isNewFeed, err := app.GtfsRepository.DoesFeedExist(madisonGtfs.StaticHash)

	if err != nil {
		slog.Error("failed to check if feed is new", "msg", err)
		return
	}

	if !isNewFeed {
		err := app.GtfsRepository.SaveStaticGTFS(ctx, madisonGtfs)
		if err != nil {
			slog.Error("failed to save gtfs to database", "msg", err)
			return
		}

		err = madisonGtfs.StreamAndSaveStopTimes(ctx, app.GtfsRepository)

		if err != nil {
			slog.Error("failed to save stop times to database", "msg", err)
			return
		}
	}

}
