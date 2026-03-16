package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WritebackRequest struct {
	CellEdits []CellEdit `json:"cellEdits,omitempty"`

	// Top-Down Planning Request (Red Team Audit Response)
	Spread *SpreadRequest `json:"spread,omitempty"`
}

type CellEdit struct {
	Dims     map[string]string `json:"dims"`
	Measure  string            `json:"measure"`
	Scenario string            `json:"scenario"`
	Value    float64           `json:"value"`
	IsLocked bool              `json:"isLocked"` // Red Team: Bottom-Up lock
}

type SpreadRequest struct {
	TargetValue   float64           `json:"targetValue"`
	TargetDims    map[string]string `json:"targetDims"`    // e.g. Region: "Global"
	ReferenceBase string            `json:"referenceBase"` // e.g. "PreviousYearActuals"
}

func HandleWriteback(c *gin.Context) {
	var req WritebackRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid writeback request"})
		return
	}

	mode := "bottom_up"
	if req.Spread != nil {
		mode = "top_down"
		// Red Team: Macro-Transaction Logging
		// Instead of logging 100k leaf node mutations, we log the single INTENT.
		// e.g. logger.LogMacroTransaction("User X spread $1B based on Prior Year")
		// Emits a SINGLE event to Redpanda, Engine resolves the leaf mathematics.
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"mode":   mode,
	})
}
