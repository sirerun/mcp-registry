# MCP Registry

A curated registry of OpenAPI specifications for use with the [mint](https://github.com/sirerun/mint) CLI. Mint converts OpenAPI specs into MCP (Model Context Protocol) servers, and this registry provides a catalog of verified API specs that work well with mint.

## How It Works

The registry is a single `registry.json` file containing metadata for each API, including:

- **name** -- unique identifier used by mint (e.g., `mint generate stripe`)
- **spec_url** -- URL to the OpenAPI specification
- **auth_type** -- authentication method (`bearer`, `api_key`, `oauth2`, `none`)
- **auth_env_var** -- environment variable for the API credential
- **tags** -- categorization labels

## Using with mint

```bash
# List available APIs from the registry
mint registry list

# Search by tag
mint registry search --tag payments

# Generate an MCP server from a registry entry
mint generate stripe
```

When you run `mint generate <name>`, mint looks up the entry in this registry, downloads the OpenAPI spec from `spec_url`, and generates a fully functional MCP server.

## Available APIs

| Name | Category | Auth Type |
|------|----------|-----------|
| twitter-v2 | Social | Bearer |
| github | Social / Dev Tools | Bearer |
| stripe | Payments | Bearer |
| slack | Messaging | Bearer |
| discord | Messaging | Bearer |
| openai | AI/ML | Bearer |
| anthropic | AI/ML | API Key |
| jira | Dev Tools | Bearer |
| linear | Dev Tools | Bearer |
| notion | Dev Tools | Bearer |
| gitlab | Dev Tools | Bearer |
| datadog | Infrastructure | API Key |
| pagerduty | Infrastructure | Bearer |
| cloudflare | Infrastructure | Bearer |
| sendgrid | Communication | Bearer |
| twilio | Communication | API Key |
| shopify | E-commerce | Bearer |
| mixpanel | Analytics | API Key |
| segment | Analytics | Bearer |
| petstore | Testing | API Key |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to add a new API entry.

## License

Apache 2.0 -- see [LICENSE](LICENSE).
