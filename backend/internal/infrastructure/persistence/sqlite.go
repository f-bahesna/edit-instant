package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// NewSQLiteConnection opens or creates a local SQLite database file with WAL mode enabled.
func NewSQLiteConnection(dbPath string) (*sql.DB, error) {
	// Ensure the parent directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	log.Printf("📂 Opening SQLite database file at: %s", dbPath)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Optimize SQLite performance for concurrent read/write (WAL Mode)
	_, err = db.Exec("PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000;")
	if err != nil {
		log.Printf("⚠️ Warning: Failed to set SQLite PRAGMA journal_mode=WAL: %v", err)
	} else {
		log.Println("⚡ SQLite Write-Ahead Logging (WAL) enabled successfully!")
	}

	log.Println("✅ SQLite database initialized successfully!")
	return db, nil
}

// RunMigrations creates all required tables if they don't exist and ensures all required columns are present.
func RunMigrations(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		script_idea TEXT,
		product_url TEXT NOT NULL,
		status TEXT NOT NULL,
		progress INTEGER DEFAULT 0,
		raw_footage_path TEXT,
		voice_over_path TEXT,
		subtitle_path TEXT,
		final_video_path TEXT,
		final_video_url TEXT,
		subtitle_url TEXT,
		caption TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create jobs table: %w", err)
	}

	// Dynamic schema migration: ensure all columns exist
	requiredColumns := map[string]string{
		"raw_footage_path": "TEXT",
		"voice_over_path":  "TEXT",
		"subtitle_path":    "TEXT",
		"final_video_path": "TEXT",
		"final_video_url":  "TEXT",
		"subtitle_url":     "TEXT",
		"caption":          "TEXT",
		"updated_at":       "TIMESTAMP",
	}

	// Query current columns
	rows, err := db.Query("PRAGMA table_info(jobs)")
	if err != nil {
		return fmt.Errorf("failed to query table info for jobs: %w", err)
	}
	defer rows.Close()

	existingColumns := make(map[string]bool)
	for rows.Next() {
		var cid int
		var name, dfltVal sql.NullString
		var ctype string
		var notnull, pk int
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltVal, &pk); err != nil {
			return fmt.Errorf("failed to scan table info row: %w", err)
		}
		if name.Valid {
			existingColumns[name.String] = true
		}
	}

	for colName, colType := range requiredColumns {
		if !existingColumns[colName] {
			log.Printf("🔄 Migrating: Adding missing column '%s' (%s) to 'jobs' table...", colName, colType)
			alterQuery := fmt.Sprintf("ALTER TABLE jobs ADD COLUMN %s %s", colName, colType)
			if _, err := db.Exec(alterQuery); err != nil {
				return fmt.Errorf("failed to add column %s: %w", colName, err)
			}
			log.Printf("✅ Migration completed: column '%s' added.", colName)
		}
	}

	log.Println("📋 Database schema synchronized (jobs table verified).")
	return nil
}
