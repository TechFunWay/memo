#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_PORT="${PORT:-8905}"

cd "${SCRIPT_DIR}"
exec ./memo \
    -port "${APP_PORT}" \
    -web-dir "${SCRIPT_DIR}/www" \
    -data-dir "${SCRIPT_DIR}/data" \
    "$@"
