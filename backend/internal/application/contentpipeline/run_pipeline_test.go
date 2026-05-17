package contentpipeline

import (
	"context"
	"errors"
	"testing"

	domain "github.com/fbahesna/tiktok-agent-backend/internal/domain/contentpipeline"
)

// ━━━ MOCK ADAPTERS ━━━

type mockLLM struct {
	strategyErr error
	planErr     error
}

func (m *mockLLM) GenerateStrategy(ctx context.Context, title, scriptIdea, productURL string) (*domain.ContentStrategy, error) {
	if m.strategyErr != nil {
		return nil, m.strategyErr
	}
	return &domain.ContentStrategy{
		Hook:      domain.Hook{Text: "Mock hook for " + title, Duration: domain.Duration{Seconds: 3}},
		CTA:       domain.CTA{Text: "Buy Now!", Position: "end"},
		Narration: []string{"First we open the site", "Then we explore features"},
		Flow:      []string{"intro", "demo", "cta"},
	}, nil
}

func (m *mockLLM) GenerateVideoPlan(ctx context.Context, strategy *domain.ContentStrategy, style domain.VideoStyle) (*domain.VideoPlan, error) {
	if m.planErr != nil {
		return nil, m.planErr
	}
	return &domain.VideoPlan{
		Scenes: []domain.Scene{
			{ID: 1, Goal: "Landing page", BrowserAction: "navigate", VoiceOver: "Welcome", Duration: domain.Duration{Seconds: 5}},
			{ID: 2, Goal: "Feature tour", BrowserAction: "scroll", VoiceOver: "Check this out", Duration: domain.Duration{Seconds: 10}},
		},
		MusicStyle:    "upbeat",
		TotalDuration: domain.Duration{Seconds: 20},
	}, nil
}

type mockRecorder struct {
	err error
}

func (m *mockRecorder) Record(ctx context.Context, plan *domain.VideoPlan, productURL string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return "/tmp/mock_footage.mp4", nil
}

type mockVoice struct {
	err error
}

func (m *mockVoice) Synthesize(ctx context.Context, narrationTexts []string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return "/tmp/mock_voice.mp3", nil
}

type mockSubtitle struct {
	err error
}

func (m *mockSubtitle) Generate(ctx context.Context, audioPath string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return "/tmp/mock_subtitles.srt", nil
}

type mockRenderer struct {
	err error
}

func (m *mockRenderer) Render(ctx context.Context, footagePath, audioPath, srtPath string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return "/tmp/mock_final.mp4", nil
}

type mockStorage struct {
	err error
}

func (m *mockStorage) Store(ctx context.Context, filename string, data []byte) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return "http://localhost:8080/assets/" + filename, nil
}

func (m *mockStorage) GetURL(ctx context.Context, filename string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return "http://localhost:8080/assets/" + filename, nil
}

// ━━━ PIPELINE HELPER ━━━

func newTestPipeline(repo *mockJobRepository, overrides ...func(*RunPipelineUseCase)) *RunPipelineUseCase {
	uc := NewRunPipelineUseCase(
		repo,
		&mockLLM{},
		&mockRecorder{},
		&mockVoice{},
		&mockSubtitle{},
		&mockRenderer{},
		&mockStorage{},
	)
	for _, fn := range overrides {
		fn(uc)
	}
	return uc
}

func seedDraftJob(t *testing.T, repo *mockJobRepository, jobID string) {
	t.Helper()
	style := domain.VideoStyle{Pacing: "fast", Editing: "modern_tiktok"}
	job, err := domain.NewProductionJob(jobID, "Test TikTok Video", "Tutorial idea", "https://example.com", style)
	if err != nil {
		t.Fatalf("failed to create seed job: %v", err)
	}
	_ = repo.Save(context.Background(), job)
}

// ━━━ HAPPY PATH ━━━

func TestRunPipeline_HappyPath(t *testing.T) {
	repo := newMockJobRepository()
	seedDraftJob(t, repo, "job_happy")

	uc := newTestPipeline(repo)
	output, err := uc.Execute(context.Background(), "job_happy")
	if err != nil {
		t.Fatalf("unexpected pipeline error: %v", err)
	}

	// Verify output
	if output.JobID != "job_happy" {
		t.Errorf("expected job ID job_happy, got %s", output.JobID)
	}
	if output.VideoTitle != "Test TikTok Video" {
		t.Errorf("unexpected video title: %s", output.VideoTitle)
	}
	if output.VideoURL == "" {
		t.Error("expected non-empty video URL")
	}
	if output.SubtitleURL == "" {
		t.Error("expected non-empty subtitle URL")
	}
	if output.Caption == "" {
		t.Error("expected non-empty caption")
	}
	if len(output.Hashtags) == 0 {
		t.Error("expected non-empty hashtags")
	}

	// Verify final state in repo
	job, _ := repo.FindByID(context.Background(), "job_happy")
	if job.Status != domain.StatusCompleted {
		t.Errorf("expected completed status, got %s", job.Status)
	}
	if job.Progress != 100 {
		t.Errorf("expected progress 100, got %d", job.Progress)
	}
}

// ━━━ STEP FAILURE TESTS ━━━

func TestRunPipeline_JobNotFound(t *testing.T) {
	repo := newMockJobRepository()
	uc := newTestPipeline(repo)

	_, err := uc.Execute(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent job")
	}
}

