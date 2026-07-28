// Enumeração por RAMO, não por lista escolhida a mão.
const providerDisplayNames = { mercado_livre: "Mercado Livre" };
const INJECTIVE_PROVIDER_SLUG = /^[a-z0-9]+(_[a-z0-9]+)*$/;

function typesetSlug(code) {
  return code.split("_").map((w) => w.charAt(0).toUpperCase() + w.slice(1)).join(" ");
}

// devolve [saida, ramo]
function providerDisplayNameB(code) {
  const mapped = providerDisplayNames[code];
  if (mapped) return [mapped, "MAPPED"];
  if (!INJECTIVE_PROVIDER_SLUG.test(code)) return [code, "VERBATIM"];
  return [typesetSlug(code), "TYPESET"];
}

// População derivada do MECANISMO: para cada chave mapeada, o seu literal de saída
// também é um provider_code candidato (registry.go dedupa por igualdade exata de string).
const inputs = new Set();
for (const [k, v] of Object.entries(providerDisplayNames)) {
  inputs.add(k);
  inputs.add(v);
}
for (const c of ["amazon", "Amazon", "amazon_marketplace", "Amazon Marketplace", "shopee", "Shopee", "amazon-marketplace", "amazon__market", "mercado livre", "MERCADO_LIVRE"]) {
  inputs.add(c);
}

console.log("--- SAIDA POR RAMO ---");
const byOut = new Map();
for (const i of [...inputs].sort()) {
  const [out, branch] = providerDisplayNameB(i);
  console.log(JSON.stringify(i).padEnd(24), "->", JSON.stringify(out).padEnd(24), branch);
  if (!byOut.has(out)) byOut.set(out, []);
  byOut.get(out).push(`${i} [${branch}]`);
}

console.log("--- COLISOES, com o RAMO de cada origem ---");
let n = 0;
for (const [out, ins] of byOut) {
  if (ins.length > 1) {
    n += 1;
    console.log("COLLIDE", JSON.stringify(out), "<-", ins.join(" | "));
  }
}
console.log("classes em colisao:", n);

// --- A opcao (a) do docstring: "render every unmapped code verbatim (injective)" ---
console.log();
console.log('--- OPCAO (a) do docstring: "render every unmapped code verbatim (injective, uglier)" ---');
function optionA(code) {
  const mapped = providerDisplayNames[code];
  if (mapped) return [mapped, "MAPPED"];
  return [code, "VERBATIM"];
}
const byOutA = new Map();
for (const i of [...inputs].sort()) {
  const [out, branch] = optionA(i);
  if (!byOutA.has(out)) byOutA.set(out, []);
  byOutA.get(out).push(`${i} [${branch}]`);
}
let nA = 0;
for (const [out, ins] of byOutA) {
  if (ins.length > 1) {
    nA += 1;
    console.log("COLLIDE", JSON.stringify(out), "<-", ins.join(" | "));
  }
}
console.log("classes em colisao sob a OPCAO (a):", nA);
console.log(nA === 0 ? "opcao (a) e injetiva nesta populacao" : "opcao (a) NAO e injetiva — o docstring chama de injetiva");
