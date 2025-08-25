package main

import (
	"log"
	"net/http"

	"gps-backend/internal/api"
	"gps-backend/internal/config"
	"gps-backend/internal/model"
	"gps-backend/internal/mqtt"
	"gps-backend/internal/service"
	"gps-backend/internal/storage"
)

func main() {
	cfg := config.Load()

	// --- Storage ---
	var store storage.Storage
	if cfg.DBUrl == "memory" {
		store = storage.NewMemoryStorage()
	} else {
		pg, err := storage.NewPostgresStorage(cfg.DBUrl)
		if err != nil {
			log.Fatal("cannot connect postgres:", err)
		}
		if err := pg.InitSchema(); err != nil {
			log.Fatal("cannot init schema:", err)
		}
		store = pg
	}
	svc := service.NewService(store)

	// --- API endpoints ---
	http.HandleFunc("/health", api.HealthHandler)
	http.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "GPS backend is running"}`))
	})

	// --- Tracks API ---
	http.HandleFunc("/api/tracks", func(w http.ResponseWriter, r *http.Request) {
		tracks, err := svc.GetTracks()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		api.WriteJSON(w, tracks)
	})

	// --- Web UI (static files from ./web) ---
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	// --- MQTT ---
	mqttHandler := mqtt.NewHandler(cfg.MQTTBroker)
	mqttHandler.OnTrackReceived = func(track model.Track) {
		if err := svc.SaveTrack(track); err != nil {
			log.Println("Error saving track:", err)
		} else {
			log.Println("Track saved:", track)
		}
	}

	if err := mqttHandler.Connect(); err != nil {
		log.Fatal(err)
	}
	mqttHandler.Subscribe("gps/+/track")

	log.Printf("Server starting on port %s...", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/tracks", api.TracksHandler(svc))
	http.HandleFunc("/api/tracks/by-device", api.TracksByDeviceHandler(svc))
	http.HandleFunc("/api/last-track", api.LastTrackHandler(svc))

}

// mosquitto_pub -h localhost -p 1883 -t "gps/device123/track" -m '{
//   "device_id": "device123",
//   "latitude": 50.4510,
//   "longitude": 30.5230,
//   "speed": 70.2,
//   "timestamp": "2025-08-24T18:40:00Z"
// }'
