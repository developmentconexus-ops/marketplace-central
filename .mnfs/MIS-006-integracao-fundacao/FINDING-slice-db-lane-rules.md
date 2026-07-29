# FINDING — lane de DB só prova com env carregado E base migrada; o brief é quem parte isso

`status: profile-candidate · vinculante para CHIP-VENDAVEL desde já (autoridade do hub sobre protocolo de chip); entrada no HARNESS-PROFILE pende ratificação do operador`
`provenance: 2026-07-29 · CHIP-VENDAVEL S5B · worker lane RUN 27 / PASS 1 / SKIP 26`

## O facto

O brief da S5B omitiu a linha `. testdb-env.ps1` que S1..S5 carregavam. Sem
`MPC_TEST_DATABASE_URL`, todo teste de DB cai em `SkipWithoutTarget` e o pacote imprime `ok`
com exit 0 — assinatura 3 do §11 do perfil, à letra. A run do worker foi RUN 27 / PASS 1 /
SKIP 26, e o teste de round-trip — a razão da slice existir — era um dos skips.

O vácuo só ficou visível porque o worker contou os SKIPs e escreveu na evidência *"this is not
a claim of live integration GREEN"*. O instrumento é do perfil; a mão foi do worker; funcionou
contra o resultado do próprio worker.

Segunda camada, descoberta na re-medição do chip: a migração 0085 existia na árvore e **nunca
tinha sido aplicada à base da slice** — 13 falhas `column "usoprod" does not exist (42703)`.
`OpenPool` valida e liga; não migra. Migração nova é invisível até alguém correr `cmd/migrate`.
Ruidoso (42703), não silencioso — mas custa uma rodada inteira.

## Regras que saem (as duas do chip, ratificadas pelo hub para este chip)

1. **Brief que toca DB carrega a linha de dot-source do env** (`. testdb-env.ps1` ou
   equivalente do repo) **E exige o relatório de SKIP explícito** (RUN/PASS/SKIP/FAIL contados
   por linha). A segunda metade salvou a S5B; a primeira é o que a partiu.
2. **Depois de uma slice adicionar migração, o orquestrador migra a base da slice antes de
   rever** (`go run ./cmd/migrate` no worktree). `applied N` também serve de prova de não-drift:
   N deve ser exatamente o nº de migrações novas da slice.

## Corolário de forma (do mesmo commit)

`%v` sobre `*string` imprime ENDEREÇO — asserção que discrimina o valor mas reporta ponteiro
(`got 0x34f73c6925b0, want R`) recusa-se a dizer o que viu; quem lê não distingue valor errado
de nil. Mensagem de teste deref'a antes de imprimir, e o conserto da MENSAGEM se prova por
re-injeção como qualquer conserto de teste (regra do S4, um grau mais suave que o panic).
