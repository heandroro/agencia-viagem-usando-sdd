# [NOME DA FEATURE] — Design

## Visão Geral

[PREENCHER: Descrição técnica de alto nível da solução]

## Interface/API

### Endpoints

```
[METHOD] /api/v1/[endpoint]
```

**Request**:
```json
{
  "campo1": "tipo",
  "campo2": "tipo"
}
```

**Response 200**:
```json
{
  "id": "string",
  "status": "string"
}
```

**Response 4xx/5xx**:
```json
{
  "error": "codigo_erro",
  "message": "descrição"
}
```

### Contratos de Evento (se aplicável)

```json
{
  "event_type": "nome.evento",
  "payload": {
    "campo": "valor"
  },
  "timestamp": "ISO8601",
  "correlation_id": "uuid"
}
```

## Modelo de Dados

```
Collection/Table: [nome]
- campo1: tipo (descrição)
- campo2: tipo (descrição)
- indexes: [campos indexados]
```

## Diagrama de Sequência

```
Frontend → BFF → Backend → Database
   │         │        │         │
   │ ───────>│        │         │ (1. Request)
   │         │───────>│         │ (2. Process)
   │         │        │───────>│ (3. Query)
   │         │        │<───────│ (4. Result)
   │         │<───────│         │ (5. Response)
   │<───────│        │         │ (6. Render)
```

## APM / Observability Design

### Métricas (APM-Mx)

| ID | Nome | Tipo | Labels | Descrição |
|----|------|------|--------|-----------|
| APM-M1 | [nome] | counter/histogram/gauge | [labels] | [descrição] |

### Eventos (APM-Ex)

| ID | Nome | Atributos | Trigger |
|----|------|-----------|---------|
| APM-E1 | [nome] | [attrs] | [quando logar] |

### Traces

- **Span name**: `[nome]`
- **Attributes**: `[lista de atributos]`

## Decisões de Design

### DD-001: [Título da decisão]
**Alternativas consideradas**:
- Opção A: [descrição]
- Opção B: [descrição]

**Decisão**: [Qual foi escolhida]
**Justificativa**: [Por que]

## Alinhamento com Architecture

- [ ] Consistente com ADR-001 (se aplicável)
- [ ] Usa dependências listadas em architecture.md
- [ ] Não introduz novas dependências externas
- [ ] Segue padrões de código definidos

## Checklist de Validação

- [ ] Cada REQ-x.x tem detalhamento técnico
- [ ] OBS-x mapeado para APM-Mx ou APM-Ex
- [ ] Contratos de API documentados
- [ ] Modelo de dados definido
- [ ] Diagrama de sequência para fluxo principal
- [ ] Não há conflito com constitution.md ou architecture.md

---
**Status**: [draft | review | approved]
**Próximo passo**: Após aprovação, gerar tasks.md
