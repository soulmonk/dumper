package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

// Represents database server and credentials
type Config struct {
	// HTTP/REST gateway start parameters section
	// HTTPPort is TCP port to listen by HTTP/REST gateway
	HTTPPort string

	// Log parameters section
	// LogLevel is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
	LogLevel int
	// LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00
	LogTimeFormat              string
	PostgresqlConnectionString string
}

// TODO application.Config vs config.Get()
var config Config

func Load() *Config {
	// load .env file from given path
	// we keep it empty it will load .env from current directory
	// if no file it's okay
	_ = godotenv.Load(".env")

	logLevel, _ := strconv.Atoi(getEnv("LOG_LEVEL", "-4"))
	config = Config{
		HTTPPort:                   getEnv("HTTP_PORT", "51101"),
		LogLevel:                   logLevel,
		LogTimeFormat:              getEnv("LOG_TIME_FORMAT", "2006-01-02T15:04:05.999999999Z07:00"),
		PostgresqlConnectionString: getEnv("POSTGRESQL_CONNECTION_STRING", ""),
	}

	return Get()
}

func getEnv(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getBoolEnv(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v == "true"
}

func Get() *Config {
	return &config
}
