package compute

import (
	"context"
	"fmt"
	"sort"

	"quantatomai/grid-service/domain"
)

//
// ─────────────────────────────────────────────────────────────
//   CALENDAR PROVIDER (FISCAL-AWARE)
// ─────────────────────────────────────────────────────────────
//

// CalendarProvider abstracts calendar/fiscal logic for time intelligence operations.
type CalendarProvider interface {
	YearOf(periodID int64) int
	QuarterOf(periodID int64) int
	MonthOf(periodID int64) int

	// FiscalStartMonth allows fiscal calendars (e.g., July = 7).
	FiscalStartMonth() int
}

// TimeAggType defines how to aggregate values over a time range.
type TimeAggType int

const (
	TimeAggSum     TimeAggType = iota // Standard additive aggregation.
	TimeAggAverage                    // Mean value across the period range.
	TimeAggLast                       // The value of the final period in the range (e.g., for Headcount).
)

// TimeAggregationKind defines the time intelligence logic (YTD, Roll, etc.).
type TimeAggregationKind int

const (
	TimeAggYTD     TimeAggregationKind = iota // Year-To-Date
	TimeAggQTD                                // Quarter-To-Date
	TimeAggMTD                                // Month-To-Date
	TimeAggRolling                            // Rolling Window
)

// FillPolicy defines how to handle missing periods in the grid.
type FillPolicy int

const (
	FillNone    FillPolicy = iota // Skip periods with no data.
	FillForward                   // Forward-fill cumulative values across gaps (conceptual).
)

//
// ─────────────────────────────────────────────────────────────
//   CONFIG
// ─────────────────────────────────────────────────────────────
//

// TimeIntelligenceConfig defines rules for calculating time-based metrics.
type TimeIntelligenceConfig struct {
	// Per-measure time rules.
	Measures map[int64]TimeMeasureConfig

	// Which dimension ID in AtomKey represents the time axis.
	TimeDimensionID int64

	// Enable audit trail for calculation transparency.
	EnableAudit bool

	// How to handle gaps in time.
	FillPolicy FillPolicy
}

// TimeMeasureConfig defines the specific logic for a single measure.
type TimeMeasureConfig struct {
	Kind TimeAggregationKind

	// Rolling window length (in periods) when Kind == TimeAggRolling.
	RollingWindow int

	// Aggregation type over time (Sum, Average, Last).
	AggType TimeAggType

	// Optional: override measure ID for the output.
	OutputMeasureID int64

	// Optional: scenario filter (if zero, applies to all scenarios).
	ScenarioID int64
}

//
// ─────────────────────────────────────────────────────────────
//   AUDIT (ENHANCED WINDOW PEDIGREE)
// ─────────────────────────────────────────────────────────────
//

// TimeAudit captures the inputs and result of a time intelligence calculation.
type TimeAudit struct {
	Kind       TimeAggregationKind
	AggType    TimeAggType
	ScenarioID int64

	// For cumulative: base → current.
	BasePeriodID    int64
	CurrentPeriodID int64

	// For rolling: window size and explicit window range.
	RollingWindow int
	WindowStartID int64
	WindowEndID   int64

	// Explicit list of period IDs participating in the calculation.
	WindowPeriods []int64

	Value float64
}

//
// ─────────────────────────────────────────────────────────────
//   TRANSFORMER
// ─────────────────────────────────────────────────────────────
//

// TimeIntelligenceTransformer v4 implements high-performance time aggregation.
type TimeIntelligenceTransformer struct {
	cfg      TimeIntelligenceConfig
	calendar CalendarProvider
	audit    map[domain.AtomKey]TimeAudit
}

//
// ─────────────────────────────────────────────────────────────
//   CONSTRUCTOR
// ─────────────────────────────────────────────────────────────
//

// NewTimeIntelligenceTransformer constructs a new transformer.
func NewTimeIntelligenceTransformer(cfg TimeIntelligenceConfig, cal CalendarProvider) *TimeIntelligenceTransformer {
	t := &TimeIntelligenceTransformer{
		cfg:      cfg,
		calendar: cal,
	}
	if cfg.EnableAudit {
		t.audit = make(map[domain.AtomKey]TimeAudit)
	}
	return t
}

