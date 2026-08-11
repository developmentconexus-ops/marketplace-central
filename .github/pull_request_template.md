Closes #

## O que aterra

| Camada | Onde |
|---|---|
|  |  |

## Aceite do issue, respondido

Uma linha por aceite declarado no issue. A evidência é a saída do comando, não
a afirmação de que ele passou.

| Aceite declarado | Comando | Saída |
|---|---|---|
|  |  |  |

**Controlo positivo:** que sabotagem foi feita, o que ela imprimiu, e que o
aceite voltou ao verde depois de restaurada. Sem isto, o aceite não se sabe
falhar.

## Fora do escopo declarado

Cada ficheiro tocado que o issue não previa, com a razão. Vazio só se for
mesmo vazio — excesso não declarado é a metade do desvio de escopo que ninguém
vê, porque nada fica vermelho por causa dele.

| Ficheiro | Porquê entrou |
|---|---|
|  |  |

## Gate

```
gate: PASS
lanes_selected= lanes_ran= lanes_failed=
```

Baselines em `contracts/gate/baselines.json`: inalteradas / encolheram (qual e
porquê). Levantar baseline não é opção.

## Revisão fria de escopo

Veredito diff × issue, por quem não implementou: `COBRE` / `FALTA` / `EXCEDE`,
com `file:line`.
