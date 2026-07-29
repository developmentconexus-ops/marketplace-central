# `TestCase3EANAloneYieldsMediaConfirm` — a asserção, antes e depois

Exigido pelo HUB RULING R2. A asserção foi mudada, o código não. Ela fica aqui verbatim para que
quem julgar não dependa do raciocínio do chip.

## ANTES — `git show 5441fe18:…/generation_service_test.go`, linhas 1441-1447

```go
	skuReason, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionIncomparable)
	if !ok {
		t.Fatalf("reasons=%#v, want seller_sku INCOMPARABLE", candidate.Reasons)
	}
	if !strings.HasPrefix(skuReason.Detail, "sem CODPROD para corroborar o EAN") {
		t.Fatalf("seller_sku reason detail=%q, want the missing anchor named", skuReason.Detail)
	}
```

## DEPOIS — HEAD, mesmo teste

```go
	// UNAVAILABLE, not INCOMPARABLE: the ERP product resolved by the EAN always
	// carries a CODPROD, so the seller_sku comparison has a value on BOTH sides
	// (they simply disagree) and the seeded reason falls to the excluded branch.
	// This is the same shape TestExactSKUWithUnmatchedListingEANKeepsSeededEANReason
	// already pins for `ean` — seller_sku was the outlier only because it read
	// `refforn` (B-01). Which direction both-present-and-DIFFERENT should carry is
	// the open AGAINST branch of A2-R1, the operator's decision, not this chip's.
	skuReason, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionUnavailable)
	if !ok {
		t.Fatalf("reasons=%#v, want seller_sku UNAVAILABLE (excluded branch)", candidate.Reasons)
	}
	if skuReason.Side != "" {
		t.Fatalf("seller_sku reason=%#v, want side empty on the excluded branch", skuReason)
	}
	if !strings.HasPrefix(skuReason.Detail, "sem CODPROD para corroborar o EAN") {
		t.Fatalf("seller_sku reason detail=%q, want the missing anchor named", skuReason.Detail)
	}
```

O `Detail` não mudou. O que mudou é a direção, mais uma asserção NOVA sobre `Side` que antes não
existia — a asserção depois é estritamente mais forte que a de antes, não mais fraca.

## Por que a asserção antiga só valia sobre o comportamento errado

O fixture do teste (`InternalProductID: canonicalIDPtr(300)`, sem `ReferenceCode`) só produzia
`INCOMPARABLE`/`side=erp` porque `identityAnchorValues` lia `product.ReferenceCode` — `refforn`,
que neste fixture é vazio. O ramo `product == nil || productValue == ""` disparava, e o motivo
emitido dizia, sobre um produto com CODPROD 300, que o produto ERP **não tem CODPROD**. Era essa
frase falsa que a asserção antiga fixava.

Com o CODPROD canônico do lado ERP, os dois lados têm valor (`"SKU-NO-MATCH"` contra `"300"`), o
ramo `default` dispara, e a direção é `UNAVAILABLE` com `Side` vazio.

## A simetria alegada, verificável no arquivo

`TestExactSKUWithUnmatchedListingEANKeepsSeededEANReason` (mesmo arquivo) já fixa exatamente essa
forma para `ean`: anúncio `EAN: "EAN-LISTING"`, produto `EAN: stringPtr("EAN-ERP")` — presentes dos
dois lados e diferentes — e assere `LinkCandidateReasonDirectionUnavailable` com `Side != ""`
proibido. `seller_sku` era a única âncora fora dessa simetria, e era por ler `refforn`.

## O resíduo, que este chip NÃO fecha (R-24 / HUB R3)

`UNAVAILABLE` sob D-B significa "o provider não fornece essa âncora". Aqui o provider FORNECE
(`"SKU-NO-MATCH"`). O balde está errado para o caso both-present-and-DIFFERENT nas DUAS âncoras —
`ean` inclusive, desde antes deste chip.

**O CORR-1 não cria esse defeito; torna-o frequente.** Antes, `seller_sku` só caía no `default`
quando o produto ERP tinha `refforn` preenchido. Agora cai sempre que o anúncio traz um
`seller_sku`, o que é o caminho comum de uma âncora primária. É essa mudança de peso — não uma
mudança de comportamento nova — que o hub está levando ao operador junto do ramo `AGAINST` de
A2-R1.
