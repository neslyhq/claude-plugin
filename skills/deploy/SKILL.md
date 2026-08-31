---
name: deploy
description: Deploy an application to nesly and get it live — use when the user wants to deploy, host or publish an app or site, or when a nesly deployment failed and needs fixing.
---

# Deploying to nesly

The `nesly` CLI is the whole interface. Every command accepts `--json`
(exactly one JSON object on stdout, progress on stderr; exit codes:
0 ok, 1 error, 2 usage, 3 auth, 4 not found, 5 conflict). The MCP
tools from this plugin (`deploy_project`, `diagnose_deployment`, …)
are the same operations; use whichever is closer at hand.

If `nesly` is not installed: `brew install neslyhq/tap/nesly`
(macOS/Linux). Signed in? `nesly login` opens a browser; headless
environments use a token: `nesly token create agent` once, then
`NESLY_TOKEN=nesly_pat_...`.

## The application contract

One rule: listen on `$PORT`, bind `0.0.0.0`. Static sites (or anything
with a build output directory) need no server at all — the platform
serves the files.

## Deploy

From the project directory:

```
nesly deploy --json
```

No project yet? `nesly new` scaffolds from a template, or deploy an
existing directory and a project is created on first deploy. The build
is detected automatically (Node frameworks, Python, Go, .NET, Docker,
Hugo, Jekyll, MkDocs). Before deploying something unusual, preview the
plan:

```
nesly build show
```

If it prints "a deploy would be refused", fix that first — the same
check would reject the deploy before building. Override detection in
`nesly.toml` under `[build]` (keys outside `[build]` are ignored —
diagnose will tell you) or with `nesly build set`.

## When a deploy fails

Do not start with raw logs. Diagnose reads them for you:

```
nesly diagnose --json
```

It returns the probable cause, evidence lines from the log, env vars
the app reads but the project does not define, and the command that
fixes each finding. Apply the fix, deploy again. `nesly logs` exists
for when the diagnosis is not enough.

## Everything after the first deploy

- Env vars: `nesly env set KEY=VALUE`, `nesly env ls`
- Databases: `nesly addon add postgres` or `nesly addon add mssql` —
  both inject `DATABASE_URL`; mssql also `ConnectionStrings__Default`
- Backups: `nesly backup create`, `nesly backup ls`,
  `nesly backup restore <backup-id>`
- Custom domain: `nesly domain add example.com` prints the DNS records,
  `nesly domain verify example.com` activates
- Push-to-deploy: `nesly git connect <owner/repo>`
- Instant rollback: `nesly rollback <deployment-id>`
- Destructive commands (`addon rm`, `backup rm`, project deletion)
  permanently delete data — confirm with the user first.

Full agent contract: https://nesly.net/docs/agents/ and
https://nesly.net/llms.txt
