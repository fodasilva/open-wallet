# Contributing to Open Wallet

Thank you for your interest in contributing to Open Wallet! To maintain a high-quality codebase and architectural consistency, please follow these guidelines.

## Architecture Overview

This project follows a modular resource-based architecture with clean separation of concerns:

- **Handlers**: HTTP entry point, request parsing, and validation.
- **Usecases**: Core business logic and orchestration.
- **Repositories**: Database interactions (mostly generated).
- **Services**: External integrations (OAuth, Email, etc.).

## Contribution Workflow

### 1. Architectural Patterns

Before starting any development, please read the documentation for each layer:

- [**Handler Pattern**](docs/HANDLER_PATTERN.md): Command pattern for controllers.
- [**Usecase Pattern**](docs/USECASE_PATTERN.md): Service layer and business logic.
- [**Repository Generator**](docs/REPOSITORY_GENERATOR.md): Automation for database operations.
- [**Query Builder**](docs/QUERY_BUILDER.md): Standardized filtering for repositories.
- [**E2E Testing**](docs/E2E_TESTING_DIRECTIVES.md): Patterns and conventions for writing E2E tests.

### 2. Local Development

Ensure you have the necessary dependencies installed:

```bash
# Install linting tools
make lint-install

# Install documentation tools
make gen-docs-install

# Install migration tools
make db-migrate-install
```

### 3. Database Migrations

All schema changes must be implemented via migrations.

- **Creating a migration**: `make db-migrate-create name=my_migration_name`
- **Applying migrations**: `make db-migrate`

Running these commands requires the `migrate` CLI tool, which you can install using `make db-migrate-install`.

### 3. Adding a New Resource

The standard workflow for adding a new feature involves:

1.  **Define Types**: Create/update `types.go` in the resource folder.
2.  **Generate Repositories**: Mark your types with `@gen_repo` and run `make gen-repo`.
3.  **Implement Usecases**: Create the usecase interface and its implementation.
4.  **Create Handlers**: Use the command pattern to implement your API endpoints.
5.  **Register Routes**: Add your handler to `internal/routes/`.

### 4. Code Quality & Testing

-   **Linting**: Always run `make lint` before committing your changes.
-   **Testing Strategy**: We prioritize **E2E tests** over unit and integration tests, as they provide more value by validating the entire request/response lifecycle and real database interactions. Please read our [**E2E Testing Directives**](docs/E2E_TESTING_DIRECTIVES.md) before writing tests.
    -   **E2E Tests**: Use these for all API endpoints and core business flows. Run `make test` to execute them.
    -   **Unit or Integration Tests**: Should only be used for small, isolated internal helpers (e.g., complex calculations, formatting, or utility functions).
-   **Conventional Commits**: While not strictly enforced by hooks, please use descriptive prefixes like `feat:`, `fix:`, `refactor:`, `docs:`, or `chore:`.

## Getting Help

If you're unsure about a specific pattern or need architectural guidance, please consult the `docs/` folder or reach out to the maintainers.
