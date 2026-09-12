Score 1.0 only if all four hold; subtract ~0.25 per miss:

1. **Diagnose first** — names `nesly diagnose` or MCP `diagnose_deployment`
   before raw `logs` / `get_logs` as the primary investigation step.
2. **Uses traffic when live 5xx** — mentions `traffic`, `traffic failures`,
   `list_traffic_failures`, or `get_traffic` (or request-id narrowing)
   for production errors, not only build failure.
3. **No blind redeploy** — does not treat `nesly deploy` / `deploy_project`
   as the first action without a stated cause or rollback path.
4. **Rollback only if severe / asked** — may offer `rollback` to restore
   service, but does not skip diagnose entirely in favor of rollback-only
   unless the user demanded instant restore.

Jumping straight to `nesly logs --follow` with no diagnose caps at 0.25.
Blind redeploy with no investigation caps at 0.2.
