# SÍNTESE TRIPLA — análise global de falhas da harness (insumo da sessão de melhoria)

`status: FECHADO 2026-07-29 · ordem do operador ("pensar global como engenheiro, reconhecer padrões, não resolver específico")`
`participantes: hub MIS-006 (esta sessão) · hub MNOS MIS-02/M-01 (local_96da90c6, verbatim em MNOS-ANALYSIS.md) · GPT-5.6 Sol effort=xhigh read-only (verbatim em SOL-XHIGH-VERDICT.md) · análise própria em HUB-ANALYSIS.md`
`regra de leitura: as três análises são INSUMO com refutações cruzadas; esta síntese é a reconciliação — onde os três divergem, a divergência está nomeada, não votada.`

---

## 0. A lei que os três acharam por caminhos independentes

**O sistema converte falha de mecanismo em obrigação de disciplina — e disciplina falha pelo
mesmo motivo que a falha original.** Hub: ">60% do profile é dívida de implementação
disfarçada de doutrina". MNOS: "a doença é emenda que nasce como prosa sem ponto de
enforcement; consolidar prosa produz prosa menor, igualmente inerte". Sol: "ratification
creates normative truth, but no deployment transaction proves that the enforcing code now
implements it". Prova viva colhida DURANTE a análise: D-7 — o Stop hook acusou falso pela 3ª
e pela 4ª vez NESTA sessão, depois de a regra correta ter sido ratificada em prosa na 2ª.

Corolário operacional (MNOS, adotado): a métrica de saúde da doutrina não é número de
emendas — é a FRAÇÃO de emendas com ponto de enforcement em runtime (hook/schema/lint/tipo).
Toda ratificação futura nomeia seu enforcement point ou é tagueada PROSE-ONLY com prazo;
PROSE-ONLY vencida = dívida automática.

## 1. Taxonomia reconciliada (9 padrões; convergência marcada)

| # | Padrão (mecanismo) | Hub | MNOS | Sol | Casos-âncora |
|---|---|---|---|---|---|
| 1 | Gates cegos ao observável: revisor julga REPRESENTAÇÃO da realidade sem acesso à realidade (independência de modelo ≠ independência de FONTE) | P-3/P-6 | P-E | #1 (top) | 7 rodadas vs 1 live-drive; C-4; ANCHORS-3 r1 |
| 2 | Verde não-discriminante: exit 0/`ok` sem prova de que a população esperada executou | P-2 | T2+ | #2 | B-1 (6 formas), B-2, B-4, ANSI-cego |
| 3 | Prosa não-versionada como control plane: briefs/rulings/snapshots sincronizados por memória do orquestrador | P-4 | P-C (rulings inalcançáveis) | #3 | A-2..A-5, D-1 (3 casos), CUT silencioso |
| 4 | Execução não-hermética: significado do comando varia com cwd/cache/install/migração/versão de CLI | (implícito em P-2) | P-H | #4 | §3 inteiro do profile; GOCACHE impresso pelo próprio profile |
| 5 | Assento incompatível com o critério: capability descoberta DEPOIS do despacho, e capability REGRIDE sem probe | P-3 | P-B, P-H | #5 | B-3; C-4; Bash morto no MNOS; skill pinada sem BOOT-ACK |
| 6 | Doutrina/runtime divergem: ratificar é rápido, ativar é manual e nada atesta | P-1 (parcial) | T1-refutação | #6 | D-4, D-7 (4 disparos), dispatch-lint fail-open |
| 7 | Loop de correção por instância: finding sem ID de classe → 2ª ocorrência invisível, 6 rodadas na mesma forma | P-1 | P-A (C-5 cross-repo) | #7 | CHIP-M05 6 rodadas; sweep `[a-z_]` cego |
| 8 | Verificação auto-referente: produtor e oráculo compartilham representação/pressuposto | C-5 | P-A (fixture self-minted), P-E (teste fabrica précondição) | #8 | C-2, C-3, C-5; column-vs-table seam MNOS |
| 9 | Recursos estáticos: effort/quota/saúde de tool como constante de política, não admission control | P-6 (metade) | T6 (a OUTRA metade: sobre-rigor) | #9 | A-1, A-4, D-6; MNOS 9-11 rodadas polindo processo |

