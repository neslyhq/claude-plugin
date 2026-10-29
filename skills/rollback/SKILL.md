---
name: rollback
description: Restore a broken nesly production deployment by rolling back to a previous revision. Use when production is down, a bad deploy shipped, or the user asks to undo / revert / rollback.
---

# Rolling back a nesly deployment

Rollback **repoints traffic** at a deployment that already exists.
There is no rebuild. Prefer this when production is broken and service
must come back before a code fix.

## Auth first (mandatory)

Call `whoami` before listing deployments or rolling back. If not signed
in: stop; `nesly login` + `nesly mcp setup`, restart the agent. On MCP
401: stop retrying MCP; use the CLI.

## 1. Confirm the project and list revisions

```
nesly project info --json
nesly ls --limit 10 --json
```

MCP: `get_project`, `list_deployments`.

Pick a **previous** deployment that was healthy (not the live bad one).
IDs look like `dep_…`.

## 2. Restore service

```
nesly rollback <deployment-id> --json
```

MCP: `rollback` with that deployment id.

Tell the user which id is live again and that this did **not** rebuild
or re-run migrations.

## 3. Verify

```
nesly project info --json
nesly diagnose --json
```

Optionally hit the live URL. If 5xx continue, you rolled to the wrong
revision — list again and roll to an older good one.

## 4. Find the first bad revision (optional)

When several deploys sit between “good” and “bad” and the user needs
the breaking change named:

```
nesly bisect --good <good-id> --bad <bad-id> -- \
  curl -sf "$NESLY_URL/health"
```

MCP: `bisect` (local `command` argv; exit 0 = fine). Nothing is
deployed or rolled back during bisect — each candidate already has its
own URL via `NESLY_URL`.

## 5. After service is stable

Fix the root cause (use the **debug** skill: diagnose → traffic →
logs), then deploy a **new revision**. Do not leave production on an
emergency rollback without a forward fix unless the user asks to stop
there.

## Contract

- Never invent a deployment id.
- Never `rollback` the currently live bad id hoping it “restarts”.
- Destructive project/addon deletes are unrelated — ask first; rollback
  is not delete.
