package config

import (
	"log"
	"os"
)

type Config struct {
	ServerPort string
	DBUrl      string
	MQTTBroker string
}

func Load() *Config {
	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBUrl:      getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/gps?sslmode=disable"),
		MQTTBroker: getEnv("MQTT_BROKER", "tcp://localhost:1883"),
	}

	log.Printf("Config loaded: %+v", cfg)
	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
