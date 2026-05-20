---
name: sdd-workflow
description: Use when the user is working on files inside .sdd/ or mentions spec, requirements, design, tasks, memory bank, or APM in a spec-driven development context. Provides the complete SDD workflow, checklists, and behavioral rules for each stage of the cycle.
---

# SDD Workflow

## Quando consultar este guia

Consulte quando o usuário estiver:
- Criando ou editando arquivos em `.sdd/specs/` ou `.sdd/memory-bank/`
- Perguntando sobre o ciclo SDD (requirements, design, tasks)
- Executando tasks de um `tasks.md`
- Promovendo decisões para `architecture.md` ou `adr/`
- Usando os comandos `/criar-spec`, `/revisar-spec`, `/gerar-tasks`, `/promover-adr`

## Visão geral do ciclo

```
Memory Bank (lido no início de toda sessão)
     ↓
requirements.md → [aprovação humana] → design.md → [aprovação humana] → tasks.md → execução task-a-task
```

Para detalhes completos de cada etapa, LOAD references/workflow.md.
Para os checklists de cada etapa, LOAD references/checklists.md.

## Regras fundamentais (constitution.md)

1. **Observabilidade obrigatória** — todo código em produção deve ser instrumentado; T-APM-xx nunca são opcionais
2. **Separação funcional × técnica** — `requirements.md` = o QUÊ; `design.md` = o COMO
3. **Contrato de interface** — interfaces documentadas em `design.md` antes de qualquer código
4. **Rastreabilidade** — cada task referencia `[REQ-x.x]`; cada task APM referencia `[APM-Mx]`/`[APM-Ex]`
5. **Tamanho de task** — uma task = uma sessão de trabalho; subdividir se parecer grande
6. **Revisão humana** — aprovação explícita entre cada etapa; nunca avance sozinho
7. **Consistência com memory bank** — `design.md` sempre alinhado com `architecture.md`
8. **Commit semântico** — todos os commits seguem [Conventional Commits](https://www.conventionalcommits.org/pt-br/)

## Comportamentos proibidos

- **NÃO** incluir dados PII/sensíveis em exemplos de telemetria
- **NÃO** avançar para a próxima etapa sem aprovação humana explícita
- **NÃO** introduzir dependências externas não listadas em `architecture.md`
- **NÃO** alterar contratos de interface estabelecidos sem aprovação
- **NÃO** pular as tasks T-APM-xx mesmo que o humano não as mencione
- **NÃO** incluir tokens/secrets literais — sempre usar `${VAR}`
- **NÃO** misturar detalhes técnicos em `requirements.md`

## Quando reportar ao humano antes de prosseguir

- Conflito entre spec e `constitution.md` ou `architecture.md`
- Task que precisa alterar contrato de interface existente
- Dependência não listada em `architecture.md`
- Task ambígua com múltiplas interpretações válidas
- Requisito de observabilidade incompleto ou inconsistente

## Tamanho ideal de problema

| Tamanho | Recomendação |
|---|---|
| Bug pequeno | Não use SDD — vá direto ao código |
| Feature média (3–8 pontos) | **Use o fluxo completo** |
| Feature grande | Quebre em sub-specs antes de iniciar |
| Produto novo | Comece pelo memory bank; uma spec por feature |