func TestRunPipeline_StrategyFails(t *testing.T) {
	repo := newMockJobRepository()
	seedDraftJob(t, repo, "job_strat_fail")

	uc := NewRunPipelineUseCase(
		repo,
		&mockLLM{strategyErr: errors.New("LLM API timeout")},
		&mockRecorder{},
		&mockVoice{},
		&mockSubtitle{},
		&mockRenderer{},
		&mockStorage{},
	)

	_, err := uc.Execute(context.Background(), "job_strat_fail")
	if err == nil {
		t.Fatal("expected error when strategy generation fails")
	}

	// Job should be marked failed
	job, _ := repo.FindByID(context.Background(), "job_strat_fail")
	if job.Status != domain.StatusFailed {
		t.Errorf("expected failed status after strategy error, got %s", job.Status)
	}
}

func TestRunPipeline_PlanFails(t *testing.T) {
	repo := newMockJobRepository()
	seedDraftJob(t, repo, "job_plan_fail")

	uc := NewRunPipelineUseCase(
		repo,
		&mockLLM{planErr: errors.New("plan generation error")},
		&mockRecorder{},
		&mockVoice{},
		&mockSubtitle{},
		&mockRenderer{},
		&mockStorage{},
	)

	_, err := uc.Execute(context.Background(), "job_plan_fail")
	if err == nil {
		t.Fatal("expected error when plan generation fails")
	}

	job, _ := repo.FindByID(context.Background(), "job_plan_fail")
	if job.Status != domain.StatusFailed {
		t.Errorf("expected failed status after plan error, got %s", job.Status)
	}
}

func TestRunPipeline_RecordingFails(t *testing.T) {
	repo := newMockJobRepository()
	seedDraftJob(t, repo, "job_rec_fail")

	uc := NewRunPipelineUseCase(
		repo,
		&mockLLM{},
		&mockRecorder{err: errors.New("browser crashed")},
		&mockVoice{},
		&mockSubtitle{},
		&mockRenderer{},
		&mockStorage{},
	)

	_, err := uc.Execute(context.Background(), "job_rec_fail")
	if err == nil {
		t.Fatal("expected error when recording fails")
	}

	job, _ := repo.FindByID(context.Background(), "job_rec_fail")
	if job.Status != domain.StatusFailed {
		t.Errorf("expected failed after recording error, got %s", job.Status)
	}
}

func TestRunPipeline_VoiceFails(t *testing.T) {
	repo := newMockJobRepository()
	seedDraftJob(t, repo, "job_voice_fail")

	uc := NewRunPipelineUseCase(
		repo,
		&mockLLM{},
		&mockRecorder{},
		&mockVoice{err: errors.New("TTS API unavailable")},
		&mockSubtitle{},
		&mockRenderer{},
		&mockStorage{},
	)

	_, err := uc.Execute(context.Background(), "job_voice_fail")
	if err == nil {
		t.Fatal("expected error when voice synthesis fails")
	}

	job, _ := repo.FindByID(context.Background(), "job_voice_fail")
	if job.Status != domain.StatusFailed {
		t.Errorf("expected failed after voice error, got %s", job.Status)
	}
}

func TestRunPipeline_SubtitleFails(t *testing.T) {
	repo := newMockJobRepository()
	seedDraftJob(t, repo, "job_sub_fail")

	uc := NewRunPipelineUseCase(
		repo,
		&mockLLM{},
		&mockRecorder{},
		&mockVoice{},
		&mockSubtitle{err: errors.New("subtitle generation error")},
		&mockRenderer{},
		&mockStorage{},
	)

	_, err := uc.Execute(context.Background(), "job_sub_fail")
	if err == nil {
		t.Fatal("expected error when subtitle generation fails")
	}

	job, _ := repo.FindByID(context.Background(), "job_sub_fail")
	if job.Status != domain.StatusFailed {
		t.Errorf("expected failed after subtitle error, got %s", job.Status)
	}
}

func TestRunPipeline_RenderFails(t *testing.T) {
	repo := newMockJobRepository()
	seedDraftJob(t, repo, "job_render_fail")

	uc := NewRunPipelineUseCase(
		repo,
		&mockLLM{},
		&mockRecorder{},
		&mockVoice{},
		&mockSubtitle{},
		&mockRenderer{err: errors.New("FFmpeg failed")},
		&mockStorage{},
	)

	_, err := uc.Execute(context.Background(), "job_render_fail")
	if err == nil {
		t.Fatal("expected error when render fails")
	}

	job, _ := repo.FindByID(context.Background(), "job_render_fail")
	if job.Status != domain.StatusFailed {
		t.Errorf("expected failed after render error, got %s", job.Status)
	}
}

// ━━━ STATE VERIFICATION ━━━

func TestRunPipeline_StateProgression(t *testing.T) {
	repo := newMockJobRepository()
	seedDraftJob(t, repo, "job_progress")

	uc := newTestPipeline(repo)
	_, err := uc.Execute(context.Background(), "job_progress")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// After full pipeline, repo.Update should have been called multiple times
	// (once per state transition: strategized, planned, recording, recorded, rendering, completed)
	if repo.updateCalls < 6 {
		t.Errorf("expected at least 6 update calls for state transitions, got %d", repo.updateCalls)
	}
}
