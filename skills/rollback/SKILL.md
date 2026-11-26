---
name: rollback
description: Restores a broken nesly production deployment by repointing traffic to a previous revision, and finds the first bad revision. Use when production is down, a bad deploy shipped, or the user asks to roll back, revert, undo, geri al, önceki sürüme dön, or says the last deploy broke the site.
---

# Rolling back a nesly deployment

Rollback **repoints traffic** at a deployment that already exists.
There is no rebuild. Prefer it when production is broken and service must come
back before a code fix.

## Auth first (mandatory)

Call `whoami` before listing deployments or rolling back. If `signed_in` is
false, read `reason` and act as the **setup** skill describes
(`references/auth.md`): `api-unreachable` is a stored session and a quiet API,
so retry and never ask for a login; `no-credential` / `rejected` /
`credential-unreadable` → run `nesly login` yourself; `account-mismatch` is the
wrong account. On MCP 401 use the CLI.

## 1. Confirm the project and list revisions

```
nesly project info --json
nesly ls --limit 10 --json
```

MCP `get_project` (its `deployments`), or `call_api` `DeploymentService.ListDeployments`. Pick a **previous** deployment that was
healthy, not the live bad one. Ids look like `dep_…`; never invent one.

## 2. Restore service

```
nesly rollback <deployment-id> --json
```

MCP `rollback` with that id. Tell the person which deployment is live again and
that this did not rebuild or re-run migrations.

## 3. Verify

```
nesly project info --json
nesly diagnose --json
```

GET the live URL and see it answer. If 5xx continue, you rolled to the wrong
revision: list again and go further back. Never roll back to the currently live
bad id hoping it restarts.

## 4. Find the first bad revision (optional)

When several deploys sit between good and bad and the person needs the
breaking change named:

```
nesly bisect --good <good-id> --bad <bad-id> -- \
  curl -sf "$NESLY_URL/health"
```

MCP `bisect` (local `command` argv; exit 0 = fine). Nothing is deployed or
rolled back during a bisect; each candidate already has its own URL in
`NESLY_URL`.

## 5. After service is stable

Fix the root cause with the **debug** skill (diagnose, traffic, logs), then
deploy a new revision. Do not leave production on an emergency rollback without
a forward fix unless the person asks to stop there.

## Contract

- Rollback is not delete. Deletes (`addon rm`, project delete) are unrelated:
  explicit yes first, then the phone approval (exit 6).
- Reply directly, in the person's language: the live URL first, which revision
  is live again, at most one question.
- Exit codes and JSON fields: the **deploy** skill's `references/output.md`.
