# API Handler Pattern

This project implements a **Command Pattern** for HTTP request handlers. The structure is heavily inspired by Kubernetes (`kubectl`) source code and cleanly separates the concern of parsing requests, validating state, and executing business logic.

By using this standardized structure, we avoid massive, undocumented monolithic functions and duplicated error handling.

## The Lifecycle

Instead of writing everything inside an `api.MyHandler(ctx *gin.Context)` function, every handler must define and execute an options struct that implements `utils.HandlerCommand`:

1. **`Complete(ctx *gin.Context) error`**: Extracts everything needed from the HTTP request (params, queries, headers, JSON body payloads) and assigns them to the fields of the struct.
2. **`Validate() error`**: Verifies that the populated data is correct and makes sense (e.g., minimum character length, mutually exclusive fields).
3. **`Run() error`**: Executes the core business logic using the `UseCases` layer and responds back via the Context.

Finally, you run these steps securely using `utils.RunCommand(ctx, options)`.

---

## Handler Template

Use the following template below when creating a new endpoint in a `/handlers` package (e.g. `internal/resources/{resource}/handlers/{action}.go`).

```go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/felipe1496/open-wallet/internal/utils"

	// Replace the following import with your actual resource use cases:
	// "github.com/felipe1496/open-wallet/internal/resources/{resource}/usecases"
)

// -------------------------------------------------------------------------
// 1. Define the Options Struct
// -------------------------------------------------------------------------

type ActionOptions struct {
	// Base Dependencies
	Ctx      *gin.Context
	UseCases interface{} // Replace with e.g., usecases.MyResourceUseCases

	// Request State (extracted later by Complete)
	UserID string

	// Example: Add fields for your specific requests
	// Body PayloadDTO
	// ID   string 
}

// -------------------------------------------------------------------------
// 2. Implement the Lifecycle Methods
// -------------------------------------------------------------------------

// Complete extracts and parses data from the Context
func (o *ActionOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.UserID = ctx.GetString("user_id")

	// Example: bind JSON body
	// if err := ctx.ShouldBindJSON(&o.Body); err != nil {
	//     return err 
	// }

	// Example: get path parameter
	// o.ID = ctx.Param("id")

	return nil
}

// Validate executes strict business rules against the parsed structural state
func (o *ActionOptions) Validate() error {
	// Example validation
	// if o.Body.Name == "" {
	//     return errors.New("name cannot be empty")
	// }
	return nil
}

// Run executes the core database or external actions
func (o *ActionOptions) Run() error {
	// Execute your use cases logic here
	// data, err := o.UseCases.Action(o.ID, o.Body)
	// if err != nil {
	//     return err
	// }

	// Send an HTTP success response
	// o.Ctx.JSON(http.StatusOK, data)
	return nil
}

// -------------------------------------------------------------------------
// 3. The API Wrapper Hook
// -------------------------------------------------------------------------

// @Summary Action Example
// @Description Action Example
// @Tags {resource}
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} interface{} "Success"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /{resource}/action [post]
func (api *API) Action(ctx *gin.Context) {
	cmd := &ActionOptions{
		UseCases: api.useCases, // Inject use cases from the API struct
	}
	utils.RunCommand(ctx, cmd)
}
```
