package projection

import "context"

type ProjectionEngine interface {
	ProjectGrid(
		ctx context.Context,
		planID, viewID, windowHash string,
		target *GridResult,
		opts ProjectionOptions,
	) error
}
