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

## Already live? (mandatory before first deploy in the turn)

When the user says publish / deploy / host / “yayınla”, **do not** jump
straight to `nesly deploy` or `deploy_project`.

1. Resolve the project: `.nesly/project.json` in the working directory,
   or `nesly project info --json` / MCP `get_project`, or
   `nesly project list --json` / MCP `list_projects` (match by name).
2. **If it already has a live URL** (or an active deployment): tell the
   user that URL in one sentence. State clearly that the next action is
   a **new revision of the same project**, not a first-time publish.
   Then deploy only if they asked to publish/update (not if they only
   asked whether it is live).
3. **If there is no project / not linked**: first deploy (create + link
   as usual).

Skipping this check is a contract violation — the user must never be
left wondering whether the agent noticed an existing site.

## Deploy

From the project directory:

```
nesly deploy --json
```

Or a prebuilt image (no local upload/build):

```
nesly deploy --image registry.example.com/app:1.2.3 --json
```

Preview without touching production:

```
nesly deploy --preview --json
```

MCP `deploy_project` accepts `dir` or `image` (exactly one), and
optional `environment=preview`.

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
for when the diagnosis is not enough — MCP `get_logs` (optional
`request` = X-Nesly-Request-Id, `worker` = Procfile process name).

## Everything after the first deploy

- Env vars: `nesly env set KEY=VALUE`, `nesly env ls`, `nesly env rm KEY`,
  `nesly env share KEY --with preview`, `nesly env pull` — MCP `set_env` /
  `list_env` / `delete_env` / `share_env` / `pull_env` (pass `reveal=true`
  on list only when you need values; `share_env` copies production→preview
  on the server and never returns the value; `pull_env` writes a local
  `.env` to an absolute `output` path including derived DB_* — values
  never in the tool answer; `force=true` to replace)
- Databases: `nesly addon add postgres` or `nesly addon add mssql` —
  both inject `DATABASE_URL`; mssql also `ConnectionStrings__Default`
- Backups: `nesly backup create`, `nesly backup ls`,
  `nesly backup restore <backup-id>`, `nesly backup download` —
  MCP `list_backups` / `create_backup` / `restore_backup`
  (`confirm=true`) / `remove_backup` (`confirm=true`) /
  `download_backup` (absolute `output` path; bytes never in the
  tool answer)
- Custom domain: `nesly domain add example.com` prints the DNS records,
  `nesly domain verify example.com` activates, `nesly domain rm` detaches,
  `nesly domain pin <fqdn> <deployment-id>` / `unpin` — MCP `add_domain` /
  `verify_domain` / `list_domains` / `remove_domain` (`confirm=true` on
  remove) / `pin_domain` (omit `deployment_id` to unpin) /
  `delegate_domain` (apex nameserver handoff); availability —
  `nesly domain search <name>` — MCP `search_domains` (no purchase)
- Certificates: `nesly certs` — MCP `list_certificates` (measured TLS
  handshake per host; broken hosts sort first); `upload_certificate`
  (PEM or absolute cert_file/key_file; private key never returned) /
  `remove_certificate` (`confirm=true`)
- Push-to-deploy: `nesly git connect <owner/repo>` — MCP
  `get_github_status` / `list_git_repositories` then
  `connect_github_repository` / `disconnect_github_repository`
  (`confirm=true`); legacy `connect_repository` (https URL) remains
- Instant rollback: `nesly rollback <deployment-id>`
- Find the break: `nesly bisect --good <id> --bad <id> -- <check>` —
  MCP `bisect` (local `command` argv; `NESLY_URL` = candidate preview;
  exit 0 = fine; nothing is deployed)
- Gradual rollout: `nesly deploy --percent 10`, `nesly rollout`,
  `nesly rollout set|promote|abort` — MCP `deploy_project` with
  `percent`, then `get_rollout` / `set_rollout` / `promote_rollout` /
  `abort_rollout`
- Remove deployment: `nesly rm <deployment-id>` — MCP `delete_deployment`
  (`confirm=true`; not the live version)
- Image scan: `nesly scan` — MCP `get_deployment_scan` (registry vulns;
  static = not_applicable, not clean; optional min_severity / summary_only)
