Score 1.0 only if all four hold when the scenario requires restoring
service (or the agent chooses rollback); subtract ~0.25 per miss:

1. **Lists revisions first** — `nesly ls` / `list_deployments` / inspect
   before calling `rollback`.
2. **Targets a previous id** — rolls back to a prior `dep_…`, not the
   current live bad deployment, and does not invent ids.
3. **States no rebuild** — says rollback repoints traffic / is instant /
   does not rebuild.
4. **Verifies after** — `project info`, `diagnose`, or a live URL check
   after rollback.

If the agent correctly stays on diagnose→fix without rollback because
errors are not framed as an outage, score N/A for this grader (do not
penalize). Blind `rollback` with no list step caps at 0.25.
