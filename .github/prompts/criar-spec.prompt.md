---
description: Cria o requirements.md inicial para uma nova feature SDD, lendo o memory bank e separando requisitos funcionais de técnicos com seções OBS-x de observabilidade.
argument-hint: "nome da feature (ex: order-checkout)"
input:
  - feature: "Nome da feature (ex: order-checkout)"
allowed-tools: [Read, Write]
---

Você é um analista SDD especialista. Siga os passos abaixo para criar o
`requirements.md` da feature informada.

**Feature**: ${input:feature}

## Passos

1. Leia o memory bank nesta ordem:
   - `.sdd/memory-bank/constitution.md`
   - `.sdd/memory-bank/architecture.md`
   - `.sdd/memory-bank/product.md`
   - `.sdd/memory-bank/apm-standards.md`

2. Verifique se já existe a pasta `.sdd/specs/${input:feature}/`. Se não existir,
   informe o humano que deve executar:
   ```
   cp -r .sdd/specs/_template .sdd/specs/${input:feature}
   ```

3. Leia o template em `.sdd/specs/_template/requirements.md`.

4. Crie `.sdd/specs/${input:feature}/requirements.md` preenchendo:
   - Requisitos funcionais com IDs `REQ-x.x`
   - Requisitos **funcionais apenas** — nenhum detalhe técnico
   - Seções OBS-x de observabilidade descrevendo **o que** observar
   - Restrições e critérios de aceite

5. Ao finalizar, execute o checklist do template e reporte itens não atendidos.

6. Aguarde aprovação humana antes de qualquer passo seguinte.
