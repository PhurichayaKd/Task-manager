package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port             string
	DBDSN            string
	JWTAccessSecret  string
	JWTRefreshSecret string
	AccessTTLMin     int
	RefreshTTLHours  int

	// เพิ่มสองฟิลด์นี้
	GinMode     string
	FrontendURL string
}

func MustLoad() Config {
	return Config{
		Port:             get("APP_PORT", "8080"),
		DBDSN:            must("DB_DSN"),
		JWTAccessSecret:  must("JWT_ACCESS_SECRET"),
		JWTRefreshSecret: must("JWT_REFRESH_SECRET"),
		AccessTTLMin:     atoi(get("JWT_ACCESS_TTL_MIN", "15")),
		RefreshTTLHours:  atoi(get("JWT_REFRESH_TTL_HR", "168")),
		GinMode:          get("GIN_MODE", "release"),
		FrontendURL:      get("FRONTEND_URL", "http://localhost:5173"),
	}
}

func get(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}
func must(key string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	panic("missing required environment variable: " + key)
}
func atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic("invalid int for " + s)
	}
	return n
}
