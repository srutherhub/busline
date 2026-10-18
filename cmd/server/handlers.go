package server

import (
	"busline/internal/app"
	"encoding/json"
	"net/http"
)

func RoutesCityHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

type HealthResponse struct {
	Status string `json:"status"`
}

func HealthzHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := HealthResponse{
			Status: "ok",
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
