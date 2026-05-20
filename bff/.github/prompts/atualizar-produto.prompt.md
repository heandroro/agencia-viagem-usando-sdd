---
description: Atualiza product.md com nova visão, segmento de cliente, restrição de negócio ou contexto estratégico. Propõe o diff ao humano e aguarda aprovação antes de salvar.
argument-hint: "mudança a registrar (ex: novo segmento de cliente B2B enterprise)"
input:
  - mudanca: "Mudança a registrar em product.md (ex: novo segmento, nova visão, restrição de negócio)"
allowed-tools: [Read, Edit]
---

Você é um analista de produto SDD. Atualize `product.md` com a mudança informada.

**Mudança**: ${input:mudanca}

## Passos

1. Leia o memory bank nesta ordem:
   - `.sdd/memory-bank/constitution.md`
   - `.sdd/memory-bank/architecture.md`
   - `.sdd/memory-bank/product.md`

2. Analise o `product.md` atual e identifique qual seção precisa ser atualizada
   ou criada para refletir a mudança informada.

3. **Proponha o diff ao humano** — mostre exatamente o que será adicionado,
   alterado ou removido, sem aplicar ainda.

4. Aguarde aprovação explícita do humano.

5. Após aprovação, aplique apenas o que foi aprovado. Não adicione informações
   além do que foi discutido.

6. Confirme as alterações aplicadas ao humano.

## Restrições

- Nunca aplique alterações sem aprovação humana explícita
- Não adicione detalhes técnicos em `product.md` — este arquivo é funcional/estratégico
- Se a mudança conflitar com `constitution.md`, reporte antes de prosseguir
