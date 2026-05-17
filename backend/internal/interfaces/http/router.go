package http

import (
	"fmt"
	"net/http"
	"path/filepath"

	app "github.com/fbahesna/tiktok-agent-backend/internal/application/contentpipeline"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// NewRouter creates and configures the Chi router with all handlers and middleware.
func NewRouter(
	createJobUC *app.CreateJobUseCase,
	pipelineUC *app.RunPipelineUseCase,
	healthChecker *HealthChecker,
	assetsDir string,
) chi.Router {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Handlers
	jobHandler := NewJobHandler(createJobUC, pipelineUC)

	r.Get("/api/health", healthChecker.Handle)
	r.Post("/api/generate-video", jobHandler.HandleGenerateVideo)

	// Static asset file server with forced download (Content-Disposition: attachment)
	fileServer := http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsDir)))
	r.Handle("/assets/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := filepath.Base(r.URL.Path)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
		fileServer.ServeHTTP(w, r)
	}))

	return r
}
