# V-round9 — B-1 discharged by execution, and R-1 measured NOT closed with it

Tree: `bfc1d9bb` + round-9 working changes. Measured before commit.

## What the order said

> **B-1.** Duas saídas, escolha é sua porque o arquivo é seu: restringir a aplicação ao
> domínio onde a transformação É injetiva (`^[a-z0-9]+(_[a-z0-9]+)*$`) — caixa e hífen caem
> em verbatim e a frase vira verdadeira **sem reescrita**; ou DELETAR a frase (R-25) e
> declarar o escopo real. **A primeira fecha B-1 e R-1 juntos.**

Option 1 taken. The last sentence is the one this file answers.

## The new function, executed

`providerDisplayName` after the domain restriction, run in node against the same inputs
the hub used to prove B-1:

```
"amazon-marketplace"     -> "amazon-marketplace"
"amazon"                 -> "Amazon"
"Amazon"                 -> "Amazon"
"amazon__market"         -> "amazon__market"
"_amazon"                -> "_amazon"
"amazon_marketplace"     -> "Amazon Marketplace"
"mercado_livre"          -> "Mercado Livre"
```

## B-1: CLOSED

The hub's counterexample was `"amazon-marketplace" -> "Amazon-marketplace"` — the docstring
said hyphenated codes render verbatim and the code capitalized them. It now renders verbatim.
The total sentence is true over its stated domain, and the domain is CHECKED rather than
inferred from a round-trip.

## R-1: NOT CLOSED — the ruling is false, and here is why

`amazon` is INSIDE the domain and typesets to `Amazon`. `Amazon` is OUTSIDE the domain and
therefore renders VERBATIM as `Amazon`. Two distinct, simultaneously-registrable wire codes
(`registry.go:100-114` dedupes by exact string equality) paint the same name.

The domain restriction makes the TRANSFORM injective. It does not make the FUNCTION
injective, because the function is two branches sharing one codomain, and the second branch
is the IDENTITY. Every string the transform can produce is also a string the fallback can
produce, so no narrowing of the domain can separate them — narrowing only moves codes from
the transform branch to the fallback branch, which is where the collision partner lives.

The collision did not disappear. It moved from inside the transform to the seam between the
transform and its own escape hatch.

Two shapes actually close it, both outside this chip's write-set:

- render every unmapped code verbatim (injective; uglier, and drops the feature the
  typesetting exists for);
- require a display name in the registry beside the code (injective; a contract change).

R-1 therefore keeps its trigger: **the second registered adapter.** Today `mercado_livre` is
mapped literally, it is the only declaration that exists, and no live row reaches either
branch.

## What was written down as a result

- `QueueRow.tsx`: the docstring's stale `ROUND-TRIPS` premise corrected at source, and a
  closing paragraph that states the case collision is NOT closed and why narrowing cannot
  close it. The docstring no longer claims what the code does not do (R-24).
- `QueueTab.test.tsx`: **no** `Amazon` fixture and **no** case assertion. Asserting today's
  output would write the surviving defect down as the requirement — which is exactly how the
  hyphen case survived: the V10 test asserted `getByText("Amazon-marketplace")`, the old
  buggy output, and stayed green for it. That assertion was corrected to
  `getByText("amazon-marketplace")` plus a negative on `Amazon-marketplace`.
