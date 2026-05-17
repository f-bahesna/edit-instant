package persistence

import (
	"context"
	"testing"

	domain "github.com/fbahesna/tiktok-agent-backend/internal/domain/contentpipeline"
)

func TestSQLiteJobRepository_SaveAndFind(t *testing.T) {
	// Create an in-memory SQLite database for fast, isolated integration testing
	db, err := NewSQLiteConnection(":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	repo := NewSQLiteJobRepository(db)
	ctx := context.Background()

	style := domain.VideoStyle{
		Pacing:  "fast",
		Editing: "modern_tiktok",
	}

	job, err := domain.NewProductionJob("job_test_001", "TikTok Integration Test", "Drafting integration tests", "https://tiktok.com", style)
	if err != nil {
		t.Fatalf("failed to create domain job: %v", err)
	}

	// 1. Test Save
	if err := repo.Save(ctx, job); err != nil {
		t.Fatalf("failed to save job to SQLite: %v", err)
	}

	// 2. Test FindByID
	found, err := repo.FindByID(ctx, "job_test_001")
	if err != nil {
		t.Fatalf("failed to find job: %v", err)
	}

	if found.ID != job.ID {
		t.Errorf("expected job ID %s, got %s", job.ID, found.ID)
	}
	if found.Title != job.Title {
		t.Errorf("expected title %s, got %s", job.Title, found.Title)
	}
	if found.Status != domain.StatusDraft {
		t.Errorf("expected status draft, got %s", found.Status)
	}

	// 3. Test Update (transition to Strategized)
	strategy := domain.ContentStrategy{
		Hook: domain.Hook{Text: "Test hook"},
		CTA:  domain.CTA{Text: "Test CTA"},
	}
	if err := job.ApplyStrategy(strategy); err != nil {
		t.Fatalf("failed to apply strategy: %v", err)
	}

	if err := repo.Update(ctx, job); err != nil {
		t.Fatalf("failed to update job: %v", err)
	}

	// Verify update
	updated, err := repo.FindByID(ctx, "job_test_001")
	if err != nil {
		t.Fatalf("failed to find updated job: %v", err)
	}

	if updated.Status != domain.StatusStrategized {
		t.Errorf("expected status strategized, got %s", updated.Status)
	}
	if updated.Progress != 20 {
		t.Errorf("expected progress 20, got %d", updated.Progress)
	}
}
