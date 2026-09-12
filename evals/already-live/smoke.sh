#!/usr/bin/env bash
# already-live smoke — G4 exit proof without pretending an LLM ran.
#
# 1) Deterministic contract gate (must be green).
# 2) If `nesly` is signed in, print one live project URL and the exact
#    sentence shape an agent must use before deploy_project.
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/../../../.." && pwd)

bash "$ROOT/tools/discipline/check-agent-already-live.sh"
bash "$ROOT/tools/discipline/check-agent-already-live.sh" --self-test

echo "--- live lookup (optional) ---"
if ! command -v nesly >/dev/null 2>&1; then
  echo "nesly not on PATH; contract gate only"
  echo "already-live smoke: PASS"
  exit 0
fi

json=$(nesly project list --json 2>/dev/null || true)
if [[ -z "$json" ]]; then
  echo "not signed in or empty output; contract gate only"
  echo "already-live smoke: PASS"
  exit 0
fi

/usr/bin/python3 - "$json" <<'PY'
import json, sys
raw = sys.argv[1]
try:
    data = json.loads(raw)
except json.JSONDecodeError:
    print("project list not JSON; skip live URL")
    raise SystemExit(0)
items = data.get("projects") if isinstance(data, dict) else data
if not isinstance(items, list) or not items:
    print("no projects in account; contract gate only")
    raise SystemExit(0)
p = items[0]
name = p.get("name") or p.get("Name") or "?"
url = (
    p.get("default_url")
    or p.get("defaultUrl")
    or p.get("url")
    or ""
)
print(f"example project: {name}")
print(f"example live URL: {url or '(none on first row)'}")
print()
print("AGENT MUST SAY (before deploy):")
if url:
    print(
        f"  Already live at {url}. "
        "Next deploy publishes a new revision of the same project."
    )
else:
    print(
        "  (no URL on first row — get_project for a live one "
        "before claiming already-live)"
    )
PY

echo "already-live smoke: PASS"
exit 0
