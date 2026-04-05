# Query Builder Documentation

## Motivation

Traditional dynamic query implementations often suffer from two major problems:

1.  **Redundant Domain Types**: Developers are usually forced to create unique `Filter` or `Search` structs for every single resource (`UserFilter`, `TransactionFilter`, etc.), leading to massive boilerplate.
2.  **Disconnected Logic**: Handling query parameters often happens exclusively in the controller/middleware, and then the extracted data is passed down through several layers as loose primitives (many arguments).

### The "Living Object" Architecture

The `querybuilder.Builder` serves as a **Living Object** that travels through the entire request lifecycle:

-   **Middleware**: Receives the raw web parameters and performs initial parsing and basic validation.
-   **Usecase**: Can augment the builder with mandatory business rules (e.g., `builder.And("user_id", "eq", session.UID)` to enforce ownership) without knowing about web concerns.
-   **Repository**: Finalizes the process by converting the rich state in the builder into a SQL query using `ToSquirrel`.

This approach provides:
-   **State Composition**: You can compose filters at different levels of the application without overwriting or losing previous constraints.
-   **Zero Boilerplate**: No need to define or maintain custom filter types for each domain.
-   **Layer Isolation**: The repository doesn't need to know about the filter string syntax, and the middleware doesn't need to know about Squirrel or SQL.

---

## Web Syntax

The following query parameters are supported in the API:

### 1. Filtering (`filter`)

The `filter` parameter uses an OData-like syntax.

**Basic Comparisons:**
- `eq`: Equal (`name eq 'John'`)
- `ne`: Not Equal (`status ne 'deleted'`)
- `gt`: Greater Than (`age gt 25`)
- `gte`: Greater Than or Equal (`age gte 25`)
- `lt`: Less Than (`total lt 100`)
- `lte`: Less Than or Equal (`total lte 100`)
- `like`: Case-insensitive partial match (`description like 'rocket'`)

**logical Operators:**
- `and`: Combine multiple conditions (`active eq true and category_id eq 1`)
- `or`: Group conditions inside parentheses (`(status eq 'pending' or status eq 'failed') and active eq true`)

**Value Types:**
- **Strings**: Must be enclosed in single quotes (`'value'`). To include a single quote in a value, double it (`'O''Brien'`).
- **Numbers**: Raw numbers (`10`, `25.5`).
- **Booleans**: `true` or `false` (lowercase).
- **Null**: `null` (lowercase, without quotes).

### 2. Sorting (`order_by`)

Standardized sorting format: `field:direction`.

- **Single field**: `order_by=created_at:desc`
- **Multiple fields**: `order_by=priority:desc,name:asc`
- **Default direction**: If omitted, defaults to `asc` (`order_by=name`).

### 3. Pagination (`page` and `per_page`)

- `page`: The page number (starts at 1).
- `per_page`: Number of records per page (defaults to 10).

---

## Developer Guide

### 1. Using the Middleware

To enable dynamic queries on a route, use the `QueryBuilderMiddleware`. It is highly recommended to provide a `ParseConfig` to restrict which fields can be filtered and sorted.

```go
filterConfig := &querybuilder.ParseConfig{
    AllowedFields: map[string]querybuilder.FieldConfig{
        "name":   {AllowedOperators: []string{"eq", "like"}},
        "status": {AllowedOperators: []string{"eq"}},
        "age":    {AllowedOperators: []string{"gt", "lt"}},
    },
    AllowedSortFields: []string{"name", "created_at"},
}

r.GET("/items", 
    middlewares.QueryBuilderMiddleware(filterConfig), 
    handler.ListItems,
)
```

### 2. Retrieving the Builder in Handlers

The middleware injects the `querybuilder.Builder` into the Gin context.

```go
func (h *Handler) ListItems(ctx *gin.Context) {
    builder := ctx.MustGet("query_builder").(*querybuilder.Builder)
    
    // Pass it to the usecase/repository
    items, err := h.usecase.List(builder)
}
```

### 3. Manual Building (Fluent API)

You can also build queries programmatically:

```go
builder := querybuilder.New().
    And("active", "eq", true).
    OrderBy("priority", "desc").
    Limit(10).
    Offset(0)

// With OR groups
builder.InitOr().
    Or("status", "eq", "failed").
    Or("status", "eq", "cancelled").
    EndOr()
```

### 4. Converting to Squirrel (SQL)

The builder integrates seamlessly with `Masterminds/squirrel`.

```go
query := squirrel.Select("*").From("users")
query = querybuilder.ToSquirrel(query, builder)

sql, args, err := query.ToSql()
```

### 5. Repository Integration

Repository methods generated with the `@gen_repo` tag automatically receive and use the `querybuilder.Builder` in their `Select`, `Count`, and `Delete` methods.

```go
// @method: Select | fields: id, name, status
func (r *Repo) Select(db utils.Executer, filter *querybuilder.Builder) ([]Item, error) {
    // ... generated code will use querybuilder.ToSquirrel ...
}
```
