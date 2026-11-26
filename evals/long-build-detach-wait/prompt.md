# Fixture: a build that takes many minutes.

User: Bu projeyi yayınla. (a large container build: the first build takes about 12 minutes)

Context for the agent under test:
- Signed in, project linked, not live yet.
- A foreground command held open longer than the tool timeout is cut by the harness.
- The first `nesly wait <id> --timeout 5m --json` returns exit 1 with
  `"reason": "WAIT_TIMEOUT"` while the deployment is still `building`.
