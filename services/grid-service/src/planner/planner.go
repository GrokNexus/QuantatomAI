package planner

import (
	"context"
	"fmt"
	"time"

	"quantatomai/grid-service/domain"
)

// -----------------------------
// Public Query Types
// -----------------------------

// MemberInfo represents a resolved dimension member.
type MemberInfo struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type AxisKey struct {
	ModelID string
	ViewID  string
}

type ProjectionWindow struct {
	RowStart int `json:"rowStart"`
	RowEnd   int `json:"rowEnd"`
	ColStart int `json:"colStart"`
	ColEnd   int `json:"colEnd"`
}

type DefaultValuePolicy interface {
	DefaultFor(measureID int64) float64
}

type FixedDefaultPolicy struct {
	Values map[int64]float64
}

func (p *FixedDefaultPolicy) DefaultFor(measureID int64) float64 {
	return p.Values[measureID]
}

type CellValue struct {
	RowIndex int     `json:"rowIndex"`
	ColIndex int     `json:"colIndex"`
	Value    float64 `json:"value"`
}

type QueryResult struct {
	Rows    [][]MemberInfo `json:"rows"`    // nested row labels
	Columns [][]MemberInfo `json:"columns"` // nested column labels
	Cells   []CellValue    `json:"cells"`
}

// -----------------------------
// Internal Plan Types
// -----------------------------

// SparsedBlockRef prevents atom-block explosion.
type SparsedBlockRef struct {
	RowIDs      []int64
	ColIDs      []int64
	MeasureIDs  []int64
	ScenarioIDs []int64
	Filters     map[string][]int64
}

type QueryPlan struct {
	ID           string
	RowAxes      [][]MemberInfo
	ColAxes      [][]MemberInfo
	PageAxis     []MemberInfo
	Filters      map[string][]MemberInfo
	BlockRef     SparsedBlockRef
	HotPreferred bool

	// Advanced IO orchestration
	MetadataRequest any             // mapping.ResultRequest or similar
	AtomRequest     SparsedBlockRef
	AtomRevision    string
	AtomRevisionUnix int64
}

// -----------------------------
// Metadata Resolver Interface
// -----------------------------

type MetadataResolver interface {
	ResolveMembers(ctx context.Context, dim string, codes []string) ([]MemberInfo, error)
	ResolveMeasureIDs(ctx context.Context, measures []string) ([]int64, error)
	ResolveScenarioIDs(ctx context.Context, scenarios []string) ([]int64, error)
}

// ProjectionCache defines a pluggable cache for axis projections.
type ProjectionCache interface {
	GetAxes(key AxisKey) (rows [][]MemberInfo, rowIDs [][]int64, cols [][]MemberInfo, colIDs [][]int64, ok bool)
	SetAxes(key AxisKey, rows [][]MemberInfo, rowIDs [][]int64, cols [][]MemberInfo, colIDs [][]int64)
}

// Planner builds query plans from GridQuery.
type Planner struct {
	metadata MetadataResolver
	cache    ProjectionCache
}

func NewPlanner(metadata MetadataResolver) *Planner {
	return &Planner{metadata: metadata}
}

// -----------------------------
// BuildQueryPlan
// -----------------------------

func (p *Planner) BuildQueryPlan(ctx context.Context, q domain.GridQuery) (*QueryPlan, error) {

	// -----------------------------
	// 1. Validate Measures & Scenarios
	// -----------------------------
	measures, ok := q.Members["Measure"]
	if !ok || len(measures) == 0 {
		return nil, fmt.Errorf("no measures provided")
	}

	scenarios, ok := q.Members["Scenario"]
	if !ok || len(scenarios) == 0 {
		return nil, fmt.Errorf("no scenarios provided")
	}

	// -----------------------------
	// 2. Resolve Measures & Scenarios
	// -----------------------------
	measureIDs, err := p.metadata.ResolveMeasureIDs(ctx, measures)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve measures: %w", err)
	}

	scenarioIDs, err := p.metadata.ResolveScenarioIDs(ctx, scenarios)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve scenarios: %w", err)
	}

	// -----------------------------
	// 3. Resolve Multi-Stacked Row Axes
	// -----------------------------
	var rowAxes [][]MemberInfo
	for _, dim := range q.Dimensions.Rows {
		mems, err := p.metadata.ResolveMembers(ctx, dim, q.Members[dim])
		if err != nil {
			return nil, fmt.Errorf("failed to resolve row dimension %s: %w", dim, err)
		}
		rowAxes = append(rowAxes, mems)
	}

	// -----------------------------
	// 4. Resolve Multi-Stacked Column Axes
	// -----------------------------
	var colAxes [][]MemberInfo
	for _, dim := range q.Dimensions.Columns {
		mems, err := p.metadata.ResolveMembers(ctx, dim, q.Members[dim])
		if err != nil {
			return nil, fmt.Errorf("failed to resolve column dimension %s: %w", dim, err)
		}
		colAxes = append(colAxes, mems)
	}

	// -----------------------------
	// 5. Resolve Filters
	// -----------------------------
	filters := make(map[string][]MemberInfo)
	for dim, codes := range q.Dimensions.Filters {
		mems, err := p.metadata.ResolveMembers(ctx, dim, codes)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve filter %s: %w", dim, err)
		}
		filters[dim] = mems
	}

	// -----------------------------
	// 6. Build SparsedBlockRef (no atom explosion)
	// -----------------------------
	block := SparsedBlockRef{
		RowIDs:      flattenMemberIDs(rowAxes),
		ColIDs:      flattenMemberIDs(colAxes),
		MeasureIDs:  measureIDs,
		ScenarioIDs: scenarioIDs,
		Filters:     flattenFilterIDs(filters),
	}

	// -----------------------------
	// 7. Build QueryPlan
	// -----------------------------
	plan := &QueryPlan{
		RowAxes:      rowAxes,
		ColAxes:      colAxes,
		Filters:      filters,
		BlockRef:     block,
		HotPreferred: true,
		// In production, populate requests for orchestration
		AtomRequest:      block,
		AtomRevision:     "rev-0",
		AtomRevisionUnix: time.Now().Truncate(time.Hour).Unix(),
	}

	return plan, nil
}

// -----------------------------
// Helpers
// -----------------------------

func flattenMemberIDs(axes [][]MemberInfo) []int64 {
	var out []int64
	for _, axis := range axes {
		for _, m := range axis {
			out = append(out, m.ID)
		}
	}
	return out
}

func flattenFilterIDs(filters map[string][]MemberInfo) map[string][]int64 {
	out := make(map[string][]int64)
	for dim, mems := range filters {
		for _, m := range mems {
			out[dim] = append(out[dim], m.ID)
		}
	}
	return out
}

// -----------------------------
// Execution Interfaces (Stubs)
// -----------------------------

type AtomFetcher interface {
	FetchAtoms(ctx context.Context, block SparsedBlockRef) (map[domain.AtomKey]float64, error)
}

type ComputeEngine interface {
	PostProcess(ctx context.Context, values map[domain.AtomKey]float64) (map[domain.AtomKey]float64, error)
}
