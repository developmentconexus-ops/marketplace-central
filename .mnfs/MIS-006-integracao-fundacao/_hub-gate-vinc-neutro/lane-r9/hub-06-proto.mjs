// Verifica por EXECUCAO os dois achados que os assentos de VINC r9 devolveram.
const providerDisplayNames = { mercado_livre: "Mercado Livre" };
const INJECTIVE_PROVIDER_SLUG = /^[a-z0-9]+(_[a-z0-9]+)*$/;
function typesetSlug(code) {
  return code.split("_").map((w) => w.charAt(0).toUpperCase() + w.slice(1)).join(" ");
}
function providerDisplayName(providerCode) {
  const mapped = providerDisplayNames[providerCode];
  if (mapped) return mapped;
  if (!INJECTIVE_PROVIDER_SLUG.test(providerCode)) return providerCode;
  return typesetSlug(providerCode);
}

console.log("=== ACHADO 1: cadeia de prototipo (os DOIS assentos) ===");
for (const k of ["toString", "constructor", "valueOf", "hasOwnProperty", "__proto__", "isPrototypeOf", "toLocaleString", "propertyIsEnumerable"]) {
  const out = providerDisplayName(k);
  console.log(
    String(k).padEnd(22),
    "typeof:", String(typeof out).padEnd(10),
    "verbatim?", out === k,
    "|", typeof out === "string" ? JSON.stringify(out) : String(out).slice(0, 40).replace(/\n/g, " "),
  );
}

console.log();
console.log("=== ACHADO 2 (Sol): saida da TRANSFORMACAO que o VERBATIM nao produz ===");
// A frase: "the fallback is identity, so every string a transform can produce is
// also a string the fallback can produce." Contra-exemplo precisa de saida do
// TYPESET que PASSE no regex (pois verbatim so emite quem FALHA no regex).
const probes = ["2amazon", "0", "9x", "amazon", "amazon_marketplace", "123_456"];
for (const p of probes) {
  if (!INJECTIVE_PROVIDER_SLUG.test(p)) { console.log(String(p).padEnd(20), "-> nao entra no TYPESET"); continue; }
  const out = typesetSlug(p);
  const reachableByVerbatim = !INJECTIVE_PROVIDER_SLUG.test(out); // verbatim so emite quem falha
  console.log(
    String(p).padEnd(20), "TYPESET ->", JSON.stringify(out).padEnd(24),
    "verbatim consegue emitir essa saida?", reachableByVerbatim,
  );
}
