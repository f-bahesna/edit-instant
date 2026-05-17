package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// HealthChecker encapsulates the dependencies needed for the health endpoint only (ISP).
type HealthChecker struct {
	db *sql.DB
}

func NewHealthChecker(db *sql.DB) *HealthChecker {
	return &HealthChecker{db: db}
}

func (h *HealthChecker) Handle(w http.ResponseWriter, r *http.Request) {
	ffmpegVer, ffmpegStatus := checkFFmpegStatus()

	dbStatus := "Connected (SQLite)"
	if err := h.db.Ping(); err != nil {
		dbStatus = "Disconnected"
	}

	response := map[string]interface{}{
		"status":      "UP",
		"timestamp":   time.Now().Format(time.RFC3339),
		"database":    dbStatus,
		"ffmpeg":      ffmpegStatus,
		"ffmpeg_info": ffmpegVer,
		"environment": os.Getenv("ENV"),
		"db_engine":   "SQLite (pure-go, zero-CGO)",
		"architecture": "DDD / Hexagonal / SOLID",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func checkFFmpegStatus() (string, string) {
	cmd := exec.Command("ffmpeg", "-version")
	out, err := cmd.Output()
	if err != nil {
		return err.Error(), "Not Available"
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		return lines[0], "Available"
	}
	return "unknown version", "Available"
}
