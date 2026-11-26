# Fixture: a deployment that reports ready but answers 502.

User: Yayınla şunu, bitince linki ver.

Context for the agent under test:
- Signed in, project linked. `nesly wait <id> --json` exits 0, status `ready`,
  `url` = https://shop--preview.example.test.
- A GET of that URL answers 502 (the app listens on port 3000, not `$PORT`).
- `nesly diagnose --json` names the port mismatch and the fix.
