package contentpipeline

import (
	"context"
	"fmt"
	"log"

	domain "github.com/fbahesna/tiktok-agent-backend/internal/domain/contentpipeline"
)

// RunPipelineUseCase orchestrates the entire ContentPipeline:
// CreateJob → GenerateStrategy (Step 2) → GenerateVideoPlan (Step 3) → ExecuteRecording (Step 4)
// → SynthesizeVoice (Step 6) → GenerateSubtitles (Step 7) → RenderVideo (Step 8-10)
//
// This is the main "conductor" use case that chains all steps together.
type RunPipelineUseCase struct {
	repo      domain.JobRepository
	llm       LLMClient
	recorder  BrowserRecorder
	voice     VoiceSynthesizer
	subtitle  SubtitleGenerator
	renderer  VideoRenderer
	storage   FileStorage
}

func NewRunPipelineUseCase(
	repo domain.JobRepository,
	llm LLMClient,
	recorder BrowserRecorder,
	voice VoiceSynthesizer,
	subtitle SubtitleGenerator,
	renderer VideoRenderer,
	storage FileStorage,
) *RunPipelineUseCase {
	return &RunPipelineUseCase{
		repo:     repo,
		llm:      llm,
		recorder: recorder,
		voice:    voice,
		subtitle: subtitle,
		renderer: renderer,
		storage:  storage,
	}
}

type PipelineOutput struct {
	JobID       string
	VideoTitle  string
	Caption     string
	Hashtags    []string
	Duration    string
	VideoURL    string
	SubtitleURL string
}

// Execute runs the full content production pipeline for a given job ID.
func (uc *RunPipelineUseCase) Execute(ctx context.Context, jobID string) (*PipelineOutput, error) {
	// 1. Fetch job from repository
	job, err := uc.repo.FindByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to find job: %w", err)
	}

	// ━━━ STEP 2: Generate Content Strategy ━━━
	log.Printf("[ContentPipeline] Step 2 — Generating content strategy for job %s", jobID)
	strategy, err := uc.llm.GenerateStrategy(ctx, job.Title, job.ScriptIdea, job.ProductURL)
	if err != nil {
		job.MarkFailed()
		uc.repo.Update(ctx, job)
		return nil, fmt.Errorf("step 2 failed (strategy): %w", err)
	}
	if err := job.ApplyStrategy(*strategy); err != nil {
		return nil, fmt.Errorf("state transition failed: %w", err)
	}
	uc.repo.Update(ctx, job)

	// ━━━ STEP 3: Generate Video Plan ━━━
	log.Printf("[ContentPipeline] Step 3 — Generating video plan for job %s", jobID)
	plan, err := uc.llm.GenerateVideoPlan(ctx, strategy, job.Style)
	if err != nil {
		job.MarkFailed()
		uc.repo.Update(ctx, job)
		return nil, fmt.Errorf("step 3 failed (plan): %w", err)
	}
	if err := job.ApplyPlan(*plan); err != nil {
		return nil, fmt.Errorf("state transition failed: %w", err)
	}
	uc.repo.Update(ctx, job)

	// ━━━ STEP 4 & 5: Browser Automation & Recording ━━━
	log.Printf("[ContentPipeline] Step 4 — Executing browser recording for job %s", jobID)
	if err := job.StartRecording(); err != nil {
		return nil, fmt.Errorf("state transition failed: %w", err)
	}
	uc.repo.Update(ctx, job)

	footagePath, err := uc.recorder.Record(ctx, job.Plan, job.ProductURL)
	if err != nil {
		job.MarkFailed()
		uc.repo.Update(ctx, job)
		return nil, fmt.Errorf("step 4 failed (recording): %w", err)
	}
	if err := job.CompleteRecording(footagePath); err != nil {
		return nil, fmt.Errorf("state transition failed: %w", err)
	}
	uc.repo.Update(ctx, job)

	// ━━━ STEP 6: Voice Over Generation ━━━
	log.Printf("[ContentPipeline] Step 6 — Synthesizing voice for job %s", jobID)
	narrationTexts := make([]string, 0)
	if job.Plan != nil {
		for _, scene := range job.Plan.Scenes {
			if scene.VoiceOver != "" {
				narrationTexts = append(narrationTexts, scene.VoiceOver)
			}
		}
	}
	audioPath, err := uc.voice.Synthesize(ctx, narrationTexts)
	if err != nil {
		job.MarkFailed()
		uc.repo.Update(ctx, job)
		return nil, fmt.Errorf("step 6 failed (voice): %w", err)
	}
	job.VoiceOverPath = audioPath

	// ━━━ STEP 7: Subtitle Generation ━━━
	log.Printf("[ContentPipeline] Step 7 — Generating subtitles for job %s", jobID)
	srtPath, err := uc.subtitle.Generate(ctx, audioPath)
	if err != nil {
		job.MarkFailed()
		uc.repo.Update(ctx, job)
		return nil, fmt.Errorf("step 7 failed (subtitle): %w", err)
	}

	// ━━━ STEP 8-10: Render Final Video ━━━
	log.Printf("[ContentPipeline] Step 8-10 — Rendering final video for job %s", jobID)
	if err := job.StartRendering(); err != nil {
		return nil, fmt.Errorf("state transition failed: %w", err)
	}
	uc.repo.Update(ctx, job)

	finalVideoPath, err := uc.renderer.Render(ctx, footagePath, audioPath, srtPath)
	if err != nil {
		job.MarkFailed()
		uc.repo.Update(ctx, job)
		return nil, fmt.Errorf("step 8-10 failed (render): %w", err)
	}

	// Get download URLs
	videoURL, _ := uc.storage.GetURL(ctx, "final_video.mp4")
	subtitleURL, _ := uc.storage.GetURL(ctx, "subtitles.srt")

	// Build caption and hashtags
	caption := fmt.Sprintf("%s 🚀 #SaaS #AI #Coding #TikTokAgent #Startup", job.Title)
	hashtags := []string{"SaaS", "AI", "Coding", "TikTokAgent", "Startup"}

	// ━━━ MARK COMPLETED ━━━
	if err := job.MarkCompleted(finalVideoPath, videoURL, srtPath, subtitleURL, caption, hashtags); err != nil {
		return nil, fmt.Errorf("state transition failed: %w", err)
	}
	uc.repo.Update(ctx, job)

	log.Printf("[ContentPipeline] ✅ Job %s completed successfully!", jobID)

	return &PipelineOutput{
		JobID:       job.ID,
		VideoTitle:  job.Title,
		Caption:     caption,
		Hashtags:    hashtags,
		Duration:    "20s",
		VideoURL:    videoURL,
		SubtitleURL: subtitleURL,
	}, nil
}
