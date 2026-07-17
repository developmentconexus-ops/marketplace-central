# P1 Clarified Decisions — MIS-004 / MIS-005 replan

> Registro do gate P1 (P1a domain scan · P1b architecture clarify · P1c quality scan), operador respondeu 2026-07-17 via AskUserQuestion. Fonte de intake: `.mnfs/REPLAN-BRIEF-2026-07-17.md` + pacote design + research de pricing.

## P1a — Domain scope (capabilities)

| Item | Decisão | Nota |
|---|---|---|
| Brief §4 itens 1–9 (identidade, modelo preço, adapter ML, shell, Vínculos, Anúncios, Produto Detalhe parcial, Simulador, Pedidos) | IN (lean-core confirmado) | corte do brief ratificado |
| Dashboard | IN, cortável | última prioridade, agregações de dados já trazidos, corta se apertar |
| DIFAL seed mínimo | IN | tabela 27 UFs seed "padrão 2026" + toggle no drawer do Simulador + chip/coluna em Pedidos; SEM agendamento/lembrete/marcar-pago (MIS-005) |
| Webhooks | OUT do MVP | polling/GET cobre demo read-only; webhooks completos = MIS-005 |
| **Import .xlsx de produtos ERP** | **IN (novo, obrigatório p/ demo)** | empresa da apresentação NÃO estará conectada ao Sankhya; demo usa planilha. Operador delegou o contrato de colunas ao planning ("apenas o que você julgar necessário"). Sankhya/Oracle adapter existente permanece como caminho alternativo. |

## P1b — Architecture clarifications

| Taxon | Decisão do operador |
|---|---|
| Persistência / fonte ERP demo | Oracle/Sankhya continua + criar import .xlsx; demo de segunda usa .xlsx |
| UI convergence / sequência FE | Retheme/shell PRIMEIRO milestone da trilha FE; telas novas herdam; roda paralelo ao backend |
| Writes ML na demo | ZERO writes live; fila de sync mostra preview+protocolo com execução contra ML desabilitada; write real gated pós-demo |
| Validação / gate de close | Doutrina completa por milestone (dual gate Opus+Sol + QA live-drive); sem compressão |
| Runtime da demo | Docker dev stack local (npm run docker:dev; :8080/:5174) |

## P1c — Quality attributes

| Attr | Decisão | Alvo |
|---|---|---|
| Observabilidade de coleta | TARGET | evidência de coleta visível em TODA UI de preço (fonte, fetched_at, idade, n_offers/n_sellers) + telemetria na rota flag-gated products/{id}/items — critérios de validação obrigatórios |
| Durabilidade ADR-17 | TARGET | snapshot válido nunca sobrescrito por zero/falha; teste negativo no contrato de validação |
| Security | baseline | sem superfície nova de auth no MVP; tenancy/secrets per profile §7; auth/multi-tenant real = MIS-005 |
| Maintainability | baseline | L0-L2 per profile |
| Performance | DECLINED | volume single-seller, demo local — sem p95 explícito |
| A11y | DECLINED | herda padrões do design system; sem critério WCAG neste prazo |

## Accepted assumptions (forced-assumption ledger)

- Conta ML conectada = a existente do operador/tenant demo (probes do research já a usaram); pré-anúncio p/ produtos do cliente via catálogo oficial por EAN. Reversível (trocar installation).
- Coleta de sinais/snapshots: on-demand na preparação da demo (sem scheduler; scheduler diário = MIS-005). Reversível.
- Flag `products/{id}/items` default OFF; ligar na demo = decisão do operador na hora. Reversível.
- Contrato de colunas xlsx: definido pelo planning como IC (obrigatórias: CODPROD, DESCRPROD, CUSTO, ESTOQUE_FISICO; opcionais: ESTOQUE_RESERVADO, EAN, REFFORN, MARCA, NCM; opcional ausente → estados honestos: sem EAN ⇒ matching só REVIEW/sem auto-ACCEPT por título). Reversível até P4 (IC ratificado).
- Nav pills Mercado/Repasses no MVP: visíveis desabilitadas "em breve" (não ocultas) — evita nav mentirosa e mostra roadmap na demo. Reversível (CSS/flag).
- Dashboard do mock é rascunho multi-canal desatualizado → recriar no domínio MPC (ML-only). Irreversibilidade baixa.
- MIS-005 planejado em grão de milestone agora; P4–P7 completos antes da execução dele.
