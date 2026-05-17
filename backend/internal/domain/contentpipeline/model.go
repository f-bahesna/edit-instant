package contentpipeline

import (
	"errors"
	"time"
)

// ━━━ VALUE OBJECTS ━━━

// JobStatus represents the state machine states of a ProductionJob.
type JobStatus string

const (
	StatusDraft       JobStatus = "draft"
	StatusStrategized JobStatus = "strategized"
	StatusPlanned     JobStatus = "planned"
	StatusRecording   JobStatus = "recording"
	StatusRecorded    JobStatus = "recorded"
	StatusRendering   JobStatus = "rendering"
	StatusCompleted   JobStatus = "completed"
	StatusFailed      JobStatus = "failed"
)

// VideoStyle encapsulates the visual editing preferences for a TikTok video.
type VideoStyle struct {
	Pacing   string // "fast", "moderate"
	Editing  string // "modern_tiktok", "clean_tutorial", "energetic_creator"
	HasZoom  bool
	HasMusic bool
	HasSFX   bool
}

// Duration represents a time duration in seconds for video segments.
type Duration struct {
	Seconds int
}

// ━━━ ENTITIES ━━━

// Hook is the attention-grabbing opening of a TikTok video (first 3 seconds).
type Hook struct {
	Text     string
	Duration Duration
}

// CTA is the Call-To-Action element placed at the end or as overlay.
type CTA struct {
	Text     string
	Position string // "end", "overlay"
}

// Scene represents a single scene/shot in the video plan.
type Scene struct {
	ID            int
	Goal          string
	BrowserAction string
	VoiceOver     string
	Subtitle      string
	CameraEffect  string
	Duration      Duration
}

// ContentStrategy holds the creative strategy generated from the script idea.
type ContentStrategy struct {
	Hook      Hook
	CTA       CTA
	Narration []string
	Flow      []string
}

// VideoPlan holds the detailed scene-by-scene execution plan.
type VideoPlan struct {
	Scenes        []Scene
	MusicStyle    string
	EditingStyle  string
	TotalDuration Duration
}

// ━━━ AGGREGATE ROOT ━━━

// ProductionJob is the aggregate root of the ContentPipeline bounded context.
// It enforces the state machine invariant: Draft → Strategized → Planned → Recording → Recorded → Completed.
type ProductionJob struct {
	ID             string
	Title          string
	ScriptIdea     string
	ProductURL     string
	Status         JobStatus
	Progress       int
	Style          VideoStyle
	Strategy       *ContentStrategy
	Plan           *VideoPlan
	RawFootagePath string
	VoiceOverPath  string
	SubtitlePath   string
	FinalVideoPath string
	FinalVideoURL  string
	SubtitleURL    string
	Caption        string
	Hashtags       []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ━━━ FACTORY ━━━

// NewProductionJob creates a new ProductionJob in Draft status with validation.
func NewProductionJob(id, title, scriptIdea, productURL string, style VideoStyle) (*ProductionJob, error) {
	if title == "" {
		return nil, ErrTitleRequired
	}
	if productURL == "" {
		return nil, ErrProductURLRequired
	}
	now := time.Now()
	return &ProductionJob{
		ID:         id,
		Title:      title,
		ScriptIdea: scriptIdea,
		ProductURL: productURL,
		Status:     StatusDraft,
		Progress:   0,
		Style:      style,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// ━━━ DOMAIN BEHAVIOR (state machine guards) ━━━

// ApplyStrategy transitions Draft → Strategized.
func (j *ProductionJob) ApplyStrategy(s ContentStrategy) error {
	if j.Status != StatusDraft {
		return errors.New("can only apply strategy to a draft job")
	}
	j.Strategy = &s
	j.Status = StatusStrategized
	j.Progress = 20
	j.UpdatedAt = time.Now()
	return nil
}

// ApplyPlan transitions Strategized → Planned.
func (j *ProductionJob) ApplyPlan(p VideoPlan) error {
	if j.Status != StatusStrategized {
		return errors.New("can only apply plan after strategy is set")
	}
	j.Plan = &p
	j.Status = StatusPlanned
	j.Progress = 40
	j.UpdatedAt = time.Now()
	return nil
}

// StartRecording transitions Planned → Recording.
func (j *ProductionJob) StartRecording() error {
	if j.Status != StatusPlanned {
		return errors.New("can only start recording after plan is set")
	}
	j.Status = StatusRecording
	j.Progress = 50
	j.UpdatedAt = time.Now()
	return nil
}

// CompleteRecording transitions Recording → Recorded.
func (j *ProductionJob) CompleteRecording(footagePath string) error {
	if j.Status != StatusRecording {
		return errors.New("can only complete recording while recording")
	}
	j.RawFootagePath = footagePath
	j.Status = StatusRecorded
	j.Progress = 70
	j.UpdatedAt = time.Now()
	return nil
}

// StartRendering transitions Recorded → Rendering.
func (j *ProductionJob) StartRendering() error {
	if j.Status != StatusRecorded {
		return errors.New("can only start rendering after recording is done")
	}
	j.Status = StatusRendering
	j.Progress = 80
	j.UpdatedAt = time.Now()
	return nil
}

// MarkCompleted transitions Rendering → Completed with all final output paths.
func (j *ProductionJob) MarkCompleted(videoPath, videoURL, subtitlePath, subtitleURL, caption string, hashtags []string) error {
	if j.Status != StatusRendering {
		return errors.New("can only complete after rendering")
	}
	j.FinalVideoPath = videoPath
	j.FinalVideoURL = videoURL
	j.SubtitlePath = subtitlePath
	j.SubtitleURL = subtitleURL
	j.Caption = caption
	j.Hashtags = hashtags
	j.Status = StatusCompleted
	j.Progress = 100
	j.UpdatedAt = time.Now()
	return nil
}

// MarkFailed transitions any state → Failed.
func (j *ProductionJob) MarkFailed() {
	j.Status = StatusFailed
	j.UpdatedAt = time.Now()
}
