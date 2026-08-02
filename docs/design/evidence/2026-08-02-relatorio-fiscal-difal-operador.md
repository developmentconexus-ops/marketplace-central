# Relatório fiscal — divergência de alíquota interna no DIFAL

Para: operador · Data: 2026-08-02 · **Fora do escopo do software.** Este documento existe para você levar ao contador. Nenhuma decisão fiscal foi tomada por nós e nenhum valor aqui é calculado pelo sistema.

Tudo abaixo foi **medido** na base `METALPRD` (acesso somente leitura) e nas tabelas de pedido do marketplace. Nada é estimativa.

---

## 1. O achado

Duas notas fiscais irmãs, mesmo produto (grupo fiscal 122), mesma origem (MG), mesmo destino (**Bahia**), mesmo tipo de operação, ambas para **pessoa física**, emitidas com **um dia de diferença** — usaram alíquotas internas de destino diferentes.

| nota | data | valor | alíquota interna BA usada | DIFAL |
|---|---|---|---|---|
| 895436 | véspera | R$ 189,90 | **20,5%** | R$ 25,64 |
| 895507 | dia seguinte | R$ 299,90 | **17,0%** | R$ 29,99 |

A fórmula reconcilia ao centavo nos dois casos (`valor × alíquota_destino − valor × 7%`), então não há erro de leitura nossa. Também foi verificado que não houve redução de base em nenhuma das duas (base do DIFAL = 100% do valor da nota).

## 2. Qual das duas está certa

**A de 20,5%.** A alíquota interna do ICMS na Bahia é **20,5% desde fevereiro de 2024**. A nota 895507 foi emitida com uma alíquota que estava desatualizada havia mais de dois anos.

## 3. Por que aconteceu

A matriz de ICMS do ERP (a tabela que decide a alíquota em cada operação) tinha a linha MG→BA do grupo 122 **congelada em 17,0%**. O log de alterações mostra que essa linha foi corrigida para 20,5% em **20/07/2026**, em seis edições sucessivas entre 14h02 e 14h42, feitas pelos usuários `LEANDROTH` e `SUP`.

Ou seja: **a Bahia mudou em fev/2024 e o cadastro só acompanhou em jul/2026 — 29 meses de atraso.** Notas emitidas nesse intervalo, para a Bahia nesse grupo, saíram com DIFAL recolhido a menor.

Para dimensionar: no mesmo log, o volume de alterações de alíquota vem acelerando — 3 em 2024, 34 em 2025, 62 até agosto de 2026. O atraso não é evento isolado.

## 4. O que perguntar ao contador

1. **Notas emitidas com DIFAL a menor entre fev/2024 e jul/2026 para a BA nesse grupo** — há obrigação de complementar o recolhimento? Qual o prazo e qual o instrumento (denúncia espontânea, retificação)?
2. **Outros estados na mesma situação.** Achamos a BA porque tínhamos duas notas próximas para comparar. Não fizemos varredura de todos os estados contra a legislação — não é competência nossa e não vamos fazer. Vale um levantamento pelo contador.
3. **Rotina de atualização.** Hoje a alíquota interna de destino é digitada à mão na matriz do ERP. Não existe alarme quando um estado muda. Vale definir quem acompanha e com que frequência.

## 5. Segunda lacuna, separada da primeira

Para **Paraná, Rio Grande do Sul e Santa Catarina**, a matriz do ERP diz que o DIFAL **incide**, mas **não informa a alíquota interna do estado de destino** — o campo está em branco. São 20 combinações (produto × estado) no recorte que realmente vendemos.

Isso é diferente do item 1: lá o número estava errado; aqui ele **não existe**. Curiosamente, o ERP **consegue** emitir nota com DIFAL para o Paraná (11 notas com 18,0%), então ele busca esse número em algum lugar que não é a matriz e que não conseguimos localizar na base.

**Decisão tomada:** o sistema **não vai adivinhar** esse número. Onde a matriz não responde, a tela mostra "desconhecido" com o motivo, e aponta a lacuna para ser preenchida no Sankhya. Impacto medido: 7 notas em 12 meses, de 74.

**Pergunta ao contador / à TI do Sankhya:** de onde a emissão da nota tira a alíquota interna do PR, já que a matriz está vazia? Se existir uma segunda tabela, ela precisa ser a fonte — ou a matriz precisa ser preenchida.

## 6. O que o software faz e o que não faz

**Faz:** lê a matriz de ICMS do ERP como está, calcula o imposto e o DIFAL antes da venda, e mostra na tela quando o ERP não sabe responder.

**Não faz e não vai fazer:** manter tabela própria de alíquotas legais, corrigir a matriz do ERP por conta própria, ou preencher lacuna com valor plausível. O ERP é o dono do dado fiscal. Onde ele estiver errado ou vazio, o sistema **mostra**, para que seja consertado onde deve — no Sankhya.

Esse é o mecanismo pelo qual esta divergência apareceu, aliás: ela não veio de auditoria fiscal, veio de o sistema tentar calcular a margem de um pedido e as contas não fecharem.

---

**Fontes das medições:** `TGFICM` (matriz vigente), `TGFHICM` (log de alterações), `TGFCAB`/`TGFITE` (notas 895436 e 895507), `TGFPRO.GRUPOICMS` (grupo fiscal do produto). Acesso somente leitura, nenhuma escrita.
