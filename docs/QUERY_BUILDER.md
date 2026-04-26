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

### Conjunctive Normal Form (CNF / NFC) Inspiration

A major dilemma when designing parsing for web query parameters is deciding **how much complexity to allow with parentheses**. Allowing infinitely nested parentheses like `((A and B) or (C and (D or E)))` makes parsing incredibly volatile and complex, while leaving the API vulnerable to slow logic or deep nesting attacks.

To solve this, our QueryBuilder is strictly designed around **Conjunctive Normal Form** (CNF), sometimes referred to as Forma Normal Conjuntiva (FNC) or Normal Form Conjunctive (NFC) depending on the regional academic translation.

**What is CNF?**
CNF dictates that any boolean logic can be represented strictly as an `AND` of multiple `OR` groups:
`(A OR B) AND (C OR D) AND E`

**How to make ANY boolean operation with CNF:**
Because of the rules of Boolean Algebra (specifically, Distributive Law), **any logical expression can be converted into CNF**. 

For example, what if you conceptually want `(A AND B) OR C`? The parser does not support `OR` logic wrapping an `AND`. However, by distributing the `OR` over the `AND`, you convert it natively to CNF:
`(A OR C) AND (B OR C)`

Translating this to the Web Syntax:
Instead of `(status eq 'active' and age gt 20) or category eq 1`, you write:
`(status eq 'active' or category eq 1) and (age gt 20 or category eq 1)`

This architectural decision means our API remains Turing-complete in logical filter capability, while the parser stays flat, predictable, and extremely performant without needing recursive Abstract Syntax Trees (ASTs) for parameter nesting!

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
- `in`: In List (`id in ('id1', 'id2', 'id3')`)

**logical Operators:**
- `and`: Combine multiple conditions (`active eq true and category_id eq 1`)
- `or`: Group conditions inside parentheses (`(status eq 'pending' or status eq 'failed') and active eq true`)

**Value Types:**
- **Strings**: Must be enclosed in single quotes (`'value'`). To include a single quote in a value, double it (`'O''Brien'`).
- **Numbers**: Raw numbers (`10`, `25.5`).
- **Booleans**: `true` or `false` (lowercase).
- **Null**: `null` (lowercase, without quotes).
- **Lists**: Comma-separated values enclosed in parentheses (`(val1, val2, ...)`). Values inside the list can be any of the above types.

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
        "status": {AllowedOperators: []string{"eq", "in"}},
        "age":    {AllowedOperators: []string{"gt", "lt", "in"}},
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

---

## Swagger Documentation Automation

To keep API documentation in sync with `ParseConfig`, the project uses an automated generator.

### 1. Usage

Annotate your `ParseConfig` variable with `// @gen_swagger_filter`:

```go
// @gen_swagger_filter
var MyFilterConfig = querybuilder.ParseConfig{
    AllowedFields: map[string]querybuilder.FieldConfig{
        "amount": {AllowedOperators: []string{"eq", "gt"}},
    },
    AllowedSortFields: []string{"created_at"},
}
```

### 2. Generating

Run the following command to scan the codebase and update the `@Param filter` and `@Param order_by` annotations in your handlers:

```bash
make gen-swagger-filters
```

This command is also automatically executed as part of `make gen-docs`.

### 3. CI Enforcement

The CI pipeline runs `make check-filters` to ensure that the documentation in the code matches the actual configuration. If you update a `ParseConfig` but forget to run the generator, the build will fail.