- Traffic: `nesly traffic` / `traffic failures` — MCP `get_traffic` /
  `list_traffic_failures` / `get_traffic_series` (routes + per-request
  5xx + time buckets; not billing)
- Usage over time: MCP `get_usage_series` (UTC day buckets; CLI has
  totals only via `nesly usage` / MCP `get_usage`)
- Limits: `nesly limits` — MCP `get_limits` (plan ceilings vs used)
- Billing (read): MCP `get_subscription` / `list_plans` / `list_invoices`
  (card brand/last4 only — never PAN; no plan change here)
- Protection: `nesly protect show|set|off` — MCP `get_protection` /
  `set_protection` / `clear_protection` (rate limit; get also shows
  country/IP/request/preview snapshot); countries —
  `set_country_rule` / `clear_country_rule` (allow|block + ISO codes);
  IP — `set_ip_rule` / `clear_ip_rule` (prefer block over allow);
  request rules — `add_request_rule` / `remove_request_rule` /
  `clear_request_rules` (path|agent|header refusal only);
  preview — `set_preview_protection` / `mint_preview_bypass_token` /
  `clear_preview_protection` (published site untouched)
- List / inspect: `nesly project list` / `nesly project info` — MCP
  `list_projects` / `get_project` (single-screen status: live, domains,
  addons, maintenance); `set_maintenance` (enabled + optional message);
  `whoami` (signed-in account, team, role, API); `list_accounts` /
  `list_members` /
  `invite_member` / `remove_member` (`confirm=true`; optional
  `project` for a project guest); `rename_project` (URL unchanged) /
  `transfer_project` (`confirm=true`, `to` from list_accounts);
  `get_build_settings` / `set_build_settings` / `clear_build_settings`;
  `get_notification_settings` / `set_notification_settings` /
  `list_notifications`
- Webhooks: `nesly webhooks` — MCP `list_webhook_endpoints` /
  `add_webhook_endpoint` (secret once) / `update_webhook_endpoint` /
  `roll_webhook_secret` (secret once) / `remove_webhook_endpoint`
  (`confirm=true`) / `test_webhook_endpoint` /
  `list_webhook_attempts`
- Log drains: `nesly logs drains` — MCP `list_log_drains` /
  `add_log_drain` (secret once) / `update_log_drain` /
  `roll_log_drain_secret` (secret once) / `remove_log_drain`
  (`confirm=true`) / `test_log_drain`
- Access tokens: `nesly token` — MCP `list_access_tokens` /
  `create_access_token` (secret once) / `revoke_access_token`
  (`confirm=true`). Prefer deploy hooks for rebuild-only CI.
- Activity: `nesly activity` — MCP `list_activity` (who changed what,
  including refusals; optional project / since / limit)
- Account export: `nesly account export -o <path>` — MCP
  `export_account_data` (absolute `output` path; archive never in the
  tool answer; database contents not included)
- Close account: `nesly account close --confirm <email> --yes` — MCP
  `close_account` (`confirm_email` from whoami + `confirm=true`);
  irreversible — export first
- Sessions: MCP `list_sessions` / `logout_all_sessions`
  (`confirm=true`; ends this session too — ask first)
- Processes: `nesly ps`, `nesly ps restart <name>`, `nesly ps scale name=N` —
  MCP `list_processes` / `restart_process` / `scale_process` (web is not
  restarted or scaled here — deploy or roll back)
- One-off command: `nesly run -- <cmd>` — MCP `run_command` (already-split
  words; no TTY; returns exit_code + output; non-zero is isError)
- Scheduled jobs: `nesly cron add|ls|pause|rm|run|runs|logs` — MCP
  `add_cron_job` / `list_cron_jobs` / `pause_cron_job` /
  `update_cron_job` (schedule/command/time_zone/timeout; omit
  unchanged fields) / `remove_cron_job` (`confirm=true`) /
  `run_cron_job` / `get_cron_runs` / `get_cron_run_logs` (optional
  `run` from get_cron_runs; empty = most recent)
- Destructive commands (`addon rm`, `backup rm`, project deletion)
  permanently delete data — confirm with the user first.

Full agent contract: https://nesly.net/docs/agents/ and
https://nesly.net/llms.txt
