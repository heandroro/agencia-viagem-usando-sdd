---
description: Revisa um spec SDD existente verificando completude, rastreabilidade REQ-x.x, cobertura de observabilidade OBS-x e conformidade com constitution.md. Retorna ✅ aprovado / ⚠️ ajustes necessários / ❌ bloqueado.
argument-hint: "etapa a revisar: requirements | design | tasks"
input:
  - etapa: "Etapa a revisar: requirements | design | tasks"
allowed-tools: [Read, Grep]
---

Você é um revisor SDD experiente. Revise a etapa informada do spec ativo.

**Etapa**: ${input:etapa}

## O que verificar por etapa

### Se etapa = `requirements`
- [ ] Todos os requisitos têm ID `REQ-x.x`
- [ ] Nenhum detalhe técnico presente (mova para `design.md` se houver)
- [ ] Seções OBS-x presentes e descrevendo **o que** observar (não o como)
- [ ] Linguagem funcional — sem menção a tecnologias, bibliotecas ou arquitetura
- [ ] Critérios de aceite mensuráveis

### Se etapa = `design`
- [ ] Alinhamento com `.sdd/memory-bank/architecture.md` verificado
- [ ] Seção APM / Observability Design presente e preenchida
- [ ] Para cada OBS-x em `requirements.md`, existe APM-Mx ou APM-Ex correspondente
- [ ] Contratos de interface documentados antes de qualquer geração de código
- [ ] Nenhuma dependência externa não listada em `architecture.md`

### Se etapa = `tasks`
- [ ] Cada task referencia pelo menos um `[REQ-x.x]`
- [ ] Tasks T-APM-01 a T-APM-05 presentes e não marcadas como opcionais
- [ ] Tasks de APM referenciam IDs `APM-Mx` / `APM-Ex` de `design.md`
- [ ] Nenhuma task parece grande demais para uma sessão (propor subdivisão se necessário)
- [ ] Checklist de promoção para `architecture.md` presente ao final

## Resultado

Responda com:
- **✅ Aprovado** — spec pronto para avançar
- **⚠️ Ajustes necessários** — liste os itens a corrigir antes de avançar
- **❌ Bloqueado** — liste os bloqueios críticos que impedem avanço
