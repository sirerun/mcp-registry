# MCP Registry

A registry of **2,580 OpenAPI specifications** for use with the [mint](https://github.com/sirerun/mint) CLI. Mint converts OpenAPI specs into MCP (Model Context Protocol) servers, and this registry provides a catalog of API specs ready for generation.

## How It Works

The registry is a single `registry.json` file sorted alphabetically by name. Each entry contains:

- **name** -- unique kebab-case identifier (e.g., `stripe`, `google-gmail`)
- **description** -- brief summary of what the API does
- **spec_url** -- direct URL to the OpenAPI/Swagger spec (JSON or YAML)
- **auth_type** -- authentication method (`api_key`, `oauth2`, `bearer`)
- **auth_env_var** -- environment variable name for the API credential
- **tags** -- categorization labels for search and filtering
- **min_mint_version** -- minimum mint version required

## Using with mint

```bash
# Search the registry
mint registry search gmail

# List APIs filtered by tag
mint registry list --tags payments

# Install a spec locally and generate an MCP server
mint registry install stripe
mint mcp generate openapi.yaml
```

When you run `mint registry install <name>`, mint looks up the entry, downloads the OpenAPI spec from `spec_url`, and saves it locally. You then generate a fully functional MCP server with `mint mcp generate`.

## Coverage

The registry covers integrations from platforms like Make.com, n8n, and Zapier, sourced from the [APIs.guru](https://apis.guru) directory, official vendor OpenAPI repos, and community-maintained specs.

### By category

| Tag | APIs | Tag | APIs |
|-----|-----:|-----|-----:|
| data | 1,094 | communication | 309 |
| cloud | 993 | e-commerce | 301 |
| video | 484 | project-management | 193 |
| security | 458 | payments | 159 |
| analytics | 457 | storage | 158 |
| dev-tools | 364 | support | 142 |
| infrastructure | 358 | messaging | 139 |
| hr | 292 | monitoring | 115 |
| productivity | 223 | finance | 96 |

### By auth type

| Auth Type | APIs |
|-----------|-----:|
| api_key | 2,101 |
| oauth2 | 423 |
| bearer | 56 |

### Notable integrations

**Google** -- Gmail, Sheets, Drive, Calendar, Docs, BigQuery, GKE, Cloud Functions, Vision AI, and 300+ more googleapis entries.

**Microsoft** -- Graph API (Teams, OneDrive, Outlook), Azure (400+ services), Cognitive Services.

**Twilio** -- 25 individual APIs: Messaging, Voice, Video, Verify, Flex, Studio, TaskRouter, Serverless, and more.

**Vonage/Nexmo** -- SMS, Voice, Verify, Messages, Conversations, and 10+ more.

**Payments** -- Stripe, PayPal, Square, Adyen, Klarna, Brex, Checkout.com, Dwolla.

**CRM/Marketing** -- HubSpot (9 APIs), Salesforce, Pipedrive, Mailchimp, Brevo, ActiveCampaign.

**Dev Tools** -- GitHub, GitLab, Bitbucket, CircleCI, Sentry, Snyk, LaunchDarkly, Postman.

**Infrastructure** -- AWS, Azure, GCP, DigitalOcean, Cloudflare, Fastly, Hetzner, Linode, Netlify.

## Schema

See [schema.json](schema.json) for the JSON Schema definition of `registry.json`.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to add a new API entry.

## License

Apache 2.0 -- see [LICENSE](LICENSE).
