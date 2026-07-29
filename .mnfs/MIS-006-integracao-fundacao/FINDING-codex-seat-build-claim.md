# FINDING — assento codex reporta `go build` falso, e o custo é o slice retido

`status: profile-candidate · 3 ocorrências idênticas, todas refutadas por medição do chip`
`provenance: 2026-07-29 · CHIP-VENDAVEL, slices S1 / S2 / S2B`
`locus: docs/HARNESS-PROFILE.md §3 (lanes) — candidato a nota de ambiente`

## O facto

Nas três slices, o worker codex reportou `go build ./...` a falhar — na terceira com
`error obtaining VCS status: exit status 128`. Nas três, o chip mediu por conta própria no MESMO
worktree, com `cd apps/server_core` e GOCACHE/GOMODCACHE absolutos: **EXIT=0**.

É assinatura do assento, não facto do repo. O `git` do assento do worker não consegue carimbar os
metadados VCS que o `go build` pede; o mesmo comando no assento do chip passa.

## O custo, que não é o alarme

O alarme é barato — mede-se em segundos. Caro foi o que o worker FEZ com ele: no S2B **reteve o
commit** de um slice verde. Trabalho terminado e correto ficou por entregar porque o worker tratou
uma observação de ambiente como veredito de lane.

## A regra que daí sai

1. Relato de build vindo do worker é **observação**, nunca gate. Só a medição do chip, no worktree
   do chip, descarrega ou reprova o L0. (`go vet` verde nunca serviu de prova de build — *vet ≠ ran*,
   M-03; isto é o mesmo erro pelo outro lado: build vermelho no assento errado não é prova de nada.)
2. O worker **não retém commit** por falha de build que não conseguiu atribuir. Reporta e entrega.
3. `-buildvcs=false` **não entra em lane**. Pode correr-se uma vez, lado a lado, para CLASSIFICAR a
   natureza da alegação; enterrá-lo na lane troca um alarme falso por cegueira permanente.

## Nota de instrumento — varredura por string em worktree

O chip observou que a ferramenta `Grep` enraizada no caminho absoluto do worktree devolveu 5 de 24
ocorrências (falhou o código), enquanto `git grep` de dentro do worktree devolveu todas. Refutou a
hipótese do `.gitignore` com `rg` de dentro, e registou o resto como constrangimento sem causa
atribuída — correto, não afirmou o que não mediu.

**Contra-medição do hub (2026-07-29):** a mesma `Grep`, com o caminho absoluto do worktree do chip,
a partir do assento do HUB, devolve **26 ocorrências em 8 ficheiros**, incluindo
`apps/server_core/internal/modules/tenant_config/repository.go` e as migrações. Ou seja: o
constrangimento **não reproduz do assento do hub** e é específico do assento (provável interacção
entre o cwd do chamador e as regras de ignore herdadas do repo pai, já que o worktree vive sob
`.claude/` do próprio repo — hipótese NÃO medida, fica como hipótese).

Consequência operacional que fica: quem varre por string **dentro de um worktree** usa `git grep` de
dentro do worktree. Varredura parcial é pior que nenhuma, porque tem a forma de completa — um
assento de gate a usar a ferramenta errada devolveria "zero ocorrências" com toda a confiança, e é
por varredura de string que este repo verifica alegações totais (R-24/R-26).
