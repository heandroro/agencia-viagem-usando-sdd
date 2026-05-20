---
description: Gera o tasks.md a partir de um design.md aprovado, com rastreabilidade REQ-x.x em cada task e incluindo obrigatoriamente T-APM-01 a T-APM-05.
argument-hint: "nome da pasta da feature em .sdd/specs/ (ex: order-checkout)"
input:
  - feature: "Nome da pasta da feature em .sdd/specs/ (ex: order-checkout)"
allowed-tools: [Read, Write]
---

Você é um engenheiro SDD. Gere o `tasks.md` para a feature informada.

**Feature**: ${input:feature}

## Passos

1. Confirme que `design.md` foi aprovado pelo humano. Se não houver confirmação
   explícita, pergunte antes de prosseguir.

2. Leia:
   - `.sdd/specs/${input:feature}/requirements.md`
   - `.sdd/specs/${input:feature}/design.md`
   - `.sdd/specs/_template/tasks.md`
   - `.sdd/memory-bank/apm-standards.md`

3. Gere `.sdd/specs/${input:feature}/tasks.md` com:
   - Tasks de implementação, cada uma referenciando `[REQ-x.x]`
   - **Obrigatório**: tasks T-APM-01 a T-APM-05 referenciando IDs `APM-Mx`/`APM-Ex` do `design.md`
   - Cada task deve ser implementável em uma única sessão de trabalho
   - Se uma task parecer grande, subdivida antes de incluir

4. Ao finalizar, execute o checklist do template e reporte itens não atendidos.

5. Aguarde aprovação humana antes de iniciar qualquer task.

> **Lembre**: T-APM-xx nunca são opcionais — mesmo que o humano não as mencione.
