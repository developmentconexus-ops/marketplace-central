import re

P1 = ("nunca|sempre|apenas|inalcanç|não pode|impossív|todo |toda |todos|todas|"
      "nenhum|qualquer|única|único|só o|só a|never|always|only|unreachable|cannot|every|no longer")
cs = re.compile("(" + P1 + ")")
ci = re.compile("(" + P1 + ")", re.I)

CASES = {
    "os quatro universais pré-existentes": "os quatro universais pré-existentes",
    "Those paths carry no product at all, and they are pinned by":
        "Those paths carry no product\nat all, and they are pinned by",
    "(control) The zero value itself is NOT unreachable in production":
        "The zero value itself is NOT unreachable in production",
}
for label, text in CASES.items():
    print(f"{label!r}\n  cs={cs.findall(text)}  ci={ci.findall(text)}")
