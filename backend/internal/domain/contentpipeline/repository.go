package contentpipeline

import "context"

// JobRepository is the domain port for persisting and retrieving ProductionJob aggregates.
// This interface lives in the DOMAIN layer. Implementations live in INFRASTRUCTURE.
type JobRepository interface {
	Save(ctx context.Context, job *ProductionJob) error
	FindByID(ctx context.Context, id string) (*ProductionJob, error)
	Update(ctx context.Context, job *ProductionJob) error
}
