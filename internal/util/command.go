package util

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

// HandlerCommand defines the standard kubectl-inspired lifecycle for executing a request.
type HandlerCommand interface {
	// Complete extracts data from the context (params, queries, body) and populates the command struct.
	Complete(w http.ResponseWriter, r *http.Request) error

	// Validate enforces business rules and parameter constraints on the populated struct.
	Validate() error

	// Run executes the core logic and writes the response back to the context.
	Run() error
}

// RunCommand executes the command lifecycle and translates returned errors into HTTP responses.
func RunCommand(w http.ResponseWriter, r *http.Request, cmd HandlerCommand) {
	if err := cmd.Complete(w, r); err != nil {
		apiErr := GetApiErr(err)
		httputil.JSON(w, apiErr.StatusCode, apiErr)
		return
	}

	if err := cmd.Validate(); err != nil {
		apiErr := GetApiErr(err)
		httputil.JSON(w, apiErr.StatusCode, apiErr)
		return
	}

	if err := cmd.Run(); err != nil {
		apiErr := GetApiErr(err)
		httputil.JSON(w, apiErr.StatusCode, apiErr)
		return
	}
}
