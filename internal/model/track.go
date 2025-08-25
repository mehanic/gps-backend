package model

import "time"

// Track — структура, описывающая одно событие/координату GPS-устройства
type Track struct {
	DeviceID   string    `json:"device_id"`  // идентификатор трекера
	Latitude   float64   `json:"latitude"`   // широта
	Longitude  float64   `json:"longitude"`  // долгота
	Altitude   float64   `json:"altitude"`   // высота, м
	Speed      float64   `json:"speed"`      // скорость в км/ч
	Heading    float64   `json:"heading"`    // направление, градусы
	Satellites int       `json:"satellites"` // количество спутников
	Accuracy   float64   `json:"accuracy"`   // точность GPS, м
	Timestamp  time.Time `json:"timestamp"`  // время события
}

// NewTrack создаёт новый трек с текущим временем, если оно не задано
func NewTrack(deviceID string, lat, lon, speed float64, opts ...Option) Track {
	t := time.Now()
	tr := Track{
		DeviceID:  deviceID,
		Latitude:  lat,
		Longitude: lon,
		Speed:     speed,
		Timestamp: t,
	}
	for _, o := range opts {
		o(&tr)
	}
	return tr
}

// Option — функция для настройки трека
type Option func(*Track)

// Сеттеры для опциональных полей
func WithTimestamp(ts time.Time) Option {
	return func(tr *Track) {
		tr.Timestamp = ts
	}
}

func WithAltitude(alt float64) Option {
	return func(tr *Track) {
		tr.Altitude = alt
	}
}

func WithHeading(h float64) Option {
	return func(tr *Track) {
		tr.Heading = h
	}
}

func WithSatellites(n int) Option {
	return func(tr *Track) {
		tr.Satellites = n
	}
}

func WithAccuracy(a float64) Option {
	return func(tr *Track) {
		tr.Accuracy = a
	}
}

// CopyTrack возвращает копию трека
func CopyTrack(tr Track) Track {
	return Track{
		DeviceID:   tr.DeviceID,
		Latitude:   tr.Latitude,
		Longitude:  tr.Longitude,
		Altitude:   tr.Altitude,
		Speed:      tr.Speed,
		Heading:    tr.Heading,
		Satellites: tr.Satellites,
		Accuracy:   tr.Accuracy,
		Timestamp:  tr.Timestamp,
	}
}

// t := model.NewTrack("device123", 50.45, 30.52, 60,
//     model.WithAltitude(120),
//     model.WithHeading(90),
//     model.WithSatellites(8),
//     model.WithAccuracy(3.5),
//     model.WithTimestamp(time.Now()),
// )
