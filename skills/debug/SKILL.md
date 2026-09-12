---
name: debug
description: Debug a live nesly app — production 500s, failed requests, or slow routes. Use when the user reports errors, outages, 5xx, crashes, or slowness on a deployed project.
---

# Debugging a live nesly deployment

Do **not** start with raw logs or by redeploying. The platform already
names the likely cause.

Use the `nesly` CLI (`--json` on every command) or the matching MCP
tools. Same credentials either way.

## 1. Resolve the project

`.nesly/project.json`, or `nesly project info --json` /
MCP `get_project`, or `nesly project list --json` /
MCP `list_projects`.

## 2. Diagnose first (mandatory)

```
nesly diagnose --json
```

MCP: `diagnose_deployment`.

It returns probable cause, evidence lines, missing env vars, and a
fix command per finding. Apply those fixes before inventing a theory.

Skipping diagnose for a failed or unhealthy deployment is a
**contract violation**.

## 3. Live traffic (production errors / slowness)

When the site is “up” but visitors hit errors or it feels slow:

```
nesly traffic --sort errors --since 1h --json
nesly traffic failures --since 1h --json
nesly traffic --sort slowest --since 1h --json
```

MCP: `get_traffic`, `list_traffic_failures`, `get_traffic_series`.

Read status codes carefully:

- **500** — app threw; fix code or env, then redeploy
- **502** — nothing listening (crash, bad start command, timeout)
- **503** — no capacity / cold scale; not a missing env var
- **404** — visitor path; not an outage

For one failing request, take `X-Nesly-Request-Id` from the response
(or from `traffic failures`) and narrow logs:

```
nesly logs --request <id> --json
```

MCP `get_logs` with `request=<id>`.

## 4. Fix, then verify

1. Apply the diagnose / traffic fix (`env set`, code change, start
   command, addon, …).
2. Redeploy only after the fix: `nesly deploy --json` or
   MCP `deploy_project` (already-live rules still apply — this is a
   **new revision**).
3. Re-check: `nesly diagnose --json` and/or
   `nesly traffic failures --since 15m --json`.

## 5. When service must be restored first

If errors are severe and the user needs the site healthy before a
root-cause fix, switch to the **rollback** skill: repoint traffic at a
previous deployment, then return here to debug.

## 6. Cold starts (not “slow functions”)

Free apps pause when idle; the first request after idle can take
seconds. That is expected for the free tier, not a code bug. Confirm
with `traffic --sort slowest`: if only the first hit after a gap is
slow, say so; do not “optimize” with fake keep-alives unless the user
asks for a paid warm plan.

## Done when

- You stated the cause with evidence (diagnose finding or traffic row)
- You applied a concrete fix (or explained why rollback came first)
- A follow-up diagnose / traffic check shows the symptom gone or reduced
