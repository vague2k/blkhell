# AGENTS.md

## Project Overview

Blkhell is a tool for a small music label for managing bands, band's
releases, band's projects and projects for the label.

An internal CLI is also being developed alongside the webapp to be used by the
developer for ease of use when developing.

## General principles

- Simplicity is very important
- Clarity is important
- Consistency is important

## Styling

### Naming Enforcement (Read This)

THIS RULE IS MANDATORY FOR AGENT WRITTEN CODE.

- Take into account Go specific conventions
- Use single word names by default for new locals, params, and helper functions.
- Multi-word names are allowed only when a single word would be unclear or ambiguous.

## Error handling (Read this)

- All errors should be checked, even when a method says error is always nil.
- Use `serverErrors.ErrDb` and `serverErrors.ErrInternal` from `server/errors`
  for DB and internal server errors.
- Only check for `sql.ErrNoRows` explicitly when a query returns a single items
  that may not exist (e.g., `SELECT ... WHERE id = ?`). Insert, update, delete,
  and select queries that return slices don't need this check.

## Instructions

- Any new package added to the project must be documented in both AGENTS.md and README.md.

## Tech stack
- `Go`, as the main language
- `Templ`, for html templating
- `HTMX`, for frontend interactivity
- `Tailwind`, for styling (using a binary)
- `sqlite`, as the database
- `sqlc`, for query generation 
- `golang-migrate`, for database migrations
- `just`, as an alternative to make

## Main project structure

```
.
├── cmd/                # Application entrypoints
│   ├── cli/            # CLI binary
│   └── server/         # HTTP server binary
│
├── internal/           # Private packages
│   └── blkhell/        # blkhell CLI implementation details
│
├── server/             # Server application code
│   ├── data/           # misc. structures
│   ├── database/       # Database layer and persistence logic
│   ├── errors/         # Shared error definitions
│   ├── handlers/       # HTTP handlers grouped by route/action
│   ├── middleware/     # HTTP middleware (auth, redirects)
│   ├── services/       # Application logic used by handlers
│   ├── router.go       # Route definitions and middleware wiring
│   └── server.go       # Server bootstrap and configuration
│
└── views/              # Frontend (templ + static assets)
    ├── assets/         # CSS, JS, fonts, images
    ├── components/     # Reusable templ components
    ├── layouts/        # Page layouts / shells
    ├── pages/          # User-facing pages
    └── templui/        # TemplUI component library
```
