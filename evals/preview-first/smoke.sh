#!/usr/bin/env bash
# preview-first smoke — d390/G9 exit proof without an LLM.
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/../../../.." && pwd)

bash "$ROOT/tools/discipline/check-agent-preview-first.sh"
bash "$ROOT/tools/discipline/check-agent-preview-first.sh" --self-test

echo "preview-first smoke: PASS"
exit 0
