package server

import (
	"busline/internal/app"
	"busline/internal/gtfs"
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

func GetFeedInfoHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		city := r.PathValue("city")
		if city == "" {
			SendJSONError(w, "missing city param", http.StatusBadRequest, ErrorCode(IllegalArgument))
			return
		}

		feedInfo, err := app.GtfsRepository.GetFeedInfoBy(r.Context(), gtfs.City(city))
		if err != nil {
			SendJSONError(w, err.Error(), http.StatusInternalServerError, ErrorCode(ResourceNotFound))
			return
		}

		if err := json.NewEncoder(w).Encode(feedInfo); err != nil {
			SendJSONError(w, err.Error(), http.StatusInternalServerError, ErrorCode(ResourceNotFound))
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

func SendJSONError(w http.ResponseWriter, message string, statusCode int, code ErrorCode) {
	w.WriteHeader(statusCode)
	resp := ErrorResponse{Code: string(code), Message: message}
	json.NewEncoder(w).Encode(resp)
}
