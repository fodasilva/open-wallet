# Repository Generator Documentation

This project uses a custom script-based generator to automate the creation of repository operations (Select, Insert, Update, Delete, Count). The generator reads metadata from your Go model files and uses templates to produce the corresponding boilerplate.

## How to use

1. Define your model structs in a `models.go` file.
2. Add the generator tags as comments above your structs.
3. Run the generator script or use the make target:
   ```bash
   ./scripts/repository/gen-repos.sh
   # OR run with make
   make gen-repo
   ```
4. **Mandatory Step**: Create a `repository.go` file in the same folder. Copy the example below, adjust the names, and keep ONLY the methods that you defined in your `@method` tags.

## Repository Interface Example

Copy and paste this into your `repository.go` and adjust as needed:

```go
package <your_package>

import (
    "github.com/felipe1496/open-wallet/internal/utils"
)

// Repository interface. Make sure to only include methods
// that you defined with @method tags in models.go
type <Entity>Repo interface {
    Select(db utils.Executer, filter *utils.QueryOptsBuilder) ([]<Entity>, error)
    Insert(db utils.Executer, data Create<Entity>DTO) (<Entity>, error)
    Update(db utils.Executer, data Update<Entity>DTO, filter *utils.QueryOptsBuilder) error
    Delete(db utils.Executer, filter *utils.QueryOptsBuilder) error
    Count(db utils.Executer, filter *utils.QueryOptsBuilder) (int, error)
}

// Implementation struct. Name must match @name tag in models.go
type <Entity>RepoImpl struct {
}

func New<Entity>Repo() <Entity>Repo {
    return &<Entity>RepoImpl{}
}
```

## Metadata Tags

### Root Tags

- `@gen_repo`: Mandatory tag to mark the file for generation.
- `@table: <table_name>`: The name of the database table.
- `@entity: <struct_name>`: The main struct name used for return types.
- `@name: <interface_impl_name>`: The name of the repository implementation struct (e.g., `UsersRepoImpl`).

### Method Tags

Methods are defined with the `@method:` tag, followed by a pipe-separated list of properties.

```go
// @method: <Operation> | fields: <column>:<Field>, ... | payload: <DTO> | return: <column>:<Field>, ...
```

#### Operations supported:

- `Select`: Generates a `FindAll` style method with query options support.
- `Insert`: Generates an `Insert` method returning the new entity.
- `Update`: Generates an `Update` method for partial updates.
- `Delete`: Generates a `Delete` method using a filter.
- `Count`: Generates a `Count` method using a filter.

#### Properties:

- `fields`: A comma-separated list of `db_column:GoFieldName`. Use a `?` suffix on `GoFieldName` to mark it as optional.
- `payload`: The name of the DTO struct used as input for `Insert` and `Update`.
- `return`: For `Insert`, defines which columns to return via `RETURNING` clause and how to scan them into the entity.

## Handling Optional / Nullable Fields

When a field is marked with a `?` in the `@method` definition (e.g., `name:Name?`), the generator assumes its type in the payload is `utils.OptionalNullable[T]`.

### Example Payload

```go
type UpdateUserDTO struct {
    Name  utils.OptionalNullable[string]
    Color utils.OptionalNullable[string]
}
```

### Generated Logic

For fields marked with `?`:

- **Insert**: Only included in the SQL `INSERT` if `.Set` is true.
- **Update**: Only included in the SQL `SET` clause if `.Set` is true.

This allows for:

- Omitting a field from the query (keeping current DB value).
- Explicitly setting a field to `NULL` (by setting `.Set = true` and `.Value = nil`).
- Updating a field to a new value.

### Utility Functions

Use the following utilities to create `OptionalNullable` values:

- `utils.NewValue(v)` Sets a value.
- `utils.NewNull[T]()` Sets to NULL.
- `utils.Unset[T]()` Skips the field.

## Example Configuration

```go
// @gen_repo
// @table: categories
// @entity: Category
// @name: CategoriesRepoImpl
// @method: Select | fields: id:ID, user_id:UserID, name:Name, color:Color, created_at:CreatedAt
// @method: Insert | fields: user_id:UserID, name:Name, color:Color | return: id:ID, user_id:UserID, name:Name, color:Color, created_at:CreatedAt | payload: CreateCategoryDTO
// @method: Update | fields: name:Name?, color:Color? | payload: UpdateCategoryDTO
```

## Running the Generator

The generator scans for all `.go` files containing the `@gen_repo` tag (excluding the `scripts/` directory).

```bash
# Update all repositories
./scripts/repository/gen-repos.sh
# OR run with make
make gen-repo

# Update a specific file
./scripts/repository/gen-repo.sh internal/resources/categories/repository/models.go
```
