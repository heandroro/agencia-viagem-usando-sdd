---
description: "SDD design rules — enforces architecture.md alignment check, mandatory APM/Observability Design section, OBS-x to APM-x traceability, and interface contract stability on every design.md inside .sdd/specs/."
applyTo: ".sdd/specs/**/design.md"
---

# SDD — Regras para design.md

- Antes de propor qualquer decisão técnica, leia `.sdd/memory-bank/architecture.md` e verifique alinhamento
- Se houver conflito com `architecture.md` ou `constitution.md`, reporte ao humano antes de prosseguir
- A seção **APM / Observability Design** é obrigatória — nunca omita
- Para cada requisito OBS-x em `requirements.md`, deve existir um item APM-Mx ou APM-Ex correspondente no design
- IDs de telemetria seguem o padrão `APM-Mx` (métricas) e `APM-Ex` (eventos)
- Não altere contratos de interface (APIs, eventos, SDKs) já estabelecidos sem aprovação humana explícita documentada no spec
- Não avance para `tasks.md` sem aprovação humana explícita deste arquivo
- Ao finalizar, execute o checklist da seção de design e reporte itens não atendidos
