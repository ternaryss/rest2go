# Changelog

History of **rest2go** library with changes description.

## Next release

- `FEATURE` REST API for fetching actual configuration
- `FEATURE` More visible HTTP response log and processing duration
- `FEATURE` More visible HTTP request log
- `FEATURE` Database handling refactor
- `FIX` Ant Pattern matching for ApiKeyAuthMiddleware
- `FEATURE` Initialize AI coding assistant (AGENTS.md and repository indexing)

## 1.1.1 (2026-04-14)

- `FIX` SQLite3 respect foreign keys

## 1.1.0 (2025-12-03)

- `FEATURE` Handle any expected HTTP error as unknown
- `FEATURE` Handle expected HTTP 500 error
- `FIX` Not public HTTP server implementation

## 1.0.0 (2025-11-03)

- `FEATURE` Health check
- `FEATURE` Filtering
- `FEATURE` Pagination
- `FEATURE` Database migrations by Goose
- `FEATURE` SQLite and Postgres database connection
- `FEATURE` Preconfigured HTTP server
- `FEATURE` Api-Key header authorization middleware
- `FEATURE` Log HTTP request & response middleware
- `FEATURE` REST API errors handling
- `FEATURE` Loaded once settings as YAML file
- `FEATURE` GoLang library project initialization
