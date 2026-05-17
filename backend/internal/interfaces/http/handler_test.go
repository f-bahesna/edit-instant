package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	app "github.com/fbahesna/tiktok-agent-backend/internal/application/contentpipeline"
	domain "github.com/fbahesna/tiktok-agent-backend/internal/domain/contentpipeline"
	"github.com/fbahesna/tiktok-agent-backend/internal/infrastructure/persistence"
)

// setupTestDeps creates a real in-memory SQLite database + use cases wired with mock adapters.
// This gives us a true integration test of handler → use case → repo, without needing
// any external services.
func setupTestDeps(t *testing.T) (*app.CreateJobUseCase, *app.RunPipelineUseCase, *HealthChecker) {
	t.Helper()

	db, err := persistence.NewSQLiteConnection(":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if err := persistence.RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	repo := persistence.NewSQLiteJobRepository(db)
	createJobUC := app.NewCreateJobUseCase(repo)
	pipelineUC := app.NewRunPipelineUseCase(
		repo,
		&testLLM{},
		&testRecorder{},
		&testVoice{},
		&testSubtitle{},
		&testRenderer{},
		&testStorage{},
	)
	healthChecker := NewHealthChecker(db)

	return createJobUC, pipelineUC, healthChecker
}

// ━━━ HANDLER TESTS ━━━

func TestHandleGenerateVideo_Success(t *testing.T) {
	createJobUC, pipelineUC, _ := setupTestDeps(t)
	handler := NewJobHandler(createJobUC, pipelineUC)

	body := `{"title":"My TikTok Video","script_idea":"Show how to sign up","product_url":"https://example.com","style_id":"modern"}`
	req := httptest.NewRequest(http.MethodPost, "/api/generate-video", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleGenerateVideo(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status 202, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}

	if resp["success"] != true {
		t.Errorf("expected success=true, got %v", resp["success"])
	}
	if resp["job_id"] == nil || resp["job_id"] == "" {
		t.Error("expected non-empty job_id in response")
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'data' object in response")
	}
	if data["video_url"] == nil || data["video_url"] == "" {
		t.Error("expected non-empty video_url in response data")
	}
	if data["subtitle_url"] == nil || data["subtitle_url"] == "" {
		t.Error("expected non-empty subtitle_url in response data")
	}
}

func TestHandleGenerateVideo_InvalidJSON(t *testing.T) {
	createJobUC, pipelineUC, _ := setupTestDeps(t)
	handler := NewJobHandler(createJobUC, pipelineUC)

	req := httptest.NewRequest(http.MethodPost, "/api/generate-video", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleGenerateVideo(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", rr.Code)
	}
}

func TestHandleGenerateVideo_MissingTitle(t *testing.T) {
	createJobUC, pipelineUC, _ := setupTestDeps(t)
	handler := NewJobHandler(createJobUC, pipelineUC)

	body := `{"title":"","script_idea":"idea","product_url":"https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/generate-video", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleGenerateVideo(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing title, got %d", rr.Code)
	}
}

func TestHandleGenerateVideo_MissingURL(t *testing.T) {
	createJobUC, pipelineUC, _ := setupTestDeps(t)
	handler := NewJobHandler(createJobUC, pipelineUC)

	body := `{"title":"TikTok","script_idea":"idea","product_url":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/generate-video", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleGenerateVideo(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing URL, got %d", rr.Code)
	}
}

func TestHealthHandler_Success(t *testing.T) {
	_, _, healthChecker := setupTestDeps(t)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()

	healthChecker.Handle(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse health response: %v", err)
	}

	if resp["status"] != "UP" {
		t.Errorf("expected status UP, got %v", resp["status"])
	}
	if resp["database"] != "Connected (SQLite)" {
		t.Errorf("expected Connected (SQLite), got %v", resp["database"])
	}
	if resp["architecture"] != "DDD / Hexagonal / SOLID" {
		t.Errorf("expected DDD architecture label, got %v", resp["architecture"])
	}
}

// ━━━ ROUTER INTEGRATION TEST ━━━

func TestRouter_HealthEndpoint(t *testing.T) {
	createJobUC, pipelineUC, healthChecker := setupTestDeps(t)

	tmpDir := t.TempDir()
	router := NewRouter(createJobUC, pipelineUC, healthChecker, tmpDir)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 from router, got %d", rr.Code)
	}
}

func TestRouter_GenerateVideoEndpoint(t *testing.T) {
	createJobUC, pipelineUC, healthChecker := setupTestDeps(t)

	tmpDir := t.TempDir()
	router := NewRouter(createJobUC, pipelineUC, healthChecker, tmpDir)

	body := `{"title":"Router Test","script_idea":"test","product_url":"https://example.com","style_id":"energy"}`
	req := httptest.NewRequest(http.MethodPost, "/api/generate-video", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected 202 from router, got %d. Body: %s", rr.Code, rr.Body.String())
	}
}

// ━━━ TEST MOCK ADAPTERS (minimal, satisfy interface contracts) ━━━

type testLLM struct{}

func (m *testLLM) GenerateStrategy(_ context.Context, title, scriptIdea, productURL string) (*domain.ContentStrategy, error) {
	return &domain.ContentStrategy{
		Hook: domain.Hook{Text: "Hook for " + title, Duration: domain.Duration{Seconds: 3}},
		CTA:  domain.CTA{Text: "CTA", Position: "end"},
	}, nil
}

func (m *testLLM) GenerateVideoPlan(_ context.Context, _ *domain.ContentStrategy, _ domain.VideoStyle) (*domain.VideoPlan, error) {
	return &domain.VideoPlan{
		Scenes:        []domain.Scene{{ID: 1, Goal: "test", VoiceOver: "hello", Duration: domain.Duration{Seconds: 10}}},
		TotalDuration: domain.Duration{Seconds: 20},
	}, nil
}

type testRecorder struct{}

func (m *testRecorder) Record(_ context.Context, _ *domain.VideoPlan, _ string) (string, error) {
	return "/tmp/test_footage.mp4", nil
}

type testVoice struct{}

func (m *testVoice) Synthesize(_ context.Context, _ []string) (string, error) {
	return "/tmp/test_voice.mp3", nil
}

type testSubtitle struct{}

func (m *testSubtitle) Generate(_ context.Context, _ string) (string, error) {
	return "/tmp/test_subs.srt", nil
}

type testRenderer struct{}

func (m *testRenderer) Render(_ context.Context, _, _, _ string) (string, error) {
	return "/tmp/test_final.mp4", nil
}

type testStorage struct{}

func (m *testStorage) Store(_ context.Context, filename string, _ []byte) (string, error) {
	return "http://localhost:8080/assets/" + filename, nil
}

func (m *testStorage) GetURL(_ context.Context, filename string) (string, error) {
	return "http://localhost:8080/assets/" + filename, nil
}
