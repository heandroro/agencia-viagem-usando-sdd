---
description: "SDD requirements rules — enforces functional-only content, REQ-x.x IDs, OBS-x observability sections, and no-advance-without-approval policy on every requirements.md inside .sdd/specs/."
applyTo: ".sdd/specs/**/requirements.md"
---

# SDD — Regras para requirements.md

- Este arquivo contém **requisitos funcionais** (o QUÊ e o PORQUÊ) — nunca decisões técnicas (o COMO)
- Detalhes de implementação, tecnologias, bibliotecas e arquitetura pertencem ao `design.md`
- Se detectar detalhe técnico no arquivo, mova para `design.md` e reporte ao humano
- Separe claramente requisitos funcionais de requisitos de observabilidade (seções OBS-x)
- Cada requisito recebe ID único no formato `REQ-x.x` (ex: `REQ-1.1`, `REQ-2.3`)
- Cada seção OBS-x deve descrever **o que** observar, não **como** instrumentar
- Não avance para `design.md` sem aprovação humana explícita deste arquivo
- Ao finalizar, execute o checklist da seção de requisitos e reporte itens não atendidos