//
// ─────────────────────────────────────────────────────────────
//   NAME / AUDIT
// ─────────────────────────────────────────────────────────────
//

func (t *TimeIntelligenceTransformer) Name() string {
	return "time_intelligence_v4"
}

// GetAudit returns the audit trail (if enabled).
func (t *TimeIntelligenceTransformer) GetAudit() map[domain.AtomKey]TimeAudit {
	return t.audit
}

//
// ─────────────────────────────────────────────────────────────
//   APPLY
// ─────────────────────────────────────────────────────────────
//

// Apply performs time-axial transformations using O(N) running sum and sliding window logic.
func (t *TimeIntelligenceTransformer) Apply(
	ctx context.Context,
	atoms map[domain.AtomKey]float64,
) (map[domain.AtomKey]float64, error) {

	if len(t.cfg.Measures) == 0 {
		return atoms, nil
	}

	out := cloneAtomMap(atoms)
	timeDimID := t.cfg.TimeDimensionID

	// Optimization: Resolve the index of the time dimension once.
	// We'll search for it in any key's DimIDs since they all share the same structure/order.
	timeIdx := -1
	for k := range atoms {
		for i, id := range k.DimIDs {
			if id == timeDimID {
				timeIdx = i
				break
			}
		}
		break // Found it (or didn't) in the first key
	}

	type groupKey struct {
		BaseKey   domain.AtomKey
		MeasureID int64
		Scenario  int64
	}

	grouped := make(map[groupKey]map[int64]float64)

	// Phase 1: Group by non-time coordinates + measure + scenario.
	for key, value := range atoms {
		measureCfg, ok := t.cfg.Measures[key.MeasureID]
		if !ok {
			continue
		}

		scenarioID := key.ScenarioID
		if measureCfg.ScenarioID != 0 && scenarioID != measureCfg.ScenarioID {
			continue
		}

		periodID := extractTimeKey(key, timeIdx)
		if periodID == 0 {
			continue
		}

		// Strip time to create grouping key.
		baseKey := stripTimeKey(key, timeIdx)

		gk := groupKey{
			BaseKey:   baseKey,
			MeasureID: key.MeasureID,
			Scenario:  scenarioID,
		}

		bucket, exists := grouped[gk]
		if !exists {
			bucket = make(map[int64]float64)
			grouped[gk] = bucket
		}
		bucket[periodID] = value
	}

	// Phase 2: Compute time intelligence per group.
	for gk, byPeriod := range grouped {
		measureCfg := t.cfg.Measures[gk.MeasureID]

		// Sort periods ascending to ensure deterministic running sums.
		periods := make([]int64, 0, len(byPeriod))
		for p := range byPeriod {
			periods = append(periods, p)
		}
		sort.Slice(periods, func(i, j int) bool { return periods[i] < periods[j] })

		switch measureCfg.Kind {
		case TimeAggRolling:
			t.applyRolling(out, gk, periods, byPeriod, measureCfg, timeIdx)
		default:
			t.applyCumulative(out, gk, periods, byPeriod, measureCfg, timeIdx)
		}
	}

	return out, nil
}

//
// ─────────────────────────────────────────────────────────────
//   CUMULATIVE (YTD / QTD / MTD) — O(N) RUNNING SUM
// ─────────────────────────────────────────────────────────────
//

