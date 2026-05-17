package browser

import (
	"context"
	"log"

	domain "github.com/fbahesna/tiktok-agent-backend/internal/domain/contentpipeline"
)

// MockRecorder is a development adapter that simulates browser automation without running a real browser.
type MockRecorder struct{}

func NewMockRecorder() *MockRecorder {
	return &MockRecorder{}
}

func (r *MockRecorder) Record(ctx context.Context, plan *domain.VideoPlan, productURL string) (string, error) {
	log.Printf("[MockRecorder] Simulating browser recording of %s (%d scenes)", productURL, len(plan.Scenes))
	for _, scene := range plan.Scenes {
		log.Printf("[MockRecorder]   Scene %d: %s → %s (%ds)", scene.ID, scene.Goal, scene.BrowserAction, scene.Duration.Seconds)
	}
	// In production, this would return the path to a real recorded .webm/.mp4 file
	return "./data/assets/raw_recording.mp4", nil
}
