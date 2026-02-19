package planner

import (
    "fmt"
    "sort"
)

// axisCombo represents a single combination of members and IDs.
type axisCombo struct {
    Labels []MemberInfo
    IDs    []int64
}

// AxisGenerator lazily produces combinations for stacked axes.
type AxisGenerator interface {
    Next() (axisCombo, bool)
    Err() error
}

// recursiveAxisGenerator is a lazy, depth-first generator.
type recursiveAxisGenerator struct {
    axes         [][]MemberInfo
    indexes      []int
    current      axisCombo
    started      bool
    done         bool
    err          error
    maxStackDepth int
}

// NewAxisGenerator constructs a lazy generator with safety and sorting.
func NewAxisGenerator(axes [][]MemberInfo, maxStackDepth int, sortAxes bool) (AxisGenerator, error) {
    if len(axes) == 0 {
        return &emptyAxisGenerator{}, nil
    }
    if maxStackDepth <= 0 {
        maxStackDepth = 64 // sane default
    }
    if len(axes) > maxStackDepth {
        return nil, fmt.Errorf("axis stack depth %d exceeds max %d", len(axes), maxStackDepth)
    }

    // Optional: sort each axis for predictable header ordering.
    if sortAxes {
        for i := range axes {
            sort.Slice(axes[i], func(a, b int) bool {
                // Prefer Name, fallback to Code.
                if axes[i][a].Name == axes[i][b].Name {
                    return axes[i][a].Code < axes[i][b].Code
                }
                return axes[i][a].Name < axes[i][b].Name
            })
        }
    }

    g := &recursiveAxisGenerator{
        axes:         axes,
        indexes:      make([]int, len(axes)),
        maxStackDepth: maxStackDepth,
    }
    return g, nil
}

func (g *recursiveAxisGenerator) Next() (axisCombo, bool) {
    if g.done || g.err != nil {
        return axisCombo{}, false
    }

    // First call: initialize at the first combination.
    if !g.started {
        g.started = true
        return g.buildCurrent(), true
    }

    // Increment indexes like a mixed-radix counter.
    for level := len(g.indexes) - 1; level >= 0; level-- {
        g.indexes[level]++
        if g.indexes[level] < len(g.axes[level]) {
            // Reset deeper levels to zero.
            for j := level + 1; j < len(g.indexes); j++ {
                g.indexes[j] = 0
            }
            return g.buildCurrent(), true
        }
    }

    // Exhausted all combinations.
    g.done = true
    return axisCombo{}, false
}

func (g *recursiveAxisGenerator) buildCurrent() axisCombo {
    labels := make([]MemberInfo, 0, len(g.axes))
    ids := make([]int64, 0, len(g.axes))
    for axisIdx, memberIdx := range g.indexes {
        m := g.axes[axisIdx][memberIdx]
        labels = append(labels, m)
        ids = append(ids, m.ID)
    }
    g.current = axisCombo{
        Labels: labels,
        IDs:    ids,
    }
    return g.current
}

func (g *recursiveAxisGenerator) Err() error {
    return g.err
}

// emptyAxisGenerator is a no-op generator for empty axes.
type emptyAxisGenerator struct{}

func (e *emptyAxisGenerator) Next() (axisCombo, bool) { return axisCombo{}, false }
func (e *emptyAxisGenerator) Err() error              { return nil }

// MaterializeAxisCombos materializes all combos when needed (e.g., for compatibility or small grids).
func MaterializeAxisCombos(axes [][]MemberInfo, maxStackDepth int, sortAxes bool) ([][]MemberInfo, [][]int64, error) {
    gen, err := NewAxisGenerator(axes, maxStackDepth, sortAxes)
    if err != nil {
        return nil, nil, err
    }

    var labels [][]MemberInfo
    var ids [][]int64

    for {
        c, ok := gen.Next()
        if !ok {
            break
        }
        labels = append(labels, c.Labels)
        ids = append(ids, c.IDs)
    }

    if gen.Err() != nil {
        return nil, nil, gen.Err()
    }

    return labels, ids, nil
}
