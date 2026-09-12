#!/usr/bin/env bash
# debug-rollback smoke — G5 exit proof without pretending an LLM ran.
#
# 1) Deterministic contract gate (must be green).
# 2) Print the sentence shapes an agent must use for debug vs rollback.
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/../../../.." && pwd)

bash "$ROOT/tools/discipline/check-agent-debug-skills.sh"
bash "$ROOT/tools/discipline/check-agent-debug-skills.sh" --self-test

echo "--- agent contract shapes ---"
cat <<'EOF'
DEBUG (production 500):
  1) nesly diagnose --json  (or diagnose_deployment)
  2) nesly traffic failures --since 1h --json
  3) fix → redeploy as new revision; rollback skill if outage first

ROLLBACK (restore service):
  1) nesly ls --json
  2) nesly rollback <previous-dep-id> --json
  3) verify with project info / diagnose
EOF

echo "debug-rollback smoke: PASS"
exit 0
