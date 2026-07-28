const providerDisplayNames = { mercado_livre: "Mercado Livre" };
const INJECTIVE_PROVIDER_SLUG = /^[a-z0-9]+(_[a-z0-9]+)*$/;

function typesetSlug(providerCode) {
  return providerCode
    .split("_")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}

function providerDisplayName(providerCode) {
  const mapped = providerDisplayNames[providerCode];
  if (mapped) return mapped;
  if (!INJECTIVE_PROVIDER_SLUG.test(providerCode)) return providerCode;
  return typesetSlug(providerCode);
}

const inputs = [
  "mercado_livre",
  "amazon",
  "Amazon",
  "amazon-marketplace",
  "amazon_marketplace",
  "Amazon_Marketplace",
  "amazon__market",
  "_amazon",
  "amazon_",
  "shopee",
  "Shopee",
  "AMAZON",
  "amazon marketplace",
  "amazon2",
  "2amazon",
  "Amazon Marketplace",
];

for (const i of inputs) {
  console.log(JSON.stringify(i).padEnd(24), "->", JSON.stringify(providerDisplayName(i)));
}

console.log("--- COLISOES: mesma saida, entradas distintas ---");
const seen = new Map();
for (const i of inputs) {
  const o = providerDisplayName(i);
  if (!seen.has(o)) seen.set(o, []);
  seen.get(o).push(i);
}
let n = 0;
for (const [o, ins] of seen) {
  if (ins.length > 1) {
    n += 1;
    console.log(JSON.stringify(o), "<-", JSON.stringify(ins));
  }
}
console.log("total de classes em colisao:", n);
