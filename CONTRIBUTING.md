# Contributing to MCP Registry

Thank you for contributing! This guide explains how to add a new API entry to the registry.

## Adding a New API

### Prerequisites

- The API must have a publicly accessible OpenAPI 3.x specification
- The spec URL must be stable (prefer raw GitHub URLs or official API endpoints)
- You have verified the spec URL returns valid OpenAPI JSON or YAML

### Required Fields

Each entry in `registry.json` must include:

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Unique lowercase identifier, e.g., `my-api` |
| `description` | string | What the API does (min 10 chars) |
| `tags` | string[] | 1-5 categorization tags |
| `spec_url` | string | URL to the OpenAPI spec (must be reachable) |
| `auth_type` | string | One of: `bearer`, `api_key`, `oauth2`, `none` |
| `auth_env_var` | string | Environment variable name, e.g., `MY_API_KEY` |
| `min_mint_version` | string | Minimum mint version, e.g., `0.2.0` |

### Steps

1. Fork this repository
2. Add your entry to the `apis` array in `registry.json`
3. Validate locally:
   ```bash
   cd tools && go run validate.go
   ```
4. Open a pull request with:
   - A clear title: `feat: add <api-name> to registry`
   - Confirmation that the spec URL is reachable
   - A brief note on why this API is useful

### Entry Template

```json
{
  "name": "my-api",
  "description": "My API for doing useful things with widgets and gadgets",
  "tags": ["category", "subcategory"],
  "spec_url": "https://raw.githubusercontent.com/org/repo/main/openapi.json",
  "auth_type": "bearer",
  "auth_env_var": "MY_API_KEY",
  "min_mint_version": "0.2.0"
}
```

### Naming Conventions

- Use lowercase with hyphens: `my-api`, not `MyAPI` or `my_api`
- Use the well-known name: `stripe`, not `stripe-payments-api`
- Add version suffix only when multiple major versions coexist: `twitter-v2`

### Validation

CI runs automatically on pull requests and checks:

1. `registry.json` conforms to `schema.json`
2. All `spec_url` values are reachable (HTTP HEAD, following redirects)
3. No duplicate `name` values

### Tags

Use existing tags when possible:

- **social**, **messaging**, **communication**
- **payments**, **e-commerce**
- **ai**, **ml**
- **dev-tools**, **vcs**, **project-management**
- **infrastructure**, **monitoring**, **cloud**
- **analytics**, **data**
- **email**, **testing**, **productivity**

## Review Process

Submissions are reviewed using a [Sire workflow](workflows/submission-review.yaml) that runs automatically on each pull request:

1. **Schema validation** -- the entry is checked against `schema.json`
2. **Spec URL check** -- the `spec_url` must return a successful HTTP response
3. **Human approval** -- a maintainer reviews the results and approves or requests changes
4. **Merge** -- once approved, the PR is merged into `main`

Expect a review within a few business days. If your spec URL check fails, update the URL and push a new commit.

## Code of Conduct

Be respectful. Focus on APIs that are broadly useful and well-maintained.
