package storage

import (
	"gps-backend/internal/model"
	"sync"
)

type Storage interface {
	SaveTrack(track model.Track) error
	GetTracks() ([]model.Track, error)
	GetTracksByDevice(deviceID string) ([]model.Track, error)
}

type MemoryStorage struct {
	mu     sync.RWMutex
	tracks []model.Track
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tracks: make([]model.Track, 0),
	}
}

func (m *MemoryStorage) SaveTrack(track model.Track) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tracks = append(m.tracks, track)
	return nil
}

func (m *MemoryStorage) GetTracks() ([]model.Track, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]model.Track(nil), m.tracks...), nil
}

func (m *MemoryStorage) GetTracksByDevice(deviceID string) ([]model.Track, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := []model.Track{}
	for _, t := range m.tracks {
		if t.DeviceID == deviceID {
			res = append(res, t)
		}
	}
	return res, nil
}
