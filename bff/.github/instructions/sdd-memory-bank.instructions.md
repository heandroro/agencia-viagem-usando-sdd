---
description: "SDD memory bank rules — enforces ordered read-on-session-start, constitution immutability for AI agents, and promote-to-architecture discipline on every file inside .sdd/memory-bank/."
applyTo: ".sdd/memory-bank/**"
---

# SDD — Regras para o Memory Bank

- No início de toda sessão, leia os arquivos na seguinte ordem obrigatória:
  1. `constitution.md` — princípios imutáveis
  2. `architecture.md` — decisões arquiteturais
  3. `product.md` — contexto de produto
  4. `apm-standards.md` — padrões de observabilidade e APM CLI
- `constitution.md` é **imutável para agentes de IA** — nunca proponha alterações sem solicitação humana explícita
- Alterações em `architecture.md` só devem ocorrer para promover decisões arquiteturais relevantes surgidas durante execução de features
- A promoção para `architecture.md` segue o checklist final de `tasks.md` — não ocorra antes da conclusão da task
- Se algum dos quatro arquivos do memory bank não existir, informe ao humano antes de prosseguir com qualquer ação
