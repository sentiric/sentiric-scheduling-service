// sentiric-scheduling-service/internal/config/config.go
package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	GRPCPort string
	HttpPort string
	CertPath string
	KeyPath  string
	CaPath   string
	LogLevel string
	Env      string

	// Planlama servisi bağımlılıkları (Placeholder)
	CalendarAdapter      string // Google, Outlook, etc.
	GoogleCalendarAPIKey string
}

func Load() (*Config, error) {
	godotenv.Load()

	// Harmonik Mimari Portlar (Yatay Yetenek, 173XX bloğu atandı)
	return &Config{
		GRPCPort: GetEnv("SCHEDULING_SERVICE_GRPC_PORT", "17311"),
		HttpPort: GetEnv("SCHEDULING_SERVICE_HTTP_PORT", "17310"),

		CertPath: GetEnvOrFail("SCHEDULING_SERVICE_CERT_PATH"),
		KeyPath:  GetEnvOrFail("SCHEDULING_SERVICE_KEY_PATH"),
		CaPath:   GetEnvOrFail("GRPC_TLS_CA_PATH"),
		LogLevel: GetEnv("LOG_LEVEL", "info"),
		Env:      GetEnv("ENV", "production"),

		CalendarAdapter:      GetEnv("CALENDAR_ADAPTER", "google"),
		GoogleCalendarAPIKey: GetEnv("GOOGLE_CALENDAR_API_KEY", ""),
	}, nil
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetEnvOrFail(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatal().Str("variable", key).Msg("Gerekli ortam değişkeni tanımlı değil")
	}
	return value
}
