#!/usr/bin/env bash
# Architecture gate. Run from the repository root.
#
# Exit 0 means: formatted, vetted, tested, and zero architecture findings.
# Every check prints what it measured, because a gate that only prints "ok" is
# indistinguishable from a gate that skipped.
set -euo pipefail

cd "$(dirname "$0")/.."
SERVER="apps/server_core"
export GOCACHE="$PWD/$SERVER/.gocache"

fail=0
step() { printf '\n== %s ==\n' "$1"; }

step "gofmt"
unformatted="$(gofmt -l "$SERVER/internal" || true)"
if [ -n "$unformatted" ]; then
  echo "$unformatted"
  echo "gofmt: files above are not formatted"
  fail=1
else
  echo "gofmt: clean"
fi

step "go vet"
if ! (cd "$SERVER" && go vet ./...); then fail=1; fi

step "architecture detectors"
for root in internal/kernel internal/contexts internal/adapters internal/composition; do
  if (cd "$SERVER" && go run ./internal/arch/cmd/archscan -root "$root"); then
    echo "archscan $root: zero findings"
  else
    fail=1
  fi
done
echo "note: internal/modules is the legacy tree and is deliberately NOT scanned"

step "unit tests"
if ! (cd "$SERVER" && go test ./internal/...); then fail=1; fi

step "working tree"
dirty="$(git status --porcelain --untracked-files=all)"
if [ -n "$dirty" ]; then
  echo "$dirty"
  echo "tree: dirty — a gate cannot certify a tree it did not see"
  fail=1
else
  echo "tree: clean"
fi

printf '\n'
if [ "$fail" -ne 0 ]; then
  echo "ARCH GATE: FAIL"
  exit 1
fi
echo "ARCH GATE: PASS"
