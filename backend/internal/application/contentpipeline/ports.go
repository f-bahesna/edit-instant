package contentpipeline

import (
	"context"

	domain "github.com/fbahesna/tiktok-agent-backend/internal/domain/contentpipeline"
)

// ━━━ APPLICATION PORTS (driven side — implemented by infrastructure adapters) ━━━

// LLMClient generates content strategy and video plan from text prompts.
// Adapters: OpenAI GPT-4o, Claude 3.5 Sonnet, MockLLM
type LLMClient interface {
	GenerateStrategy(ctx context.Context, title, scriptIdea, productURL string) (*domain.ContentStrategy, error)
	GenerateVideoPlan(ctx context.Context, strategy *domain.ContentStrategy, style domain.VideoStyle) (*domain.VideoPlan, error)
}

// BrowserRecorder automates browser navigation and screen recording.
// Adapters: PlaywrightRecorder, BrowserlessRecorder, MockRecorder
type BrowserRecorder interface {
	Record(ctx context.Context, plan *domain.VideoPlan, productURL string) (footagePath string, err error)
}

// VoiceSynthesizer generates AI voice narration audio from text.
// Adapters: ElevenLabsVoice, OpenAITTS, MockVoice
type VoiceSynthesizer interface {
	Synthesize(ctx context.Context, narrationTexts []string) (audioPath string, err error)
}

// SubtitleGenerator creates time-synced SRT subtitle files from audio.
// Adapters: WhisperSubtitle, MockSubtitle
type SubtitleGenerator interface {
	Generate(ctx context.Context, audioPath string) (srtPath string, err error)
}

// VideoRenderer merges footage, audio, and subtitles into the final TikTok video.
// Adapters: FFmpegRenderer, MockRenderer
type VideoRenderer interface {
	Render(ctx context.Context, footagePath, audioPath, srtPath string) (finalVideoPath string, err error)
}

// FileStorage handles persisting and serving generated asset files.
// Adapters: LocalStorage, S3Storage
type FileStorage interface {
	Store(ctx context.Context, filename string, data []byte) (url string, err error)
	GetURL(ctx context.Context, filename string) (string, error)
}
