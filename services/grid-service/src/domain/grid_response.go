package domain

import (
	"quantatomai/grid-service/planner"
)

// GridQueryResponse provides a strictly typed contract for non-streaming responses.
type GridQueryResponse struct {
	Query  GridQuery            `json:"query"`
	Result *planner.QueryResult `json:"result"`
	Audits map[string]any       `json:"audits,omitempty"`
	ETag   string               `json:"etag,omitempty"`
}
