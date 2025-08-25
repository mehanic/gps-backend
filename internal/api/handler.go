package api

import (
	"encoding/json"
	"net/http"

	"gps-backend/internal/service"
)



type HealthResponse struct {
	Status string `json:"status"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{Status: "ok"}
	WriteJSON(w, resp)
}

func WriteJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

// --- Tracks API ---

func TracksHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tracks, err := svc.GetTracks()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		WriteJSON(w, tracks)
	}
}

func TracksByDeviceHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deviceID := r.URL.Query().Get("device_id")
		if deviceID == "" {
			http.Error(w, "device_id is required", http.StatusBadRequest)
			return
		}

		tracks, err := svc.GetTracksByDevice(deviceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		WriteJSON(w, tracks)
	}
}

func LastTrackHandler(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deviceID := r.URL.Query().Get("device_id")
		if deviceID == "" {
			http.Error(w, "device_id is required", http.StatusBadRequest)
			return
		}

		tracks, err := svc.GetTracksByDevice(deviceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(tracks) == 0 {
			http.Error(w, "no tracks found", http.StatusNotFound)
			return
		}

		last := tracks[len(tracks)-1]
		WriteJSON(w, last)
	}
}

