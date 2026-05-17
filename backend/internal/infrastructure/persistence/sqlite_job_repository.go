package persistence

import (
	"context"
	"database/sql"
	"fmt"

	domain "github.com/fbahesna/tiktok-agent-backend/internal/domain/contentpipeline"
)

// SQLiteJobRepository implements domain.JobRepository using SQLite.
type SQLiteJobRepository struct {
	db *sql.DB
}

func NewSQLiteJobRepository(db *sql.DB) *SQLiteJobRepository {
	return &SQLiteJobRepository{db: db}
}

func (r *SQLiteJobRepository) Save(ctx context.Context, job *domain.ProductionJob) error {
	query := `INSERT INTO jobs (id, title, script_idea, product_url, status, progress, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		job.ID, job.Title, job.ScriptIdea, job.ProductURL,
		string(job.Status), job.Progress,
		job.CreatedAt, job.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("sqlite save job: %w", err)
	}
	return nil
}

func (r *SQLiteJobRepository) FindByID(ctx context.Context, id string) (*domain.ProductionJob, error) {
	query := `SELECT id, title, script_idea, product_url, status, progress,
	                 raw_footage_path, voice_over_path, subtitle_path,
	                 final_video_path, final_video_url, subtitle_url, caption,
	                 created_at, updated_at
	          FROM jobs WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	job := &domain.ProductionJob{}
	var status string
	var rawFootage, voiceOver, subtitlePath, finalVideo, finalVideoURL, subtitleURL, caption sql.NullString

	err := row.Scan(
		&job.ID, &job.Title, &job.ScriptIdea, &job.ProductURL,
		&status, &job.Progress,
		&rawFootage, &voiceOver, &subtitlePath,
		&finalVideo, &finalVideoURL, &subtitleURL, &caption,
		&job.CreatedAt, &job.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrJobNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("sqlite find job: %w", err)
	}

	job.Status = domain.JobStatus(status)
	job.RawFootagePath = rawFootage.String
	job.VoiceOverPath = voiceOver.String
	job.SubtitlePath = subtitlePath.String
	job.FinalVideoPath = finalVideo.String
	job.FinalVideoURL = finalVideoURL.String
	job.SubtitleURL = subtitleURL.String
	job.Caption = caption.String

	return job, nil
}

func (r *SQLiteJobRepository) Update(ctx context.Context, job *domain.ProductionJob) error {
	query := `UPDATE jobs SET
	            status = ?, progress = ?,
	            raw_footage_path = ?, voice_over_path = ?, subtitle_path = ?,
	            final_video_path = ?, final_video_url = ?, subtitle_url = ?,
	            caption = ?, updated_at = ?
	          WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		string(job.Status), job.Progress,
		nullStr(job.RawFootagePath), nullStr(job.VoiceOverPath), nullStr(job.SubtitlePath),
		nullStr(job.FinalVideoPath), nullStr(job.FinalVideoURL), nullStr(job.SubtitleURL),
		nullStr(job.Caption), job.UpdatedAt,
		job.ID,
	)
	if err != nil {
		return fmt.Errorf("sqlite update job: %w", err)
	}
	return nil
}

func nullStr(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
