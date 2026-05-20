---
name: sdd-memory-bank
description: Use when the user is reading or updating files in .sdd/memory-bank/, asking how to evolve the project memory bank, or wondering which file to update after a decision or product change. Provides reading order, editability rules per file, and promotion discipline.
---

# SDD Memory Bank

## Quando consultar este guia

Consulte quando o usuário estiver:
- Abrindo ou editando qualquer arquivo em `.sdd/memory-bank/`
- Perguntando como registrar uma nova decisão arquitetural, mudança de produto ou padrão de observabilidade
- Iniciando uma nova sessão de trabalho e precisando carregar o contexto do projeto
- Verificando se um spec ou task está em conformidade com os princípios do projeto

## Ordem de leitura obrigatória

No início de **toda sessão**, leia os arquivos na seguinte ordem:

```
1. constitution.md      ← princípios imutáveis
2. architecture.md      ← decisões arquiteturais vigentes
3. product.md           ← contexto do produto/negócio
4. apm-standards.md     ← padrões de observabilidade
```

Nunca inverta a ordem — `constitution.md` define o que é válido nas etapas seguintes.

## Matriz de editabilidade

| Arquivo | IA pode editar? | Como atualizar | Frequência |
|---|---|---|---|
| `constitution.md` | ❌ Nunca | Somente humano, explicitamente | Raramente |
| `architecture.md` | ⚠️ Só via `/promover-adr` | Após conclusão de task com decisão arquitetural | Por feature |
| `product.md` | ⚠️ Só via `/atualizar-produto` | Quando mudar escopo, segmento ou visão do produto | Por release/trimestre |
| `apm-standards.md` | ⚠️ Só via `/atualizar-apm-standards` | Quando adotar novo padrão de observabilidade | Por mudança de padrão |

## Regras de promoção

- **architecture.md** só recebe atualizações ao final de uma task, via checklist de promoção em `tasks.md` — nunca no meio da execução
- **product.md** e **apm-standards.md** são atualizados por decisão humana, não como efeito colateral de uma task de feature
- Se um dos quatro arquivos não existir, **reporte ao humano antes de prosseguir** com qualquer ação

## Comandos disponíveis

| Comando | Arquivo alvo | Quando usar |
|---|---|---|
| `/promover-adr` | `architecture.md` | Decisão técnica surgida durante execução de feature |
| `/atualizar-produto` | `product.md` | Nova visão, segmento, restrição de negócio |
| `/atualizar-apm-standards` | `apm-standards.md` | Novo pilar, campo obrigatório, mudança de ferramenta APM |

Para detalhes e regras de atualização por arquivo, LOAD references/update-guide.md.
