# SDD Workflow — Referência Completa

## Memory Bank (leitura obrigatória no início de toda sessão)

Antes de qualquer ação, leia nesta ordem:
1. `.sdd/memory-bank/constitution.md` — princípios imutáveis
2. `.sdd/memory-bank/architecture.md` — decisões arquiteturais
3. `.sdd/memory-bank/product.md` — contexto de produto
4. `.sdd/memory-bank/apm-standards.md` — padrões de observabilidade e APM CLI

Se algum destes arquivos não existir, informe ao humano antes de prosseguir.

---

## Etapa 1 — Requirements (`requirements.md`)

### O que fazer
1. Leia o memory bank completo
2. Preencha os requisitos funcionais separados de técnicos
3. Atribua IDs `REQ-x.x` a cada requisito
4. Preencha as seções OBS-x de observabilidade (o **que** observar, não o como)
5. Execute o checklist ao final e reporte qualquer item não atendido

### Regras
- Nunca incluir detalhes técnicos — pertencem ao `design.md`
- Linguagem funcional: "o sistema deve..." não "usar Redis para..."
- Critérios de aceite mensuráveis para cada requisito
- Não avançar sem aprovação humana explícita

---

## Etapa 2 — Design (`design.md`)

### O que fazer
1. Confirme que `requirements.md` foi aprovado pelo humano
2. Verifique alinhamento com `architecture.md` antes de propor qualquer design
3. Preencha **obrigatoriamente** a seção APM / Observability Design
4. Para cada requisito OBS-x em `requirements.md`, crie APM-Mx ou APM-Ex correspondente
5. Execute o checklist ao final e reporte qualquer item não atendido

### Regras
- IDs de telemetria: `APM-Mx` para métricas, `APM-Ex` para eventos de negócio
- Contratos de interface documentados antes de qualquer geração de código
- Não alterar contratos estabelecidos sem aprovação humana
- Dependências externas apenas se listadas em `architecture.md`
- Não avançar sem aprovação humana explícita

---

## Etapa 3 — Tasks (`tasks.md`)

### O que fazer
1. Confirme que `design.md` foi aprovado pelo humano
2. Gere tasks de implementação com rastreabilidade `[REQ-x.x]`
3. **Sempre inclua T-APM-01 a T-APM-05** — nunca são opcionais
4. Tasks de APM devem referenciar IDs `APM-Mx`/`APM-Ex` do `design.md`
5. Execute o checklist ao final

### Regras
- Uma task = uma sessão de trabalho; subdividir se parecer grande demais
- Tasks de instrumentação APM não podem ser removidas ou marcadas como opcionais
- O checklist final deve incluir promoção de decisões para `architecture.md`

---

## Execução de Tasks

- Execute **UMA task por vez**
- Antes de iniciar cada task, confirme com o humano
- Após completar, marque `[x]` e aguarde aprovação antes de prosseguir
- Ao final, verifique se há decisões arquiteturais a promover para `architecture.md`

---

## Quebrando features grandes em sub-specs

Proponha divisão **antes de iniciar o design** quando:
- `requirements.md` cobrir mais de 2 domínios de negócio distintos
- A estimativa de tasks ultrapassar 10 itens
- Componentes diferentes puderem ser entregues e testados independentemente

```
.sdd/specs/
  checkout/                    ← spec de alto nível (orienta, não executa)
  checkout-pagamento/          ← sub-spec com ciclo SDD próprio
  checkout-carrinho/           ← sub-spec com ciclo SDD próprio
```

`tasks.md` dentro de cada sub-spec permanece **arquivo único** — a divisão é por escopo.

---

## Ciclo de vida da spec (spec-first)

1. Spec criada para guiar a implementação
2. Requisitos → Design → Tasks → Execução → Validação APM
3. Spec **descartada** após entrega (spec-first)
4. Conhecimento relevante **promovido** para `architecture.md` antes do descarte
5. Informações arquiteturais registradas como ADR em `.sdd/adr/` se necessário
