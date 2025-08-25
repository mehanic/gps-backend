package service

import (
	"gps-backend/internal/model"
	"gps-backend/internal/storage"
)



type Service struct {
	storage storage.Storage
}

func NewService(s storage.Storage) *Service {
	return &Service{storage: s}
}

func (s *Service) SaveTrack(track model.Track) error {
	return s.storage.SaveTrack(track)
}

func (s *Service) GetTracks() ([]model.Track, error) {
	return s.storage.GetTracks()
}

func (s *Service) GetTracksByDevice(deviceID string) ([]model.Track, error) {
	return s.storage.GetTracksByDevice(deviceID)
}

