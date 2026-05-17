package config

import "os"

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port      string
	Env       string
	DBPath    string
	AssetsDir string
	AssetsURL string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		Env:       getEnv("ENV", "development"),
		DBPath:    getEnv("DB_PATH", "./data/tiktok_agent.db"),
		AssetsDir: getEnv("ASSETS_DIR", "./data/assets"),
		AssetsURL: getEnv("ASSETS_URL", "http://localhost:8080/assets"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
