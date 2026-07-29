import io
import re
import sys

path = sys.argv[1]
t = io.open(path, encoding="utf-8").read()

P1 = ("nunca|sempre|apenas|inalcanç|não pode|impossív|todo |toda |todos|todas|"
      "nenhum|qualquer|única|único|só o|só a|never|always|only|unreachable|cannot|every|no longer")
# P3 = the class P1 is structurally blind to: closed enumerations and coverage claims.
# Registered here with its exact regex, because a count without its pattern is not evidence.
P3 = (r"os (dois|tr[eê]s|quatro|cinco|seis)|as (duas|tr[eê]s|quatro|cinco)|"
      r"pinned by|cobre|coberto|cobertura|garante|prova que|exclusiv|ningu[eé]m|"
      r"\bzero\b|nothing|none\b|each |all ")

ci = re.compile("(" + P1 + ")", re.I)
c3 = re.compile("(" + P3 + ")", re.I)

L = t.split("\n")
h1 = [(i + 1, l) for i, l in enumerate(L) if ci.search(l)]
h3 = [(i + 1, l) for i, l in enumerate(L) if c3.search(l)]
print("linhas:", len(L))
print("P1 ci hits:", len(h1))
print("P3 hits:", len(h3))
print("--- P3: enumeracao fechada / cobertura (a classe cega de P1) ---")
for i, l in h3:
    print(i, "|", l[:160])
