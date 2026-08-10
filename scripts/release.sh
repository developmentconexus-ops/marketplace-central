#!/usr/bin/env bash
set -euo pipefail

# The release path: build both prod images, tag with the current commit sha,
# push to GHCR. There is no CI publisher -- the workflow that used to mirror
# this on a runner was deleted 2026-08-10 after failing all 18 of its runs with
# 403 Forbidden on the push (see the ADR-008 amendment).
#
# Prereq (once): docker login ghcr.io -u <github-user>   # PAT with write:packages
# Usage:         bash scripts/release.sh [--no-push]

OWNER="${MPC_GHCR_OWNER:-developmentconexus-ops}"
REPO_BASE="ghcr.io/${OWNER}/marketplace-central"
SHA="sha-$(git rev-parse --short=7 HEAD)"
PUSH=1
[[ "${1:-}" == "--no-push" ]] && PUSH=0

if [[ -n "$(git status --porcelain)" ]]; then
  echo "WARNING: working tree is dirty — image tag ${SHA} will not match the working tree exactly." >&2
fi

echo "==> building ${REPO_BASE}-server:${SHA}"
docker build -f docker/prod/backend.Dockerfile \
  -t "${REPO_BASE}-server:${SHA}" -t "${REPO_BASE}-server:latest" .

echo "==> building ${REPO_BASE}-web:${SHA}"
docker build -f docker/prod/frontend.Dockerfile \
  -t "${REPO_BASE}-web:${SHA}" -t "${REPO_BASE}-web:latest" .

if [[ "$PUSH" == "1" ]]; then
  for image in "${REPO_BASE}-server" "${REPO_BASE}-web"; do
    echo "==> pushing ${image} (${SHA}, latest)"
    docker push "${image}:${SHA}"
    docker push "${image}:latest"
  done
  echo "==> release complete: ${SHA}"
  echo "    deploy: set MPC_IMAGE_TAG=${SHA} in the host .env, then compose pull && up -d"
else
  echo "==> build-only complete (${SHA}); skipped push"
fi
