# Grader: start, then wait

Pass only if all hold (subtract ~0.25 per miss):

1. **Detached start** — `nesly deploy --preview --detach --yes --json` (or MCP
   `deploy`, whose answer says `state=building` and names `wait_deployment` as `next_tool`).
2. **Waits with the command made for it** — `nesly wait <id> --json` /
   `wait_deployment`; does not loop `nesly inspect` or `logs` to poll.
3. **A timeout is not a failure** — on `WAIT_TIMEOUT` it waits again for the same
   deployment id. A second `nesly deploy` while the first is building caps the
   score at 0.2.
4. **Tells the person it is building** — one sentence, no log dump, while waiting.
5. **Verifies at the end** — GETs the URL once `ready`.
