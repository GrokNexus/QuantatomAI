// File: services/grid-service/src/handlers/grid_query_handler.go
package handlers

import (
	"context"
	"net/http"

	"quantatomai/grid-service/src/projection"
	"quantatomai/grid-service/src/storage"
)

type GridQueryRequest struct {
	PlanID       string
	ViewID       string
	WindowHash   string
	AtomRevision int64
}

func (r *GridQueryRequest) CacheKey() storage.GridCacheKey {
	return storage.GridCacheKey{
		PlanID:       r.PlanID,
		ViewID:       r.ViewID,
		WindowHash:   r.WindowHash,
		AtomRevision: r.AtomRevision,
	}
}

type GridQueryHandler struct {
	cache storage.GridCache
}

func NewGridQueryHandler(cache storage.GridCache) *GridQueryHandler {
	return &GridQueryHandler{cache: cache}
}

func (h *GridQueryHandler) HandleGridQuery(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	planID := req.URL.Query().Get("planId")
	viewID := req.URL.Query().Get("viewId")
	windowHash := req.URL.Query().Get("windowHash")
	atomRev := int64(0)

	gr := projection.GetGridResult()
	defer releaseOffHeapArenas(gr)
	defer projection.PutGridResult(gr)

	err := h.cache.Read(ctx, storage.GridCacheKey{
		PlanID:       planID,
		ViewID:       viewID,
		WindowHash:   windowHash,
		AtomRevision: atomRev,
	}, gr)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func releaseOffHeapArenas(gr *projection.GridResult) {
	if gr == nil {
		return
	}
	if gr.OffHeapCells != nil && storage.GlobalArenaManager != nil {
		storage.GlobalArenaManager.Release(gr.OffHeapCells)
		gr.OffHeapCells = nil
	}
	if gr.OffHeapStringArena != nil && storage.GlobalArenaManager != nil {
		storage.GlobalArenaManager.Release(gr.OffHeapStringArena)
		gr.OffHeapStringArena = nil
	}
}