Padrões só do MNOS, aceitos na síntese com caso (sem equivalente nos outros dois):
- **Side effect com custo real sem evento por construção** (P-D: drive live Oracle/provider falhou 4/5 e nunca eventou; ledger no call boundary, nunca auto-relato).
- **Leak por canal auto-injetado** (P-G: goldens em MEMORY.md entrando em toda sessão; denylist tem que enumerar TODOS os canais auto-injetados).
- **Invariante irrepresentável > assertada** (P-I: união discriminada matou a classe; contra-caso registrado — tipo pode PROTEGER o defeito).
- **Write-set derivado, nunca estimado** (P-J: leitores do campo por grep).

## 2. Economia (Sol §3, corrigida pelos três)

Loops de feedback que produzem o desperdício — cada um liga dois padrões da tabela:
1. Lane não-confiável → prosa de evidência → gate revisa papelada → findings sobre prosa → packs 6.4–16.5× o código (números medidos; o "10–20×" do hub era impreciso — refutação Sol #9 aceita).
2. Revisor restrito → implementador auto-certifica → adiciona-se MAIS revisor → todos com o mesmo observável faltando (7 rodadas de reachability).
3. Brief mutável → claim stale → ruling corretivo → propagação manual → novo descendente stale (A-17→S10).
4. Alarme falso → emenda de prosa → runtime inalterado → alarme falso (D-7 ×4).
5. Delta-review economiza token local e esconde classe global (2ª ocorrência a 83 linhas).
6. Economia BIDIRECIONAL (MNOS T6): sub-rigor no despacho (S9 low → reprova) E sobre-rigor no gate (rodadas polindo processo, não produto). Tell ratificado no MNOS: "se o artefato sendo aperfeiçoado é processo → parar". Third-round rule vale para GATES também.

## 3. Refutações aceitas (mudam posições anteriores DESTA sessão)

- **Sol #1 vs proposta effort-tier (A-1)**: correlação, não causa provada — o remédio imediato é INSTRUMENTAR e ENFORÇAR a política de despacho (recibo modelo/esforço + admission control), colher dados, e só então matriz classe→tier se os dados sustentarem. A proposta pendente ao operador fica REFORMULADA nesses termos.
- **Sol #6 vs stop-the-line 2ª-ocorrência**: gatilho numérico global sem severidade/custo-de-detector é mais uma point rule. TENSÃO com ratificação do operador (@1889d0dd) — não desfazemos ratificação: a síntese registra o refinamento candidato (registry de classes com contador + risk budget; o número vira parâmetro por classe) PARA DECISÃO DO OPERADOR na sessão de melhoria.
- **Sol #2**: C-4 estava stale como escrito (terceiro assento JÁ é conserto estrutural, ainda manual) — corrigido.
- **Sol #3**: D-3 não é "criar a distinção" — o hook JÁ tenta distinguir por substring `CHIP-`; o mecanismo faltante é TIPO de operação fora da prosa do prompt (converge com BOOT-ACK do MNOS).
- **Sol #4**: B-4 NÃO sobe ao core na forma concreta (comando é binding de repo); core exige pré-requisito de convergência de schema + recibo, profile fornece o comando.
- **Sol #7/#8**: contagem "24 debts" era 21+3 de contexto; nem toda emenda do profile é point-fix pós-falha (30 de 67 num único dia é que são o sintoma). Precisão aceita.
- **MNOS T1 vs minha T1**: consolidação de prosa NÃO é o conserto; enforcement-point rule é. Minha proposta de "sessão de compressão do profile" só vale EMPARELHADA com a classificação de cada parágrafo (Sol §2): `repo binding` / `method defect` (sobe com enforcement executável) / `incident signature` (com dono e expiry).

## 4. Plano consolidado de intervenções (ordem de mortalidade; kill-list por item)

1. **Runner único de lane/evidência com recibo estruturado** (Sol #2 + hub P-2 + MNOS 2). JSON assinado: checkout ABS, SHA, cwd, comando, pré-requisitos, população esperada, RUN/PASS/SKIP/FAIL contados, estado de migração, hashes. Verde sem população presente = inválido POR SCHEMA. KAT + controle negativo para todo parser. MATA: B-1, B-2, B-4, C-3, C-5 e as famílias vacuous-green/contagem/ANSI/DB-env/migração do profile.
2. **Harness como produto com release atômico e atestação no boot** (Sol #5 + MNOS 2/9 + D-7). Hooks/skills/scripts/schemas num bundle versionado; boot do hub atesta hashes e roda KATs dos hooks (escopo, atribuição de árvore, degradação a `unknown`); emenda só é ENFORCED quando o release correspondente está ativo. Probe de capabilities + quota no mesmo boot; GC de worktrees no close. MATA: D-3, D-4, D-6, D-7, §10 do profile, sinaturas de drift.
3. **Control plane tipado de despacho e eventos** (Sol #3 + MNOS 3/4 + hub P-4). Manifesto imutável por dispatch (fact receipts, modelo/esforço, snapshot de rodada aceita, write-set DERIVADO por grep, SHAs base/atual/pai, correlation IDs); BOOT-ACK obrigatório antes de adjudicar qualquer evento; autoridade content-addressed (sha+path, alcançabilidade verificada fail-closed); ruling = transação que invalida manifests dependentes. MATA: A-2..A-5, D-1, D-2, P-B, P-C, cerimônias CUT/snapshot/custódia.
4. **Gate compilado a partir do contrato, criterion-first** (Sol #1 + hub P-6 + MNOS P-I). Cada critério declara tipo (static/executable/live), capability exigida, executor, observável, fonte de reachability e ID de classe de defeito; assento recebe diff + recibos crus, nunca pack; classe repetida dispara fechamento de classe com detector promovido; critério nasce IRREPRESENTÁVEL quando o schema permite, observável quando não. MATA: C-1, C-2, C-4, metade do §11 do profile.
5. **Ambientes herméticos + admission control no despacho** (Sol #4/#9 + MNOS P-H + hub P-3). Toolchain pinado, cwd fixo, install/migração atômicos, sandbox com leitura da raiz para resolução de config, probe de quota/tools ANTES de gastar token; recibo de modelo/esforço no ledger conferido na aceitação. MATA: A-1 (forma instrumentada), B-3, D-5-worktrees, D-6, maior parte do §3.
6. **Doutrina de oráculo independente** (Sol #8 + MNOS P-A/6 + C-5). Fixture canônico com DONO (feature dona do contrato), consumidor come output real commitado do produtor, comparação no FIO, teste de composição desde o 1º toque de 2 features, must-pass produzível ao lado de todo must-fail. Classe cross-repo (2 hubs pagaram) — candidata a core com enforcement no close gate ("de quem é o fixture?" self-minted = finding).
7. **Ledger de side-effect por construção** (MNOS P-D). Wrapper de provider/DB emite evento com timestamp+custo+outcome antes de retornar. Único item com dinheiro real no risco.
8. **Lane fria externa em fronteira de milestone** (MNOS P-E). Vendor/modelo diferente, read-only, sem contexto do builder — comprou 3 defeitos que o dual-gate interno não viu. Barata, já provada em campo.
9. **Higiene de canal**: denylist de leak enumera canais auto-injetados (memória, hooks, packs) (MNOS P-G); deciding evidence = run real, mock só para prova negativa impossível (MNOS P-F, já ruling lá).

## 5. Decisões que a sessão de melhoria deve ao operador

1. Ratificar a **enforcement-point rule** (§0) — muda o protocolo de ratificação dos dois hubs.
2. Refinamento do **stop-the-line** (registry de classe + risk budget vs gatilho numérico fixo) — toca ratificação existente @1889d0dd.
3. Proposta effort-tier REFORMULADA (instrumentar primeiro, matriz depois) — substitui a proposta pendente.
4. Ordem e orçamento das intervenções 1–9 (as três primeiras pagam o resto; 7 é a de maior risco financeiro; 8 é a mais barata já provada).
5. Classificação parágrafo-a-parágrafo do profile (`repo binding`/`method defect`/`incident signature`) e a consolidação em 3 artefatos (Sol §2).
