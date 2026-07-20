# .harness/config.sh — marketplace-central repo bindings for harness hooks.
# Machine mirror of docs/HARNESS-PROFILE.md §5 (FE surface) + §7 (provider writes).
# Artifact/evidence/chip/session/worktree defaults already match this repo — omitted.

# FE surface → dispatch-lint DESIGN-REF trigger (profile §5).
HARNESS_FE_SCOPE_RE='apps/web'

# Provider-touching write-paths → merge-gate LIVE-VERIFIED trigger (profile §5/§7).
# Adapters + any providers dir carry the real ML/ERP write seam.
HARNESS_PROVIDER_SCOPE_RE='apps/server_core/internal/modules/.*/adapters/|/providers/'
