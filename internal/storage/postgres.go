package storage

import (
	"database/sql"
	"gps-backend/internal/model"
	"time"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (p *PostgresStorage) InitSchema() error {
	_, err := p.db.Exec(`
	CREATE TABLE IF NOT EXISTS tracks (
		id SERIAL PRIMARY KEY,
		device_id TEXT NOT NULL,
		latitude DOUBLE PRECISION NOT NULL,
		longitude DOUBLE PRECISION NOT NULL,
		altitude DOUBLE PRECISION DEFAULT 0,
		speed DOUBLE PRECISION NOT NULL,
		heading DOUBLE PRECISION DEFAULT 0,
		satellites INT DEFAULT 0,
		accuracy DOUBLE PRECISION DEFAULT 0,
		timestamp TIMESTAMP NOT NULL
	)`)
	return err
}

func (p *PostgresStorage) SaveTrack(tr model.Track) error {
	_, err := p.db.Exec(`
		INSERT INTO tracks
		(device_id, latitude, longitude, altitude, speed, heading, satellites, accuracy, timestamp)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`,
		tr.DeviceID, tr.Latitude, tr.Longitude, tr.Altitude, tr.Speed, tr.Heading, tr.Satellites, tr.Accuracy, tr.Timestamp)
	return err
}

func (p *PostgresStorage) GetTracks() ([]model.Track, error) {
	rows, err := p.db.Query(`
		SELECT device_id, latitude, longitude, altitude, speed, heading, satellites, accuracy, timestamp
		FROM tracks ORDER BY timestamp DESC LIMIT 100
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []model.Track
	for rows.Next() {
		var t model.Track
		var ts time.Time
		if err := rows.Scan(&t.DeviceID, &t.Latitude, &t.Longitude, &t.Altitude,
			&t.Speed, &t.Heading, &t.Satellites, &t.Accuracy, &ts); err != nil {
			return nil, err
		}
		t.Timestamp = ts
		tracks = append(tracks, t)
	}
	return tracks, nil
}

func (p *PostgresStorage) GetTracksByDevice(deviceID string) ([]model.Track, error) {
	rows, err := p.db.Query(`
		SELECT device_id, latitude, longitude, altitude, speed, heading, satellites, accuracy, timestamp
		FROM tracks WHERE device_id = $1 ORDER BY timestamp DESC
	`, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []model.Track
	for rows.Next() {
		var t model.Track
		var ts time.Time
		if err := rows.Scan(&t.DeviceID, &t.Latitude, &t.Longitude, &t.Altitude,
			&t.Speed, &t.Heading, &t.Satellites, &t.Accuracy, &ts); err != nil {
			return nil, err
		}
		t.Timestamp = ts
		tracks = append(tracks, t)
	}
	return tracks, nil
}

func (p *PostgresStorage) GetLastTrack(deviceID string) (*model.Track, error) {
	row := p.db.QueryRow(`
		SELECT device_id, latitude, longitude, altitude, speed, heading, satellites, accuracy, timestamp
		FROM tracks WHERE device_id = $1 ORDER BY timestamp DESC LIMIT 1
	`, deviceID)

	var t model.Track
	var ts time.Time
	if err := row.Scan(&t.DeviceID, &t.Latitude, &t.Longitude, &t.Altitude,
		&t.Speed, &t.Heading, &t.Satellites, &t.Accuracy, &ts); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	t.Timestamp = ts
	return &t, nil
}
