---
description: Extrai uma decisão arquitetural da sessão atual e a registra como ADR em .sdd/adr/, atualizando também architecture.md com a referência ao novo ADR.
argument-hint: "descrição da decisão em uma frase (ex: adotar circuit breaker para chamadas ao gateway de pagamento)"
input:
  - decisao: "Descrição da decisão arquitetural em uma frase"
allowed-tools: [Read, Write, Edit]
---

Você é um arquiteto SDD. Registre a decisão arquitetural informada como ADR.

**Decisão**: ${input:decisao}

## Passos

1. Leia `.sdd/memory-bank/architecture.md` para contexto e para identificar o
   próximo número de ADR disponível.

2. Leia o template em `.sdd/adr/_template.md` (se existir).

3. Crie `.sdd/adr/ADR-NN-titulo-kebab-case.md` com:
   - Título derivado da decisão informada
   - Status: `Accepted`
   - Contexto: por que esta decisão foi necessária (extraído da sessão atual)
   - Decisão: o que foi decidido
   - Consequências: impactos positivos e negativos
   - Data: hoje

4. Atualize `.sdd/memory-bank/architecture.md` adicionando referência ao novo ADR
   na seção de ADRs (ou criando a seção se não existir).

5. Confirme as alterações ao humano antes de salvar, se o impacto for significativo.
