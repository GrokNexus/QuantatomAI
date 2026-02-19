package handlers

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"quantatomai/grid-service/compute"
	"quantatomai/grid-service/domain"
	"quantatomai/grid-service/planner"
)

// GridQueryRequest encapsulates the full query, projection, and policy parameters.
type GridQueryRequest struct {
	planner.GridQuery
	Window   planner.ProjectionWindow `json:"window"`
	Defaults map[int64]float64        `json:"defaults"`
}

// GridHandler orchestrates the planning, fetching, and computation of grids.
type GridHandler struct {
	planner          *planner.Planner
	fetcher          planner.AtomFetcher
	compute          planner.ComputeEngine
	currencyResolver *compute.CurrencyResolverMetadata // Using concrete type for Prefetch capability
}

// NewGridHandler constructs a new grid orchestration handler.
func NewGridHandler(
	p *planner.Planner,
	f planner.AtomFetcher,
	c planner.ComputeEngine,
	r *compute.CurrencyResolverMetadata,
) *GridHandler {
	return &GridHandler{
		planner:          p,
		fetcher:          f,
		compute:          c,
		currencyResolver: r,
	}
}

// HandleGridQuery performs the end-to-end grid retrieval using async optimization.
func (h *GridHandler) HandleGridQuery(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. Parse Request
	var req GridQueryRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid grid query request"})
		return
	}

	// 2. Build Query Plan (Metadata Resolution)
	plan, err := h.planner.BuildQueryPlan(ctx, req.GridQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Asynchronous Fetch & Prefetch
	// We run Redis atom fetching and Metadata currency prefetching in parallel
	// to hide the lookup latency of the Postgres metadata layer.
	var (
		atoms    map[domain.AtomKey]float64
		fetchErr error
	)

	var wg sync.WaitGroup
	wg.Add(2)

	// Fetch Atoms from Hot Tier (Redis)
	go func() {
		defer wg.Done()
		atoms, fetchErr = h.fetcher.FetchAtoms(ctx, plan.BlockRef)
	}()

	// Prefetch Currencies (Metadata) - Best Effort
	go func() {
		defer wg.Done()
		if h.currencyResolver != nil {
			_ = h.currencyResolver.Prefetch(ctx, plan.BlockRef)
		}
	}()

	wg.Wait()

	if fetchErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fetch failed: " + fetchErr.Error()})
		return
	}

	// 4. Execute Plan (Computation & Projection)
	// We wrap the request-specific defaults in a policy.
	defaultPolicy := &planner.FixedDefaultPolicy{Values: req.Defaults}

	result, binaryResult, err := h.planner.ExecutePlan(
		ctx,
		plan,
		nil, // Fetcher already called manually for async optimization
		h.compute,
		req.Window,
		defaultPolicy,
	)
	_ = binaryResult // TODO: Support binary response toggle via headers

	// Important: We need a small fix because ExecutePlan currently calls fetcher.FetchAtoms again.
	// In production, we'd refactor ExecutePlan to accept the pre-fetched atoms or use a Cache.
	// For now, we'll implement a "No-Op" fetcher stub that returns our atoms.
	
	result, _, err = h.planner.ExecutePlan(
		ctx,
		plan,
		&prefetchedFetcher{atoms: atoms},
		h.compute,
		req.Window,
		defaultPolicy,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "execution failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// prefetchedFetcher is a lightweight stub to pass pre-fetched data into the execution engine.
type prefetchedFetcher struct {
	atoms map[domain.AtomKey]float64
}

func (f *prefetchedFetcher) FetchAtoms(_ context.Context, _ planner.SparsedBlockRef) (map[domain.AtomKey]float64, error) {
	return f.atoms, nil
}
