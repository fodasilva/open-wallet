# Usecase Pattern

This project follows a clean architecture approach by separating business logic into **Usecases**. Each resource has a `usecases` package that defines the high-level operations available for that resource.

## Pattern Rules

1. **Interface Driven**: All usecases are defined as interfaces to allow for easy mocking in tests when necessary (see our [E2E Testing Directives](E2E_TESTING_DIRECTIVES.md) for actual test implementation strategies).
2. **Implementation Decoupling**: The implementation struct (`<Entities>UseCasesImpl`) is private to the package.
3. **Dependency Injection**: Use cases must receive their dependencies (repositories, services, or other use cases) via their constructor.
4. **Error Handling**: Use cases should return domain-specific errors or `utils.HTTPError` if the error should be propagated directly to the API.

---

## Usecase Templates

### Template Parameters

After copying a template below to your file, replace every occurrence of the following placeholders:

| Placeholder | Description | Example |
|---|---|---|
| `<entities>` | Lowercase plural resource name | `categories` |
| `<Entities>` | PascalCase plural resource name | `Categories` |
| `<Action>` | PascalCase name of the usecase method | `Create`, `GenerateReport` |

---

### `usecases.go` — Entry Point

This file defines the interface and the implementation struct with its constructor.

Create it at `internal/resources/<entities>/usecases/usecases.go`:

```go
// Parameters to replace:
// - <entities> → lowercase plural resource name (e.g. categories)
// - <Entities> → PascalCase plural resource name (e.g. Categories)

package usecases

type <Entities>UseCases interface {
	<Action>(/* params */) (/* return types */, error)
}

type <Entities>UseCasesImpl struct {
	// Add dependencies here:
	// repository <entities>Repo.Repository
}

func New<Entities>UseCases(/* dependencies */) <Entities>UseCases {
	return &<Entities>UseCasesImpl{
		// Inject dependencies
	}
}
```

---

### `<action>.go` — Method Implementation

Use separate files for each complex method to keep the codebase maintainable and avoid monolithic files.

Create it at `internal/resources/<entities>/usecases/<action>.go`:

```go
// Parameters to replace:
// - <entities> → lowercase plural resource name (e.g. categories)
// - <Entities> → PascalCase plural resource name (e.g. Categories)
// - <Action>   → PascalCase method name           (e.g. Create)

package usecases

func (uc *<Entities>UseCasesImpl) <Action>(/* params */) (/* return types */, error) {
	// 1. Business Logic
	// 2. Repository Calls
	// 3. Return results
	return nil
}
```
