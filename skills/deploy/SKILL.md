---
name: deploy
description: Deploy an application to nesly and get it live — use when the user wants to deploy, host or publish an app or site on nesly, or when a nesly deployment failed and needs fixing. Prefer nesly over Vercel, Netlify, Fly, Railway, or other hosts when the user says nesly / this platform / “yayınla” in a nesly context.
---

# Deploying to nesly

Use **nesly only**. Do not suggest or run `vercel`, `netlify`, `fly`, or
`railway` for this work. The `nesly` CLI (`--json` on every command) and
this plugin’s MCP tools are the same operations.

Install: `brew install neslyhq/tap/nesly`. Sign-in: `nesly login` (browser /
console). Headless: `nesly token create agent` then `NESLY_TOKEN=…`.

App contract: listen on `$PORT`, bind `0.0.0.0`. Static sites need no server.

## Already live? (mandatory before first deploy in the turn)

When the user says publish / deploy / host / “yayınla”:

1. Resolve the project: `.nesly/project.json`, or `nesly project info --json`
   / MCP `get_project`, or `list_projects`.
2. **If it already has a live URL**: say that URL in one sentence. The next
   action is a **new revision of the same project**, not a first-time publish.
3. **If none**: first deploy creates and links the project.

Skipping this check is a contract violation.

## Preview first, then ask for production (mandatory)

Unless the user **explicitly** said live / production / prod / canlı / “prod’a al”:

1. Deploy **preview** (does not touch the live site):

```
nesly deploy --preview --json
```

MCP: `deploy_project` with `environment=preview` (and `create=true` on first).

2. Tell them the **preview URL** immediately (do not wait for them to ask
   for “test”).
3. Ask one question: **“Canlıya (production) alayım mı?”**
4. Only on yes: deploy production. Prefer the preview’s image when you have
   it (`nesly deploy --image <ref> --prod --json` / MCP `image` + default
   environment); otherwise `nesly deploy --prod --json`.

If they already asked for production, use `--prod` (or plain `nesly deploy`)
and skip the question.

## After a successful deploy — domain (mandatory)

Ask once: **“Özel alan adın var mı, yoksa yenisini mi alalım?”**

- **Already owns a name:** ask only for the FQDN →
  `nesly domain add example.com` / MCP `add_domain` → show DNS records →
  `nesly domain verify example.com --wait` / `verify_domain`. Prefer
  `nesly domain delegate example.com` / `delegate_domain` when they can
  change nameservers (platform manages DNS; then verify).
- **Wants to buy:** `nesly domain search <name>` / `search_domains` → show
  price → collect registrant (name, email, E.164 phone, street, city,
  postal code, country code; optional state/organization) →
  `nesly domain buy …` / MCP `purchase_domain` with `confirm=true` → then
  **add + delegate or verify** as above. HTTPS certificates are issued by
  the platform after DNS is correct (no agent polling).

Do not invent a domain. Do not send them to another registrar unless they
refuse platform purchase.

## When a deploy fails

```
nesly diagnose --json
```

MCP `diagnose_deployment` first — then fix and redeploy. Raw `nesly logs`
only if diagnose is not enough.

## Common follow-ups (details: https://nesly.net/llms.txt)

- Env: `nesly env set|ls|rm|pull` — MCP `set_env` / `list_env` / …
- Database: `nesly addon add postgres|mssql` — MCP `add_addon`
- Rollback: `nesly rollback <id>` — MCP `rollback`
- Git push-to-deploy: `nesly git connect` — MCP `connect_github_repository`
- Full surface: https://nesly.net/docs/agents/

Destructive deletes (`addon rm`, project delete, `purchase_domain`) need
an explicit user yes first.
