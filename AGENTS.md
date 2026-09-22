# Project Goal

`rest2go` is a lightweight Go library that provides reusable foundations for REST API microservices. It reduces repeated boilerplate around YAML configuration, HTTP server setup, middleware, API error responses, pagination, filtering, health checks, database connections, and Goose migrations.

## Priorities

- Keep the library small and practical for microservice bootstrap code.
- Preserve simple integration with standard Go HTTP primitives, especially `net/http` and `http.ServeMux`.
- Keep configuration explicit, YAML-based, and extendable by applications through embedded settings structs.
- Maintain consistent API error DTOs and response handling across HTTP features.
- Avoid unnecessary framework-style machinery around business logic owned by consuming applications.

# Mentality

The project is a compact Go library organized around composable infrastructure helpers rather than a full application framework. Most code exposes small constructors, DTOs, interfaces, and functions that consuming services call directly. Control flow is straightforward: applications load settings, create routers and middleware chains, initialize optional database support, and then run their own service logic on top of these primitives.

## Architectural Rules

- `pkg/rest2go/` - public library package for HTTP server setup, middleware, API errors, pagination, filtering, health checks, database provider, and response writer helpers.
- `pkg/rest2go/settings/` - YAML-backed configuration structs, defaults, logging setup, and generic one-time settings loading.

## Preferred Patterns

- Prefer small exported constructors such as `NewServer`, `NewApiError`, `NewPagination`, and `NewPageDto` for library-facing objects.
- Keep dependencies explicit through function arguments and configuration structs rather than hidden application wiring.
- Use standard library types where practical, especially `net/http`, `http.Handler`, `http.HandlerFunc`, `http.ServeMux`, `database/sql`, `log/slog`, and Go errors.
- Model middleware as `func(http.Handler) http.HandlerFunc` and compose chains through `Middlewares`.
- Keep YAML field names and documented defaults synchronized between `pkg/rest2go/settings`, README examples, and tests.
- Keep public API changes consistent with README documentation and changelog expectations.

## Undesired Patterns

- Do not introduce a full web framework abstraction over `net/http` without an explicit project decision.
- Do not move application business logic into this library; keep `rest2go` focused on reusable infrastructure helpers.
- Do not add broad refactors or new abstraction layers when a small package-level function or struct matches the existing style.
- Do not silently change public JSON/YAML field names, default settings, error codes, or documented behavior without updating tests and documentation.

## Implementation Rules

- Prefer the smallest possible change.
- Do not perform refactoring unrelated to the task.
- Do not create new abstractions without a clear need.
- Preserve consistency with the existing code style.

# Project Analysis

Before analyzing code:

1. Read `ai/snapshot/index.json`.
2. Find the relevant functionality snapshot.
3. Only then analyze the code.

# Implementing New Features

1. Find the most similar existing functionality.
2. Treat it as the pattern.
3. Preserve the existing code organization.

# Constraints

- Do not scan the whole repository without need.
- Prefer analysis of specific files.
- Token efficiency is the priority.

# Important Project Decisions
