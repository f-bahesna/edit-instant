package media

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// MockVoice is a development adapter that simulates TTS voice synthesis.
type MockVoice struct{}

func NewMockVoice() *MockVoice {
	return &MockVoice{}
}

func (v *MockVoice) Synthesize(ctx context.Context, narrationTexts []string) (string, error) {
	log.Printf("[MockVoice] Simulating voice synthesis for %d narration segments", len(narrationTexts))
	audioPath := "./data/assets/voice_over.mp3"
	os.MkdirAll(filepath.Dir(audioPath), 0755)
	os.WriteFile(audioPath, []byte("MOCK VOICE AUDIO BINARY"), 0644)
	return audioPath, nil
}

// MockSubtitle is a development adapter that simulates subtitle generation.
type MockSubtitle struct{}

func NewMockSubtitle() *MockSubtitle {
	return &MockSubtitle{}
}

func (s *MockSubtitle) Generate(ctx context.Context, audioPath string) (string, error) {
	log.Printf("[MockSubtitle] Simulating SRT subtitle generation from: %s", audioPath)
	srtPath := "./data/assets/subtitles.srt"
	os.MkdirAll(filepath.Dir(srtPath), 0755)

	srtContent := `1
00:00:01,000 --> 00:00:03,000
Kebanyakan orang memakai fitur ini dengan cara yang SALAH...

2
00:00:03,000 --> 00:00:06,000
Fitur ini menghemat waktu saya berjam-jam setiap hari!

3
00:00:06,000 --> 00:00:10,000
Kamu bisa melakukan ini dalam waktu kurang dari 30 detik saja!

4
00:00:15,000 --> 00:00:20,000
Follow untuk tips teknologi lainnya!`

	os.WriteFile(srtPath, []byte(srtContent), 0644)
	return srtPath, nil
}

// MockRenderer is a development adapter that simulates FFmpeg video rendering.
type MockRenderer struct{}

func NewMockRenderer() *MockRenderer {
	return &MockRenderer{}
}

func (r *MockRenderer) Render(ctx context.Context, footagePath, audioPath, srtPath string) (string, error) {
	log.Printf("[MockRenderer] Simulating video render: footage=%s, audio=%s, srt=%s", footagePath, audioPath, srtPath)
	finalPath := "./data/assets/final_video.mp4"
	os.MkdirAll(filepath.Dir(finalPath), 0755)

	// Dynamically check if ffmpeg command is available
	_, err := exec.LookPath("ffmpeg")
	if err == nil {
		log.Println("[MockRenderer] FFmpeg detected! Generating a real, playable 9:16 vertical MP4 video placeholder...")
		
		// Run FFmpeg to generate a gorgeous 5-second 1080x1920 (9:16) MP4 video
		// We use the built-in 'testsrc' filter which is completely self-contained and does
		// not require any external fonts, ensuring 100% compatibility with lightweight container OSes.
		cmd := exec.CommandContext(ctx, "ffmpeg", "-y",
			"-f", "lavfi", "-i", "testsrc=size=1080x1920:d=5",
			"-c:v", "libx264", "-pix_fmt", "yuv420p",
			finalPath,
		)
		
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Printf("[MockRenderer] FFmpeg generation failed: %v (output: %s). Falling back to mock text file.", err, string(out))
			os.WriteFile(finalPath, []byte("MOCK TIKTOK 9:16 VIDEO CONTENT — RENDERED BY GO AGENT PIPELINE"), 0644)
		} else {
			log.Printf("[MockRenderer] Successfully generated real, openable MP4 video at %s using FFmpeg!", finalPath)
		}
	} else {
		log.Println("[MockRenderer] FFmpeg not found on path. Writing fallback mock text file.")
		os.WriteFile(finalPath, []byte("MOCK TIKTOK 9:16 VIDEO CONTENT — RENDERED BY GO AGENT PIPELINE"), 0644)
	}

	return finalPath, nil
}
