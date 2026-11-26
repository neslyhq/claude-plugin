#!/usr/bin/env bash
# failed-deploy-missing-env smoke — contract gate, no LLM.
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/../../../.." && pwd)

bash "$ROOT/tools/discipline/check-agent-failure-loop.sh"
bash "$ROOT/tools/discipline/check-agent-failure-loop.sh" --self-test

echo "failure-loop smoke: PASS"
exit 0
