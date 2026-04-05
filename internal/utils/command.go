package utils

import "github.com/gin-gonic/gin"

// HandlerCommand defines the standard kubectl-inspired lifecycle for executing a request.
type HandlerCommand interface {
	// Complete extracts data from the context (params, queries, body) and populates the command struct.
	Complete(ctx *gin.Context) error

	// Validate enforces business rules and parameter constraints on the populated struct.
	Validate() error

	// Run executes the core logic and writes the response back to the context.
	Run() error
}

// RunCommand executes the command lifecycle and translates returned errors into HTTP responses.
func RunCommand(ctx *gin.Context, cmd HandlerCommand) {
	if err := cmd.Complete(ctx); err != nil {
		apiErr := GetApiErr(err)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	if err := cmd.Validate(); err != nil {
		apiErr := GetApiErr(err)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	if err := cmd.Run(); err != nil {
		apiErr := GetApiErr(err)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}
}
