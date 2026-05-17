package main

import (
	"log"
	"net/http"
	"os"

	// Application (use cases)
	app "github.com/fbahesna/tiktok-agent-backend/internal/application/contentpipeline"

	// Infrastructure (adapters)
	"github.com/fbahesna/tiktok-agent-backend/internal/infrastructure/browser"
	"github.com/fbahesna/tiktok-agent-backend/internal/infrastructure/llm"
	"github.com/fbahesna/tiktok-agent-backend/internal/infrastructure/media"
	"github.com/fbahesna/tiktok-agent-backend/internal/infrastructure/persistence"
	"github.com/fbahesna/tiktok-agent-backend/internal/infrastructure/storage"

	// Interfaces (HTTP)
	httpapi "github.com/fbahesna/tiktok-agent-backend/internal/interfaces/http"

	// Shared
	"github.com/fbahesna/tiktok-agent-backend/pkg/config"
)

func main() {
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("  🎬 TikTok Konten AI Agent — DDD / Hexagonal Architecture")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// ━━━ LOAD CONFIG ━━━
	cfg := config.Load()

	// ━━━ INFRASTRUCTURE LAYER (Adapters) ━━━

	// 1. Database
	db, err := persistence.NewSQLiteConnection(cfg.DBPath)
	if err != nil {
		log.Fatalf("❌ Failed to initialize SQLite: %v", err)
	}
	defer db.Close()

	if err := persistence.RunMigrations(db); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	// 2. Repository
	jobRepo := persistence.NewSQLiteJobRepository(db)

	// 3. External service adapters (swap these lines to switch from mock → real!)
	llmClient := llm.NewMockClient()         // → swap to llm.NewOpenAIClient(apiKey)
	recorder := browser.NewMockRecorder()     // → swap to browser.NewPlaywrightRecorder()
	voice := media.NewMockVoice()             // → swap to media.NewElevenLabsVoice(apiKey)
	subtitler := media.NewMockSubtitle()      // → swap to media.NewWhisperSubtitle(apiKey)
	renderer := media.NewMockRenderer()       // → swap to media.NewFFmpegRenderer()
	store := storage.NewLocalStorage(cfg.AssetsDir, cfg.AssetsURL)

	// ━━━ APPLICATION LAYER (Use Cases) ━━━
	createJobUC := app.NewCreateJobUseCase(jobRepo)
	pipelineUC := app.NewRunPipelineUseCase(jobRepo, llmClient, recorder, voice, subtitler, renderer, store)

	// ━━━ INTERFACE LAYER (HTTP) ━━━
	healthChecker := httpapi.NewHealthChecker(db)

	// Ensure assets directory exists
	os.MkdirAll(cfg.AssetsDir, 0755)

	router := httpapi.NewRouter(createJobUC, pipelineUC, healthChecker, cfg.AssetsDir)

	// ━━━ START SERVER ━━━
	log.Printf("🚀 Server running on :%s (env: %s)", cfg.Port, cfg.Env)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
