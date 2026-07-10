#!/usr/bin/env bash
set -euo pipefail

if [[ ! -d node_modules/.vite ]]; then
  npm install
fi

exec npm run dev --workspace @marketplace-central/web -- --host 0.0.0.0 --port 5174
