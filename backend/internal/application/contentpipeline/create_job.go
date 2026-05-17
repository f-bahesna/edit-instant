package contentpipeline

import (
	"context"
	"fmt"
	"time"

	domain "github.com/fbahesna/tiktok-agent-backend/internal/domain/contentpipeline"
)

// CreateJobUseCase handles creating a new ProductionJob in Draft status.
type CreateJobUseCase struct {
	repo domain.JobRepository
}

func NewCreateJobUseCase(repo domain.JobRepository) *CreateJobUseCase {
	return &CreateJobUseCase{repo: repo}
}

type CreateJobInput struct {
	Title      string
	ScriptIdea string
	ProductURL string
	StyleID    string // "modern", "energy", "minimalist"
}

type CreateJobOutput struct {
	JobID  string
	Status domain.JobStatus
}

func (uc *CreateJobUseCase) Execute(ctx context.Context, input CreateJobInput) (*CreateJobOutput, error) {
	jobID := "job_tk_916_" + time.Now().Format("20060102150405")

	style := resolveVideoStyle(input.StyleID)

	job, err := domain.NewProductionJob(jobID, input.Title, input.ScriptIdea, input.ProductURL, style)
	if err != nil {
		return nil, fmt.Errorf("failed to create production job: %w", err)
	}

	if err := uc.repo.Save(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to persist job: %w", err)
	}

	return &CreateJobOutput{
		JobID:  job.ID,
		Status: job.Status,
	}, nil
}

func resolveVideoStyle(styleID string) domain.VideoStyle {
	switch styleID {
	case "energy":
		return domain.VideoStyle{
			Pacing: "fast", Editing: "energetic_creator",
			HasZoom: true, HasMusic: true, HasSFX: true,
		}
	case "minimalist":
		return domain.VideoStyle{
			Pacing: "moderate", Editing: "clean_tutorial",
			HasZoom: false, HasMusic: true, HasSFX: false,
		}
	default: // "modern"
		return domain.VideoStyle{
			Pacing: "fast", Editing: "modern_tiktok",
			HasZoom: true, HasMusic: true, HasSFX: true,
		}
	}
}
