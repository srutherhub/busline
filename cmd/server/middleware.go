package server

import (
	"busline/internal/app"
	"busline/internal/gtfs"
	"context"
	"net/http"
)

type ctxKey int

const FeedCtxKey ctxKey = iota

type FeedContext struct {
	City   gtfs.City
	FeedId string
}

func ApiMiddleWare(app *app.App, next http.HandlerFunc) http.HandlerFunc {
	return SetJSONContentType(SetCityContext(app, next))
}

func SetJSONContentType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetCityContext(app *app.App, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		city := gtfs.City(r.PathValue("city"))

		if string(city) == "" {
			SendJSONError(w, "missing city param", http.StatusBadRequest, ErrorCode(IllegalArgument), nil)
			return
		}

		feedInfo, err := app.GtfsRepository.GetFeedInfoBy(r.Context(), city)
		if err != nil {
			SendJSONError(w, "failed to retrieve feed info", http.StatusInternalServerError, ErrorCode(ResourceNotFound), err)
			return
		}

		ctx := context.WithValue(r.Context(), FeedCtxKey, FeedContext{City: city, FeedId: feedInfo.Id})
		next(w, r.WithContext(ctx))
	}
}
