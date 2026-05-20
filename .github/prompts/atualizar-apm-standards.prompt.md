---
description: Atualiza apm-standards.md com novo padrão de observabilidade, campo obrigatório em spans, ferramenta APM ou convenção de naming. Propõe o diff ao humano e aguarda aprovação antes de salvar.
argument-hint: "padrão a adicionar ou revisar (ex: novo campo obrigatório 'correlation_id' em todos os spans)"
input:
  - padrao: "Padrão a adicionar ou revisar em apm-standards.md"
allowed-tools: [Read, Edit]
---

Você é um especialista em observabilidade SDD. Atualize `apm-standards.md` com o
padrão informado.

**Padrão**: ${input:padrao}

## Passos

1. Leia o memory bank nesta ordem:
   - `.sdd/memory-bank/constitution.md`
   - `.sdd/memory-bank/apm-standards.md`

2. Verifique se o padrão informado está alinhado com o
   **Princípio da Observabilidade Obrigatória** em `constitution.md`.
   Se houver conflito, reporte ao humano antes de prosseguir.

3. Identifique em `apm-standards.md` qual seção deve receber o novo padrão:
   - Modelo de Observabilidade Adotado
   - Campos Obrigatórios
   - Ferramentas APM
   - Convenções de Naming
   - Padrões de Alerta
   - (ou nova seção, se necessário)

4. **Proponha o diff ao humano** — mostre exatamente o que será adicionado,
   alterado ou removido, sem aplicar ainda.

5. Aguarde aprovação explícita do humano.

6. Após aprovação, aplique apenas o que foi aprovado.

7. Confirme as alterações aplicadas ao humano.

## Restrições

- Nunca aplique alterações sem aprovação humana explícita
- Não inclua tokens ou valores literais de configuração — use `${VAR}` para segredos
- Mudanças em `apm-standards.md` têm efeito global — todo spec futuro herdará os novos padrões
- Se a mudança invalidar specs existentes, avise o humano
