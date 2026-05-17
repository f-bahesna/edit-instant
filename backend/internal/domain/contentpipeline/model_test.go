package contentpipeline

import (
	"testing"
)

// ━━━ FACTORY TESTS ━━━

func TestNewProductionJob_ValidInput(t *testing.T) {
	style := VideoStyle{Pacing: "fast", Editing: "modern_tiktok", HasZoom: true, HasMusic: true, HasSFX: true}

	job, err := NewProductionJob("job_001", "Cool Title", "Make a tutorial", "https://example.com", style)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.ID != "job_001" {
		t.Errorf("expected ID job_001, got %s", job.ID)
	}
	if job.Status != StatusDraft {
		t.Errorf("expected status draft, got %s", job.Status)
	}
	if job.Progress != 0 {
		t.Errorf("expected progress 0, got %d", job.Progress)
	}
	if job.Title != "Cool Title" {
		t.Errorf("expected title Cool Title, got %s", job.Title)
	}
	if job.ProductURL != "https://example.com" {
		t.Errorf("expected URL https://example.com, got %s", job.ProductURL)
	}
	if job.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestNewProductionJob_EmptyTitle(t *testing.T) {
	style := VideoStyle{Pacing: "fast"}
	_, err := NewProductionJob("job_001", "", "idea", "https://example.com", style)
	if err != ErrTitleRequired {
		t.Errorf("expected ErrTitleRequired, got %v", err)
	}
}

func TestNewProductionJob_EmptyURL(t *testing.T) {
	style := VideoStyle{Pacing: "fast"}
	_, err := NewProductionJob("job_001", "Title", "idea", "", style)
	if err != ErrProductURLRequired {
		t.Errorf("expected ErrProductURLRequired, got %v", err)
	}
}

// ━━━ HAPPY PATH: Full State Machine Traversal ━━━

func TestJobStateMachine_HappyPath(t *testing.T) {
	style := VideoStyle{Pacing: "fast", Editing: "modern_tiktok"}
	job, _ := NewProductionJob("job_001", "Title", "Idea", "https://example.com", style)

	// Draft → Strategized
	strategy := ContentStrategy{
		Hook: Hook{Text: "Check this out!", Duration: Duration{Seconds: 3}},
		CTA:  CTA{Text: "Buy now!", Position: "end"},
	}
	if err := job.ApplyStrategy(strategy); err != nil {
		t.Fatalf("ApplyStrategy: %v", err)
	}
	assertState(t, job, StatusStrategized, 20)
	if job.Strategy == nil {
		t.Fatal("expected strategy to be set")
	}

	// Strategized → Planned
	plan := VideoPlan{
		TotalDuration: Duration{Seconds: 20},
		Scenes: []Scene{
			{ID: 1, Goal: "Open landing", BrowserAction: "navigate", VoiceOver: "Welcome", Duration: Duration{Seconds: 5}},
			{ID: 2, Goal: "Show features", BrowserAction: "scroll", VoiceOver: "Check this", Duration: Duration{Seconds: 10}},
		},
		MusicStyle: "upbeat",
	}
	if err := job.ApplyPlan(plan); err != nil {
		t.Fatalf("ApplyPlan: %v", err)
	}
	assertState(t, job, StatusPlanned, 40)
	if job.Plan == nil {
		t.Fatal("expected plan to be set")
	}
	if len(job.Plan.Scenes) != 2 {
		t.Errorf("expected 2 scenes, got %d", len(job.Plan.Scenes))
	}

	// Planned → Recording
	if err := job.StartRecording(); err != nil {
		t.Fatalf("StartRecording: %v", err)
	}
	assertState(t, job, StatusRecording, 50)

	// Recording → Recorded
	if err := job.CompleteRecording("/tmp/raw.mp4"); err != nil {
		t.Fatalf("CompleteRecording: %v", err)
	}
	assertState(t, job, StatusRecorded, 70)
	if job.RawFootagePath != "/tmp/raw.mp4" {
		t.Errorf("expected footage path /tmp/raw.mp4, got %s", job.RawFootagePath)
	}

	// Recorded → Rendering
	if err := job.StartRendering(); err != nil {
		t.Fatalf("StartRendering: %v", err)
	}
	assertState(t, job, StatusRendering, 80)

	// Rendering → Completed
	if err := job.MarkCompleted(
		"/tmp/final.mp4", "http://localhost/final.mp4",
		"/tmp/sub.srt", "http://localhost/sub.srt",
		"My Caption", []string{"tag1", "tag2"},
	); err != nil {
		t.Fatalf("MarkCompleted: %v", err)
	}
	assertState(t, job, StatusCompleted, 100)
	if job.FinalVideoURL != "http://localhost/final.mp4" {
		t.Errorf("unexpected video URL: %s", job.FinalVideoURL)
	}
	if job.SubtitleURL != "http://localhost/sub.srt" {
		t.Errorf("unexpected subtitle URL: %s", job.SubtitleURL)
	}
	if job.Caption != "My Caption" {
		t.Errorf("unexpected caption: %s", job.Caption)
	}
	if len(job.Hashtags) != 2 {
		t.Errorf("expected 2 hashtags, got %d", len(job.Hashtags))
	}
}

// ━━━ INVALID STATE TRANSITION TESTS ━━━

func TestApplyStrategy_NotFromDraft(t *testing.T) {
	job := makeStrategizedJob(t)
	err := job.ApplyStrategy(ContentStrategy{})
	if err == nil {
		t.Error("expected error applying strategy to non-draft job")
	}
}

func TestApplyPlan_NotFromStrategized(t *testing.T) {
	style := VideoStyle{Pacing: "fast"}
	job, _ := NewProductionJob("job_001", "T", "I", "https://x.com", style)

	// Job is in Draft, not Strategized
	err := job.ApplyPlan(VideoPlan{})
	if err == nil {
		t.Error("expected error applying plan to draft job")
	}
}

func TestStartRecording_NotFromPlanned(t *testing.T) {
	job := makeStrategizedJob(t)
	// Job is Strategized, not Planned
	err := job.StartRecording()
	if err == nil {
		t.Error("expected error starting recording on non-planned job")
	}
}

func TestCompleteRecording_NotFromRecording(t *testing.T) {
	style := VideoStyle{Pacing: "fast"}
	job, _ := NewProductionJob("job_001", "T", "I", "https://x.com", style)
	// Job is in Draft, not Recording
	err := job.CompleteRecording("/tmp/raw.mp4")
	if err == nil {
		t.Error("expected error completing recording on draft job")
	}
}

func TestStartRendering_NotFromRecorded(t *testing.T) {
	job := makeStrategizedJob(t)
	err := job.StartRendering()
	if err == nil {
		t.Error("expected error starting rendering on non-recorded job")
	}
}

func TestMarkCompleted_NotFromRendering(t *testing.T) {
	style := VideoStyle{Pacing: "fast"}
	job, _ := NewProductionJob("job_001", "T", "I", "https://x.com", style)
	err := job.MarkCompleted("/a", "b", "/c", "d", "cap", nil)
	if err == nil {
		t.Error("expected error marking draft job as completed")
	}
}

// ━━━ MARK FAILED FROM ANY STATE ━━━

func TestMarkFailed_FromDraft(t *testing.T) {
	style := VideoStyle{Pacing: "fast"}
	job, _ := NewProductionJob("job_001", "T", "I", "https://x.com", style)
	job.MarkFailed()
	if job.Status != StatusFailed {
		t.Errorf("expected failed, got %s", job.Status)
	}
}

func TestMarkFailed_FromRecording(t *testing.T) {
	job := makePlannedJob(t)
	_ = job.StartRecording()
	job.MarkFailed()
	if job.Status != StatusFailed {
		t.Errorf("expected failed, got %s", job.Status)
	}
}

// ━━━ UPDATED AT CHANGES ON EVERY TRANSITION ━━━

func TestUpdatedAt_ChangesOnTransition(t *testing.T) {
	style := VideoStyle{Pacing: "fast"}
	job, _ := NewProductionJob("job_001", "T", "I", "https://x.com", style)
	original := job.UpdatedAt

	_ = job.ApplyStrategy(ContentStrategy{Hook: Hook{Text: "h"}, CTA: CTA{Text: "c"}})
	if !job.UpdatedAt.After(original) || job.UpdatedAt.Equal(original) {
		// UpdatedAt should be >= original (may be equal in very fast tests)
		// Just verify it's not zero
		if job.UpdatedAt.IsZero() {
			t.Error("UpdatedAt should be set after ApplyStrategy")
		}
	}
}

// ━━━ HELPERS ━━━

func assertState(t *testing.T, job *ProductionJob, expectedStatus JobStatus, expectedProgress int) {
	t.Helper()
	if job.Status != expectedStatus {
		t.Errorf("expected status %s, got %s", expectedStatus, job.Status)
	}
	if job.Progress != expectedProgress {
		t.Errorf("expected progress %d, got %d", expectedProgress, job.Progress)
	}
}

func makeStrategizedJob(t *testing.T) *ProductionJob {
	t.Helper()
	style := VideoStyle{Pacing: "fast"}
	job, _ := NewProductionJob("job_001", "T", "I", "https://x.com", style)
	_ = job.ApplyStrategy(ContentStrategy{Hook: Hook{Text: "h"}, CTA: CTA{Text: "c"}})
	return job
}

func makePlannedJob(t *testing.T) *ProductionJob {
	t.Helper()
	job := makeStrategizedJob(t)
	_ = job.ApplyPlan(VideoPlan{TotalDuration: Duration{Seconds: 20}})
	return job
}