func (t *TimeIntelligenceTransformer) applyCumulative(
	out map[domain.AtomKey]float64,
	gk struct {
		BaseKey   domain.AtomKey
		MeasureID int64
		Scenario  int64
	},
	periods []int64,
	byPeriod map[int64]float64,
	cfg TimeMeasureConfig,
	timeIdx int,
) {
	if t.calendar == nil {
		return
	}

	var currentYear, currentQuarter, currentMonth int
	var accSum float64
	var count int
	var lastVal float64
	var basePeriod int64

	for i, p := range periods {
		val := byPeriod[p]
		year := t.calendar.YearOf(p)
		quarter := t.calendar.QuarterOf(p)
		month := t.calendar.MonthOf(p)

		if i == 0 {
			currentYear, currentQuarter, currentMonth = year, quarter, month
			basePeriod = p
			accSum = 0
			count = 0
			lastVal = 0
		}

		boundaryChanged := false
		switch cfg.Kind {
		case TimeAggYTD:
			if year != currentYear {
				boundaryChanged = true
			}
		case TimeAggQTD:
			if year != currentYear || quarter != currentQuarter {
				boundaryChanged = true
			}
		case TimeAggMTD:
			if year != currentYear || month != currentMonth {
				boundaryChanged = true
			}
		}

		if boundaryChanged {
			currentYear, currentQuarter, currentMonth = year, quarter, month
			basePeriod = p
			accSum = 0
			count = 0
			lastVal = 0
		}

		accSum += val
		count++
		lastVal = val

		var outVal float64
		switch cfg.AggType {
		case TimeAggSum:
			outVal = accSum
		case TimeAggAverage:
			if count > 0 {
				outVal = accSum / float64(count)
			}
		case TimeAggLast:
			outVal = lastVal
		}

		targetMeasureID := cfg.OutputMeasureID
		if targetMeasureID == 0 {
			targetMeasureID = gk.MeasureID
		}

		outKey := withTimeKey(gk.BaseKey, timeIdx, p)
		outKey.MeasureID = targetMeasureID
		outKey.EnsureCanonical()

		out[outKey] = outVal

		if t.cfg.EnableAudit {
			t.audit[outKey] = TimeAudit{
				Kind:            cfg.Kind,
				AggType:         cfg.AggType,
				ScenarioID:      gk.Scenario,
				BasePeriodID:    basePeriod,
				CurrentPeriodID: p,
				RollingWindow:   0,
				WindowStartID:   basePeriod,
				WindowEndID:     p,
				WindowPeriods:   []int64{basePeriod, p},
				Value:           outVal,
			}
		}
	}
}

//
// ─────────────────────────────────────────────────────────────
//   ROLLING WINDOWS (SUM / AVG / LAST)
// ─────────────────────────────────────────────────────────────
//

func (t *TimeIntelligenceTransformer) applyRolling(
	out map[domain.AtomKey]float64,
	gk struct {
		BaseKey   domain.AtomKey
		MeasureID int64
		Scenario  int64
	},
	periods []int64,
	byPeriod map[int64]float64,
	cfg TimeMeasureConfig,
	timeIdx int,
) {
	if cfg.RollingWindow <= 0 {
		return
	}

	window := make([]int64, 0, cfg.RollingWindow)
	var windowSum float64

	for _, p := range periods {
		val := byPeriod[p]
		window = append(window, p)
		windowSum += val

		// Constant-time removal: slide the window.
		for len(window) > cfg.RollingWindow {
			oldest := window[0]
			window = window[1:]
			windowSum -= byPeriod[oldest]
		}

		if len(window) == 0 {
			continue
		}

		var outVal float64
		switch cfg.AggType {
		case TimeAggSum:
			outVal = windowSum
		case TimeAggAverage:
			outVal = windowSum / float64(len(window))
		case TimeAggLast:
			outVal = val
		}

		targetMeasureID := cfg.OutputMeasureID
		if targetMeasureID == 0 {
			targetMeasureID = gk.MeasureID
		}

		outKey := withTimeKey(gk.BaseKey, timeIdx, p)
		outKey.MeasureID = targetMeasureID
		outKey.EnsureCanonical()

		out[outKey] = outVal

		if t.cfg.EnableAudit {
			wp := make([]int64, len(window))
			copy(wp, window)

			t.audit[outKey] = TimeAudit{
				Kind:            cfg.Kind,
				AggType:         cfg.AggType,
				ScenarioID:      gk.Scenario,
				BasePeriodID:    window[0],
				CurrentPeriodID: p,
				RollingWindow:   cfg.RollingWindow,
				WindowStartID:   window[0],
				WindowEndID:     p,
				WindowPeriods:   wp,
				Value:           outVal,
			}
		}
	}
}
