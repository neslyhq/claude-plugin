---
name: debug
description: Debugs a live nesly app — production 500s, 502s, failed requests, crashes, slow routes, a deployment that is not answering. Use when the user reports errors, outages, 5xx, hata, çöktü, yavaş, "site açılmıyor", logs, or asks why a deployed project fails; diagnoses first, then reads traffic and request logs, then fixes and verifies.
---

# Debugging a live nesly deployment

Do not start with raw logs or by redeploying: the platform already names the
likely cause. Use the `nesly` CLI (`--json` on every command) or the matching
MCP tools; same credentials either way.

## Auth first (mandatory)

Call `whoami` (MCP or `nesly whoami --json`) before diagnose or logs. If
`signed_in` is false, branch on `reason` as the **setup** skill describes
(`references/auth.md`): `api-unreachable` → wait and retry, never ask for a
login; `no-credential` / `rejected` / `credential-unreadable` → run `nesly
login` yourself; `account-mismatch` → `nesly account ls`. On MCP 401 use the
CLI in the same turn.

## 1. Resolve the project

`.nesly/project.json`, or `nesly project info --json` / MCP `get_project`, or
`nesly project list --json` / MCP `list_projects`.

## 2. Diagnose first (mandatory)

```
nesly diagnose --json
```

MCP `diagnose`. It returns the probable cause, the supporting log
lines, missing env vars and a fix command per finding. Apply those before
inventing a theory. Skipping diagnose for a failed or unhealthy deployment is a
**contract violation**.

For a deployment that failed to build or start, the failure branches (the
platform's own repair and `nesly fix show`, addon for missing services,
self-generated secrets (the platform makes them), the person's own third-party secrets, redeploy and
verify the URL) are the loop in the **deploy** skill: follow it, then return
here only if the live site still misbehaves.

## 3. Live traffic (the site is "up" but visitors hit errors or slowness)

```
nesly traffic --sort errors --since 1h --json
nesly traffic failures --since 1h --json
nesly traffic --sort slowest --since 1h --json
```

MCP: `call_api` `TrafficService.GetTraffic`, `ListRequestFailures`, `GetTrafficSeries`. Status codes:

- **500** — the app threw; fix code or env, then redeploy
- **502** — nothing listening (crash, bad start command, timeout)
- **503** — no capacity or a cold start; not a missing env var
- **404** — a visitor path, not an outage

For one failing request take `X-Nesly-Request-Id` from the response (or from
`traffic failures`) and narrow the logs: `nesly logs --request <id> --json`
(MCP `get_logs`, `request=<id>`).

## 4. Fix, then verify

1. Apply the fix: `env set`, an addon, a start command, or a change to the
   app's own code **only when the person asks**. Never a dependency version,
   lockfile or configuration change to make something pass; if the cause is
   there, tell the person which line and let them decide.
2. Redeploy after the fix (`nesly deploy --preview --detach --yes --json`, then
   `nesly wait <id> --json`): already-live rules apply, this is a new revision.
3. GET the URL and see it answer, then `nesly diagnose --json` and/or `nesly
   traffic failures --since 15m --json`.

## 5. Service first

If errors are severe and the site must be healthy before a root-cause fix,
use the **rollback** skill to repoint traffic at a previous deployment, then
come back here.

## 6. Cold starts

Free apps pause when idle; the first request after idle can take seconds. That
is expected on the free tier, not a code bug. Confirm with `traffic --sort
slowest`: if only the first hit after a gap is slow, say so; do not add fake
keep-alives unless the person asks for a paid warm plan.

## Done when

- You stated the cause with a diagnose finding or a traffic row
- You applied a concrete fix, or explained why rollback came first
- The URL answers and a follow-up diagnose / traffic check shows the symptom
  gone or reduced

Reply to the person directly, in their language: URL first, then what was
wrong and what you changed, at most one question. Exit codes and JSON fields:
the **deploy** skill's `references/output.md`.
