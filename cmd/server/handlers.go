package server

import (
	"busline/internal/app"
	"busline/internal/gtfs"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

func GetFeedInfoHandler(app *app.App) http.HandlerFunc {
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

		if err := json.NewEncoder(w).Encode(feedInfo); err != nil {
			SendJSONError(w, "failed to encode response to json", http.StatusInternalServerError, ErrorCode(ResourceNotFound), err)
		}
	}
}

func GetActiveRoutes(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fc, ok := r.Context().Value(FeedCtxKey).(FeedContext)
		if !ok {
			SendJSONError(w, "feed context missing", http.StatusInternalServerError, ErrorCode(InternalServerError), nil)
			return
		}

		feedId := fc.FeedId

		routes, err := app.GtfsRepository.GetActiveRoutes(r.Context(), feedId)

		if err != nil {
			SendJSONError(w, "failed to retrieve routes", http.StatusInternalServerError, ErrorCode(ResourceNotFound), err)
			return
		}

		if err := json.NewEncoder(w).Encode(routes); err != nil {
			SendJSONError(w, "failed to encode response to json", http.StatusInternalServerError, ErrorCode(ResourceNotFound), err)
		}
	}
}

func GetNearbyStops(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fc, ok := r.Context().Value(FeedCtxKey).(FeedContext)
		if !ok {
			SendJSONError(w, "feed context missing", http.StatusInternalServerError, ErrorCode(InternalServerError), nil)
			return
		}
		feedId := fc.FeedId

		lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		if err != nil {
			SendJSONError(w, fmt.Sprintf("failed to parse latitude %s", lat), http.StatusBadRequest, ErrorCode(IllegalArgument), err)
			return
		}
		lon, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
		if err != nil {
			SendJSONError(w, fmt.Sprintf("failed to parse longitude %s", lat), http.StatusBadRequest, ErrorCode(IllegalArgument), err)
			return
		}

		stops, err := app.GtfsRepository.GetNearbyStops(r.Context(), feedId, lat, lon)
		if err != nil {
			SendJSONError(w, "failed to retrieve nearby stops", http.StatusInternalServerError, ErrorCode(InternalServerError), err)
			return
		}
		if err := json.NewEncoder(w).Encode(stops); err != nil {
			SendJSONError(w, "failed to encode response to json", http.StatusInternalServerError, ErrorCode(ResourceNotFound), err)
		}

	}
}

func HealthzHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		res := HealthResponse{
			Status: "ok",
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func SendJSONError(w http.ResponseWriter, message string, statusCode int, code ErrorCode, err error) {
	w.WriteHeader(statusCode)
	resp := ErrorResponse{Code: string(code), Message: message, Details: err.Error()}
	json.NewEncoder(w).Encode(resp)
}
