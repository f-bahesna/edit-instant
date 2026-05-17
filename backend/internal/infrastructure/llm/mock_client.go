package llm

import (
	"context"

	domain "github.com/fbahesna/tiktok-agent-backend/internal/domain/contentpipeline"
)

// MockClient is a development/testing adapter that returns hardcoded content strategy and video plans.
type MockClient struct{}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (c *MockClient) GenerateStrategy(ctx context.Context, title, scriptIdea, productURL string) (*domain.ContentStrategy, error) {
	return &domain.ContentStrategy{
		Hook: domain.Hook{
			Text:     "Kebanyakan orang memakai fitur ini dengan cara yang SALAH...",
			Duration: domain.Duration{Seconds: 3},
		},
		CTA: domain.CTA{
			Text:     "Follow untuk tips teknologi lainnya!",
			Position: "end",
		},
		Narration: []string{
			"Kebanyakan orang memakai fitur ini dengan cara yang SALAH.",
			"Fitur ini menghemat waktu saya berjam-jam setiap hari!",
			"Kamu bisa melakukan ini dalam waktu kurang dari 30 detik saja!",
			"Follow untuk tips teknologi lainnya!",
		},
		Flow: []string{
			"Hook → Show problem",
			"Navigate dashboard → Demonstrate feature",
			"Show result → CTA",
		},
	}, nil
}

func (c *MockClient) GenerateVideoPlan(ctx context.Context, strategy *domain.ContentStrategy, style domain.VideoStyle) (*domain.VideoPlan, error) {
	return &domain.VideoPlan{
		Scenes: []domain.Scene{
			{
				ID: 1, Goal: "Hook attention",
				BrowserAction: "open_landing_page",
				VoiceOver:     strategy.Hook.Text,
				Subtitle:      strategy.Hook.Text,
				CameraEffect:  "zoom_in",
				Duration:      domain.Duration{Seconds: 3},
			},
			{
				ID: 2, Goal: "Show dashboard navigation",
				BrowserAction: "navigate_dashboard",
				VoiceOver:     "Fitur ini menghemat waktu saya berjam-jam setiap hari!",
				Subtitle:      "Fitur ini menghemat waktu saya berjam-jam setiap hari!",
				CameraEffect:  "pan_right",
				Duration:      domain.Duration{Seconds: 5},
			},
			{
				ID: 3, Goal: "Demonstrate key feature",
				BrowserAction: "click_feature_button",
				VoiceOver:     "Kamu bisa melakukan ini dalam waktu kurang dari 30 detik saja!",
				Subtitle:      "Kamu bisa melakukan ini dalam waktu kurang dari 30 detik saja!",
				CameraEffect:  "zoom_in_focus",
				Duration:      domain.Duration{Seconds: 7},
			},
			{
				ID: 4, Goal: "CTA and closing",
				BrowserAction: "show_result",
				VoiceOver:     strategy.CTA.Text,
				Subtitle:      strategy.CTA.Text,
				CameraEffect:  "zoom_out",
				Duration:      domain.Duration{Seconds: 5},
			},
		},
		MusicStyle:    "lo-fi energetic beat",
		EditingStyle:  style.Editing,
		TotalDuration: domain.Duration{Seconds: 20},
	}, nil
}
