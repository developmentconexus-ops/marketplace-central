# CHIP-ANCHORS — rung de governance (hub-run, ruling R-b)

O chip pediu ruling porque C11 e o prompt de dispatch se contradizem: C11/profile §2 exige worktree
limpo detached, o prompt proíbe o chip de criar worktree. Ruling do hub: **R-b — a lane é seam do
hub** (mesma família da dev stack; `hub-ops` carrega "governance-lane runs in a clean worktree" no
charter). O chip deixou a rung OPEN e não segurou o CLOSED.

## Execução

Worktree limpo detached fora de `.claude/worktrees/` (o scan varre esse diretório e falsifica o
resultado quando roda do checkout do hub):

```
git worktree add --detach C:/Users/leandro.theodoro/Documents/mc-gov-anchors 8e37958
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command governance \
  -BaseSha 917f7bb58e385847fba5612201823f9db48791c6
```

Rodada duas vezes no MESMO worktree: uma no tip do chip `8e37958a` e outra no BASE-SHA
`917f7bb58e385847fba5612201823f9db48791c6`, para separar herança de regressão.

## Resultado

`status=failed`, exit 1, **53 violações** nas duas pontas:

| error_code | n |
|------------|---|
| `GOV_MODULE_COVERAGE` | 2 |
| `GOV_MODULE_DEPENDENCY` | 13 |
| `GOV_MODULE_LAYER` | 6 |
| `RCFG_READER_MISSING` | 10 |
| `RCFG_UNAPPROVED_READER` | 17 |
| `RCFG_UNDECLARED_READ` | 5 |

`Compare-Object` entre as duas saídas (175 linhas cada) devolveu **vazio** — as saídas são
idênticas linha a linha, incluindo as 14 `baseline_exception`.

## Veredicto

**C11 rung de governance: PASS na leitura diferencial.** O chip introduziu **zero** violação nova.
O `failed` é herdado e pré-existe ao BASE-SHA — nada nele nomeia arquivo escrito por este chip;
`GOV_MODULE_COVERAGE` acusa `sourcekind` e `tenant_config` (M-02, sem entrada em `modules.json`) e
`RCFG_*` acusa os adapters legados magalu/amazon/shopee mais `MC_ERP_SOURCE` em `root.go`.

Nota honesta: isto **não** é a lane verde. A lane está vermelha em `main` e continua vermelha
depois deste chip. Aceitar o chip sobre uma lane vermelha só é legítimo porque a diferença é nula e
está provada aqui; o débito de governance continua aberto e é do hub, não deste chip.

Saídas cruas: `scratchpad/gov-tip.txt`, `scratchpad/gov-base.txt` (sessão do hub).
