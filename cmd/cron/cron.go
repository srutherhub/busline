package cron

import (
	"busline/internal/app"
	"context"
	"log/slog"
	"time"
)

func StartCronJobs(app *app.App) {
	ctx := context.Background()
	go startNightlyJobs(ctx, app, 2, 0)
	go startRealTimePolling(ctx, app, 30*time.Second)
}

func startNightlyJobs(ctx context.Context, app *app.App, hour, minute int) {
	go func() {

		getMadisonStaticGTFS(ctx, app)

		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
			if !next.After(now) {
				next = next.AddDate(0, 0, 1)
			}

			select {
			case <-time.After(time.Until(next)):
				getMadisonStaticGTFS(ctx, app)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func startRealTimePolling(ctx context.Context, app *app.App, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		pollMadisonRealTime(ctx, app)

		for {
			select {
			case <-ticker.C:
				pollMadisonRealTime(ctx, app)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func getMadisonStaticGTFS(ctx context.Context, app *app.App) {
	err := app.MadisonGtfs.DownloadStatic(ctx)
	if err != nil {
		slog.Error("failed to download madison static gtfs files", "msg", err)
		return
	}

	err = app.MadisonGtfs.ParseStatic()
	if err != nil {
		slog.Error("failed to parse madison static gtfs files", "msg", err)
		return
	}

	isNewFeed, err := app.GtfsRepository.DoesFeedExist(app.MadisonGtfs.StaticHash)

	if err != nil {
		slog.Error("failed to check if feed is new", "msg", err)
		return
	}

	if !isNewFeed {
		err := app.GtfsRepository.SaveStaticGTFS(ctx, app.MadisonGtfs)
		if err != nil {
			slog.Error("failed to save gtfs to database", "msg", err)
			return
		}

		err = app.MadisonGtfs.StreamAndSaveStopTimes(ctx, app.GtfsRepository)

		if err != nil {
			slog.Error("failed to save stop times to database", "msg", err)
			return
		}
	}

}

func pollMadisonRealTime(ctx context.Context, app *app.App) {
	tripUpdates, err := app.MadisonGtfs.FetchRealtimeTripUpdates()
	if err != nil {
		slog.Error("failed to retrieve realtime trip updates", "msg", err)
	}

	err = app.GtfsRepository.SaveRealTimeTripUpdates(ctx, tripUpdates)
	if err != nil {
		slog.Error("failed to save realtime trip updates", "msg", err)
	}
}
