#!/usr/bin/env bash
# pack-measure.sh — emit the measurement block for an evidence pack at a frozen tip.
#
# Why this exists: four separate defects in one day came from a human typing a
# number about a file into that same file. Writing the claim changes the thing
# the claim measures, so the artifact can never be self-consistent — the fixed
# point does not exist. This script measures from OUTSIDE, at a named SHA, and
# stamps every figure with its unit. Paste its output; never retype it.
#
# Usage:
#   scripts/harness/pack-measure.sh <sha> <pack-dir> [code-base-sha]
#
# Every number it prints carries (a) the SHA it was measured at and (b) its unit.
# A count without those two rots the moment it is quoted somewhere else.

set -euo pipefail

SHA="${1:?usage: pack-measure.sh <sha> <pack-dir> [code-base-sha]}"
PACK="${2:?usage: pack-measure.sh <sha> <pack-dir> [code-base-sha]}"
BASE="${3:-}"

PACK="${PACK%/}"
FULL="$(git rev-parse --verify "${SHA}^{commit}")"

echo "# MEDIÇÃO DO PACK — gerada, não digitada"
echo
echo "tip            = ${FULL}"
echo "tip (curto)    = $(git rev-parse --short "$FULL")"
echo "pack           = ${PACK}"
echo "gerado por     = scripts/harness/pack-measure.sh"
echo

# --- files -------------------------------------------------------------------
# `git ls-tree` WITHOUT -r counts tree entries, not files: it answers a different
# question correctly, which is why a reader reports it in good faith.
echo "## Arquivos versionados no pack (git ls-tree -r, no tip)"
git ls-tree -r --name-only "$FULL" -- "${PACK}/" > /tmp/pm-files.$$ || true
echo "arquivos       = $(wc -l < /tmp/pm-files.$$)   [unidade: ARQUIVOS, recursivo]"
echo "entradas raiz  = $(git ls-tree "$FULL" -- "${PACK}/" | wc -l)   [unidade: ENTRADAS DE ÁRVORE — não é contagem de arquivo]"
echo

# --- per-file size -----------------------------------------------------------
# bytes != characters on any file with non-ASCII prose, and LF-terminated lines
# != split('\n') parts. Both bit us; both are printed.
echo "## Tamanho por arquivo, no blob do tip"
printf '%-58s %10s %10s %8s\n' "arquivo" "bytes" "chars" "linhas"
while IFS= read -r f; do
  [ -n "$f" ] || continue
  blob="$(git rev-parse "${FULL}:${f}" 2>/dev/null)" || continue
  bytes="$(git cat-file -s "$blob")"
  # NOT `wc -m`: without a UTF-8 locale it silently returns BYTES, so the two
  # columns come out identical on accented prose and the reader sees agreement
  # where there is none. Deleting UTF-8 continuation bytes (0x80-0xBF) leaves
  # exactly one byte per character, locale-independently.
  chars="$(git cat-file -p "$blob" | tr -d '\200-\277' | wc -c)"
  lines="$(git cat-file -p "$blob" | grep -c '' || true)"
  printf '%-58s %10s %10s %8s\n' "$(basename "$f")" "$bytes" "$chars" "$lines"
done < /tmp/pm-files.$$
echo
echo "unidades: bytes = git cat-file -s · chars = wc -m (multibyte conta 1) · linhas = terminadas por LF"
rm -f /tmp/pm-files.$$

# --- custody -----------------------------------------------------------------
# Two questions, two instruments. `git status --porcelain` alone conflates them
# and false-positives with core.autocrlf=true (stat cache reports " M" for a
# byte-identical file).
echo
echo "## Custódia — duas perguntas, dois instrumentos"
if git diff --quiet HEAD -- "${PACK}/" 2>/dev/null; then
  echo "deriva de conteúdo  = NENHUMA   [git diff --quiet HEAD, exit 0]"
else
  echo "deriva de conteúdo  = PRESENTE  [git diff --quiet HEAD, exit != 0] — confira com git diff, não com status"
fi
untracked="$(git status --porcelain --untracked-files=all -- "${PACK}/" | grep -c '^??' || true)"
echo "não versionados     = ${untracked}   [linhas '??' de git status --porcelain -uall]"
echo
echo "NOTA: um pack não afirma a própria rastreabilidade — o arquivo que carrega a frase"
echo "      é commitado depois da frase existir. Esta medição vale porque roda de FORA."

# --- code delta --------------------------------------------------------------
if [ -n "$BASE" ]; then
  BASE_FULL="$(git rev-parse --verify "${BASE}^{commit}")"
  echo
  echo "## Delta de código ${BASE}..${SHA} — pergunta SEM filtro de path"
  total="$(git diff --name-only "$BASE_FULL" "$FULL" | wc -l)"
  outside="$(git diff --name-only "$BASE_FULL" "$FULL" | grep -vc '^\.mnfs/' || true)"
  echo "arquivos mudados        = ${total}"
  echo "fora de .mnfs/          = ${outside}   [0 = a árvore de código é idêntica EM TODO LUGAR, não só em apps/]"
  if [ "$outside" -gt 0 ]; then
    echo
    echo "arquivos de código tocados:"
    git diff --name-only "$BASE_FULL" "$FULL" | grep -v '^\.mnfs/' | sed 's/^/  /'
    echo
    echo "shortstat (só código):"
    git diff --shortstat "$BASE_FULL" "$FULL" -- $(git diff --name-only "$BASE_FULL" "$FULL" | grep -v '^\.mnfs/' | tr '\n' ' ') | sed 's/^/  /'
  fi
fi
