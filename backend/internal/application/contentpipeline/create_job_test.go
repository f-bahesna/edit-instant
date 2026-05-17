package contentpipeline

import (
	"context"
	"errors"
	"testing"

	domain "github.com/fbahesna/tiktok-agent-backend/internal/domain/contentpipeline"
)

// ━━━ IN-MEMORY MOCK REPOSITORY ━━━

type mockJobRepository struct {
	jobs       map[string]*domain.ProductionJob
	saveErr    error
	findErr    error
	updateErr  error
	saveCalls  int
	updateCalls int
}

func newMockJobRepository() *mockJobRepository {
	return &mockJobRepository{jobs: make(map[string]*domain.ProductionJob)}
}

func (m *mockJobRepository) Save(ctx context.Context, job *domain.ProductionJob) error {
	m.saveCalls++
	if m.saveErr != nil {
		return m.saveErr
	}
	m.jobs[job.ID] = job
	return nil
}

func (m *mockJobRepository) FindByID(ctx context.Context, id string) (*domain.ProductionJob, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	job, ok := m.jobs[id]
	if !ok {
		return nil, domain.ErrJobNotFound
	}
	return job, nil
}

func (m *mockJobRepository) Update(ctx context.Context, job *domain.ProductionJob) error {
	m.updateCalls++
	if m.updateErr != nil {
		return m.updateErr
	}
	m.jobs[job.ID] = job
	return nil
}

// ━━━ CREATE JOB USE CASE TESTS ━━━

func TestCreateJobUseCase_Success(t *testing.T) {
	repo := newMockJobRepository()
	uc := NewCreateJobUseCase(repo)

	output, err := uc.Execute(context.Background(), CreateJobInput{
		Title:      "TikTok Tutorial",
		ScriptIdea: "Show how to use the dashboard",
		ProductURL: "https://example.com",
		StyleID:    "modern",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.JobID == "" {
		t.Error("expected non-empty JobID")
	}
	if output.Status != domain.StatusDraft {
		t.Errorf("expected status draft, got %s", output.Status)
	}
	if repo.saveCalls != 1 {
		t.Errorf("expected 1 save call, got %d", repo.saveCalls)
	}

	// Verify job was persisted
	saved, err := repo.FindByID(context.Background(), output.JobID)
	if err != nil {
		t.Fatalf("job not found in repo: %v", err)
	}
	if saved.Title != "TikTok Tutorial" {
		t.Errorf("expected title TikTok Tutorial, got %s", saved.Title)
	}
}

func TestCreateJobUseCase_EmptyTitle(t *testing.T) {
	repo := newMockJobRepository()
	uc := NewCreateJobUseCase(repo)

	_, err := uc.Execute(context.Background(), CreateJobInput{
		Title:      "",
		ScriptIdea: "idea",
		ProductURL: "https://example.com",
	})
	if err == nil {
		t.Fatal("expected error for empty title")
	}
	if repo.saveCalls != 0 {
		t.Errorf("expected 0 save calls on validation error, got %d", repo.saveCalls)
	}
}

func TestCreateJobUseCase_EmptyURL(t *testing.T) {
	repo := newMockJobRepository()
	uc := NewCreateJobUseCase(repo)

	_, err := uc.Execute(context.Background(), CreateJobInput{
		Title:      "Title",
		ScriptIdea: "idea",
		ProductURL: "",
	})
	if err == nil {
		t.Fatal("expected error for empty product URL")
	}
}

func TestCreateJobUseCase_RepoSaveFails(t *testing.T) {
	repo := newMockJobRepository()
	repo.saveErr = errors.New("db connection lost")
	uc := NewCreateJobUseCase(repo)

	_, err := uc.Execute(context.Background(), CreateJobInput{
		Title:      "Title",
		ScriptIdea: "idea",
		ProductURL: "https://example.com",
	})
	if err == nil {
		t.Fatal("expected error when repo save fails")
	}
}

func TestCreateJobUseCase_StyleResolution(t *testing.T) {
	tests := []struct {
		name    string
		styleID string
		pacing  string
		editing string
	}{
		{"modern default", "modern", "fast", "modern_tiktok"},
		{"energy", "energy", "fast", "energetic_creator"},
		{"minimalist", "minimalist", "moderate", "clean_tutorial"},
		{"unknown defaults to modern", "unknown_style", "fast", "modern_tiktok"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := newMockJobRepository()
			uc := NewCreateJobUseCase(repo)

			output, err := uc.Execute(context.Background(), CreateJobInput{
				Title:      "Title",
				ScriptIdea: "idea",
				ProductURL: "https://example.com",
				StyleID:    tc.styleID,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			saved, _ := repo.FindByID(context.Background(), output.JobID)
			if saved.Style.Pacing != tc.pacing {
				t.Errorf("expected pacing %s, got %s", tc.pacing, saved.Style.Pacing)
			}
			if saved.Style.Editing != tc.editing {
				t.Errorf("expected editing %s, got %s", tc.editing, saved.Style.Editing)
			}
		})
	}
}
