package http

import (
	"encoding/json"
	"log"
	"net/http"

	app "github.com/fbahesna/tiktok-agent-backend/internal/application/contentpipeline"
)

// JobHandler handles HTTP requests related to video generation jobs.
type JobHandler struct {
	createJobUC *app.CreateJobUseCase
	pipelineUC  *app.RunPipelineUseCase
}

func NewJobHandler(createJobUC *app.CreateJobUseCase, pipelineUC *app.RunPipelineUseCase) *JobHandler {
	return &JobHandler{
		createJobUC: createJobUC,
		pipelineUC:  pipelineUC,
	}
}

type generateVideoRequest struct {
	Title      string `json:"title"`
	ScriptIdea string `json:"script_idea"`
	ProductURL string `json:"product_url"`
	StyleID    string `json:"style_id"`
}

// HandleGenerateVideo creates a job and runs the full ContentPipeline.
func (h *JobHandler) HandleGenerateVideo(w http.ResponseWriter, r *http.Request) {
	var req generateVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 1. Create job via use case
	createOutput, err := h.createJobUC.Execute(ctx, app.CreateJobInput{
		Title:      req.Title,
		ScriptIdea: req.ScriptIdea,
		ProductURL: req.ProductURL,
		StyleID:    req.StyleID,
	})
	if err != nil {
		log.Printf("❌ Failed to create job: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("🎬 Job %s created, running ContentPipeline...", createOutput.JobID)

	// 2. Run the full pipeline
	pipelineOutput, err := h.pipelineUC.Execute(ctx, createOutput.JobID)
	if err != nil {
		log.Printf("❌ Pipeline failed for job %s: %v", createOutput.JobID, err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
			"job_id":  createOutput.JobID,
		})
		return
	}

	// 3. Return success with download URLs
	respondJSON(w, http.StatusAccepted, map[string]interface{}{
		"success": true,
		"message": "Video generation completed successfully!",
		"job_id":  pipelineOutput.JobID,
		"data": map[string]interface{}{
			"video_title":  pipelineOutput.VideoTitle,
			"caption":      pipelineOutput.Caption,
			"hashtags":     pipelineOutput.Hashtags,
			"duration":     pipelineOutput.Duration,
			"video_url":    pipelineOutput.VideoURL,
			"subtitle_url": pipelineOutput.SubtitleURL,
			"output_files": []string{
				"final_video.mp4",
				"thumbnail.png",
				"voice_over.mp3",
				"subtitles.srt",
			},
		},
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
