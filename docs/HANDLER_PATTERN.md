# API Handler Pattern

This project implements a **Command Pattern** for HTTP request handlers. The structure is heavily inspired by Kubernetes (`kubectl`) source code and cleanly separates the concern of parsing requests, validating state, and executing business logic.

By using this standardized structure, we avoid massive, undocumented monolithic functions and duplicated error handling.

## The Lifecycle

Instead of writing everything inside an `api.MyHandler(w http.ResponseWriter, r *http.Request)` function, every handler must define and execute an options struct that implements `utils.HandlerCommand`:

1. **`Complete(w http.ResponseWriter, r *http.Request) error`**: Extracts everything needed from the HTTP request (params, queries, headers, JSON body payloads) and assigns them to the fields of the struct.
2. **`Validate() error`**: Verifies that the populated data is correct and makes sense (e.g., minimum character length, mutually exclusive fields).
3. **`Run() error`**: Executes the core business logic using the `UseCases` layer and responds back via the Context.

Finally, you run these steps securely using `utils.RunCommand(ctx, options)`.

---

## API Documentation (Swagger)

We use **initial/docs** and **swag** to automatically generate OpenAPI/Swagger documentation. This is critical for:

- **API Readability**: Allows anyone to understand the contract without reading the Go source code.
- **Frontend Synchronization**: Frontend developers can use the Swagger UI to see available endpoints, models, and error codes.
- **Transparency**: Provides a single source of truth for the API state.

### How to document

Every handler method (`The API Wrapper Hook`) must have the appropriate `@Summary`, `@Description`, `@Tags`, and `@Router` comments as shown in the template above.

### Updating Documentation

After adding or modifying a handler, regenerate the Swagger files:

```bash
make gen-docs
```

The updated documentation will be available at `/api-docs/index.html` when running the application.

---

## The Resource Pattern & Mapping

To maintain a clean separation between the database schema and the API contract, we follow the **Resource Pattern**. This ensures that changes to the database do not accidentally break the API.

### 1. Dedicated Resource Types
Every handler package must define dedicated types for API responses in `handlers/types.go`. These types should:
- Use the **`Resource`** suffix (e.g., `TransactionResource`, `UserResource`).
- Include only the fields necessary for the API.
- Use explicit `json` tags.
- Use **`binding:"required"`** for all mandatory fields. This is critical because our frontend client generator (`swagger-typescript-api`) will otherwise mark these fields as optional (`field?`) in TypeScript, leading to type safety issues.

### 2. Standardized Response Wrappers
All API responses must be wrapped using the generic structures in `internal/util`:
- **Single Item/Action**: `utils.ResponseData[T]`
- **Lists/Pagination**: `utils.PaginatedResponse[T]`

### 3. Mapper Functions
Transformation logic (converting repository or use case structs to Resource structs) must be kept out of the main handler logic. 
- Create a `utils.go` file inside the `handlers/` directory.
- Implement functions following the naming convention: `Map<Entity>Resource(data repo.Entity) handlers.EntityResource`.

**Example:**
```go
// internal/resources/transactions/handlers/utils.go
func MapTransactionResource(t repository.Transaction) TransactionResource {
    return TransactionResource{
        ID:        t.ID,
        Name:      t.Name,
        CreatedAt: t.CreatedAt,
    }
}
```

---

## Handler Templates

### Template Parameters

After copying a template below to your file, replace every occurrence of the following placeholders:

| Placeholder | Description | Example |
|---|---|---|
| `<entities>` | Lowercase plural resource name | `categories` |
| `<Entities>` | PascalCase plural resource name | `Categories` |
| `<Action>` | PascalCase name of the handler action | `Create`, `List`, `Update` |
| `<action>` | Lowercase action name (used in routes) | `create`, `list`, `update` |

---

### `handlers.go` — Entry Point

Every handlers package must have a `handlers.go` file that declares the `API` struct and its constructor. This is the dependency injection root for all handlers in the resource.

Create it at `internal/resources/<entities>/handlers/handlers.go`:

```go
// Parameters to replace:
// - <entities>  → lowercase plural resource name  (e.g. categories)
// - <Entities>  → PascalCase plural resource name (e.g. Categories)

package handlers

import (
	"github.com/felipe1496/open-wallet/internal/resources/<entities>/usecases"
)

type API struct {
	<entities>UseCases usecases.<Entities>UseCases
}

func NewHandler(<entities>UseCases usecases.<Entities>UseCases) *API {
	return &API{
		<entities>UseCases: <entities>UseCases,
	}
}
```

> Each individual endpoint (e.g. `create.go`, `list.go`) lives in its own file within the same package and accesses use cases via `api.<entities>UseCases`.

---

### `<Action>.go` — Individual Endpoint

Use the following template when creating a new endpoint in a `/handlers` package (e.g. `internal/resources/<entities>/handlers/<Action>.go`).

```go
// Parameters to replace:
// - <entities>  → lowercase plural resource name  (e.g. categories)
// - <Entities>  → PascalCase plural resource name (e.g. Categories)
// - <Action>    → PascalCase handler action name  (e.g. Create, List)
// - <action>    → lowercase action name           (e.g. create, list)
// - <endpoint>  → HTTP endpoint path              (e.g. /categories)
// - <method>    → HTTP method                     (e.g. POST, GET, PUT, DELETE)
// - <status>    → HTTP status code                (e.g. 200, 201, 400, 401, 404, 500)
// - <response>  → HTTP response type              (e.g. interface{}, transactions.ListEntriesResponse)

package handlers

import (
	"net/http"

	"net/http"
	"github.com/felipe1496/open-wallet/internal/util"

	// Replace the following import with your actual resource use cases:
	// "github.com/felipe1496/open-wallet/internal/resources/<entities>/usecases"
)

// -------------------------------------------------------------------------
// 1. Define the Options Struct
// -------------------------------------------------------------------------

type <Action>Options struct {
	// Base Dependencies
	W        http.ResponseWriter
	R        *http.Request
	UseCases interface{} // Replace with e.g., usecases.<Entities>UseCases

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
func (o *<Action>Options) Complete(w http.ResponseWriter, r *http.Request) error {
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
func (o *<Action>Options) Validate() error {
	// Example validation
	// if o.Body.Name == "" {
	//     return errors.New("name cannot be empty")
	// }
	return nil
}

// Run executes the core database or external actions
func (o *<Action>Options) Run() error {
	// Execute your use cases logic here
	// data, err := o.UseCases.Action(o.ID, o.Body)
	// if err != nil {
	//     return err
	// }

	// Send an HTTP success response using the generic wrapper
	// o.Ctx.JSON(http.StatusOK, utils.ResponseData[<Action>ResponseData]{
	//     Data: <Action>ResponseData{
	//         <Entity>: Map<Entity>Resource(data),
	//     },
	// })
	return nil
}

// -------------------------------------------------------------------------
// 3. The API Wrapper Hook
// -------------------------------------------------------------------------

// @Summary <Action> Example
// @Description <Action> Example
// @Tags <entities>
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success <status> {object} utils.ResponseData[<Action>ResponseData] "Success"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router <endpoint> [<method>]
func (api *API) <Action>(w http.ResponseWriter, r *http.Request) {
	cmd := &<Action>Options{
		UseCases: api.<entities>UseCases,
	}
	utils.RunCommand(ctx, cmd)
}
```
