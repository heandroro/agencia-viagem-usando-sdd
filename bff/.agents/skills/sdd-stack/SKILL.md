---
name: sdd-stack
description: Use when the user mentions a specific technology, library, database, platform, or programming language in the context of a spec or design, or when .sdd/references/ files exist in the project. Provides cached reference documentation for the project stack so the agent does not need to fetch from the internet on every session.
---

# SDD Stack References

## Quando consultar este guia

Consulte quando o usuário estiver:
- Criando ou revisando `design.md` e mencionar uma tecnologia específica
- Perguntando sobre restrições, contratos ou padrões de uma lib/plataforma
- Iniciando uma nova sessão em um projeto com stack já definida

## Como funciona o cache

Cada tecnologia do projeto tem um arquivo de referência em `.sdd/references/`:

```
.sdd/
  references/
    postgres.md       ← tipos, limites, padrões de query
    kafka.md          ← conceitos, configs, anti-padrões
    stripe.md         ← contratos de API, webhooks, idempotência
    ...
```

Esses arquivos são gerados pelo comando `/cachear-referencia` e são versionados
no repositório — disponíveis offline, sem buscar na internet a cada sessão.

## Como usar em specs

Ao trabalhar em `design.md` que envolva uma tecnologia cacheada:

1. Verifique se `.sdd/references/<tech>.md` existe
2. Se existir, leia o arquivo antes de propor qualquer decisão de design
3. Se não existir e a tecnologia for relevante para o spec, sugira ao humano:
   > "Não encontrei referência cacheada para `<tech>`. Quer que eu busque agora
   > com `/cachear-referencia <url> <tech>`?"

## Atualizar referências desatualizadas

Se o arquivo cacheado tiver uma data antiga ou o humano mencionar que a
documentação mudou, use `/cachear-referencia` novamente — o prompt perguntará
se deve sobrescrever.

## Carregar referência relevante

Quando identificar que `.sdd/references/<tech>.md` é relevante para o contexto atual:

LOAD .sdd/references/<tech>.md
