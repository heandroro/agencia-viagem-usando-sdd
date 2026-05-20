# Reserva de Pacotes — Design

## Visão Geral

O fluxo de reserva é implementado em 3 camadas: Frontend React gerencia o wizard de reserva, BFF Python (FastAPI) orquestra chamadas e validações, Backend Go (Gin) executa regras de negócio e persistência no MongoDB. Reservas pendentes usam TTL de 30 minutos via worker assíncrono consumindo de Valkey Streams.

## Interface/API

### POST /api/v1/reservations
**Descrição**: Inicia uma nova reserva para um pacote.

**Request**:
```json
{
  "package_id": "pkg_abc123",
  "start_date": "2026-06-15",
  "end_date": "2026-06-20",
  "traveler_count": 2
}
```

**Response 201**:
```json
{
  "reservation_id": "res_xyz789",
  "status": "pending",
  "package": {
    "id": "pkg_abc123",
    "name": "Rio de Janeiro - 5 dias",
    "destination": "Rio de Janeiro"
  },
  "dates": {
    "start": "2026-06-15",
    "end": "2026-06-20",
    "nights": 5
  },
  "pricing": {
    "package_price": 500.00,
    "total": 5000.00,
    "currency": "BRL"
  },
  "expires_at": "2026-05-20T13:00:00Z"
}
```

**Response 409 (Unavailable)**:
```json
{
  "error": "package_unavailable",
  "message": "Pacote não disponível para as datas selecionadas"
}
```

### PUT /api/v1/reservations/{id}/travelers
**Descrição**: Atualiza dados dos viajantes da reserva.

**Request**:
```json
{
  "travelers": [
    {
      "type": "primary",
      "full_name": "João Silva",
      "document_type": "cpf",
      "document_number": "12345678901",
      "birth_date": "1985-03-15"
    },
    {
      "type": "companion",
      "full_name": "Maria Silva",
      "document_type": "cpf",
      "document_number": "98765432101",
      "birth_date": "1987-07-20"
    }
  ]
}
```

**Response 200**:
```json
{
  "reservation_id": "res_xyz789",
  "status": "pending",
  "travelers": [
    {
      "traveler_id": "trv_001",
      "type": "primary",
      "full_name": "João Silva"
    },
    {
      "traveler_id": "trv_002",
      "type": "companion",
      "full_name": "Maria Silva"
    }
  ]
}
```

**Response 400 (Validation)**:
```json
{
  "error": "validation_error",
  "message": "Dados dos viajantes inválidos",
  "details": [
    {
      "field": "travelers[0].document_number",
      "error": "invalid_cpf"
    }
  ]
}
```

### GET /api/v1/reservations/{id}/summary
**Descrição**: Retorna resumo completo da reserva para revisão.

**Response 200**:
```json
{
  "reservation_id": "res_xyz789",
  "status": "pending",
  "package": {
    "name": "Rio de Janeiro - 5 dias",
    "hotel": "Hotel Copacabana Palace",
    "flight": "LATAM - Ida e Volta",
    "destination": "Rio de Janeiro"
  },
  "dates": {
    "check_in": "2026-06-15",
    "check_out": "2026-06-20",
    "nights": 5
  },
  "travelers": [
    {
      "full_name": "João Silva",
      "document_masked": "***45678901",
      "type": "primary"
    },
    {
      "full_name": "Maria Silva",
      "document_masked": "***65432101",
      "type": "companion"
    }
  ],
  "pricing": {
    "package_price": 500.00,
    "subtotal": 5000.00,
    "taxes": 150.00,
    "total": 5150.00,
    "currency": "BRL"
  },
  "policies": {
    "cancellation": "Cancelamento gratuito até 48h antes do check-in. Após isso, taxa de 10%.",
    "modification": "Modificações permitidas até 72h antes, sujeito a disponibilidade"
  },
  "expires_at": "2026-05-20T13:00:00Z"
}
```

### Contratos de Evento

#### reservation.created
```json
{
  "event_type": "reservation.created",
  "payload": {
    "reservation_id": "res_xyz789",
    "user_id": "usr_sha256_abc",
    "package_id": "pkg_abc123",
    "amount": 5150.00,
    "currency": "BRL",
    "traveler_count": 2,
    "expires_at": "2026-05-20T13:00:00Z"
  },
  "timestamp": "2026-05-20T12:30:00Z",
  "correlation_id": "abc123-def456"
}
```

#### reservation.expired
```json
{
  "event_type": "reservation.expired",
  "payload": {
    "reservation_id": "res_xyz789",
    "user_id": "usr_sha256_abc",
    "elapsed_minutes": 30,
    "reason": "timeout"
  },
  "timestamp": "2026-05-20T13:00:00Z",
  "correlation_id": "abc123-def456"
}
```

#### reservation.travelers_updated
```json
{
  "event_type": "reservation.travelers_updated",
  "payload": {
    "reservation_id": "res_xyz789",
    "traveler_count": 2,
    "primary_traveler": "João Silva"
  },
  "timestamp": "2026-05-20T12:35:00Z",
  "correlation_id": "abc123-def456"
}
```

## Modelo de Dados

### Collection: reservations
```
- _id: ObjectId (reservation_id)
- user_id: String (hashed user id)
- package_id: String (referência ao pacote)
- status: Enum [pending, confirmed, cancelled, expired]
- created_at: DateTime
- updated_at: DateTime
- expires_at: DateTime (TTL de 30 min)
- dates: Object
  - start_date: Date
  - end_date: Date
  - nights: Number
- pricing: Object
  - package_price: Decimal128
  - subtotal: Decimal128
  - taxes: Decimal128
  - total: Decimal128
  - currency: String
- travelers: Array
  - traveler_id: String
  - type: Enum [primary, companion]
  - full_name: String
  - document_type: Enum [cpf, passport]
  - document_encrypted: String (AES-256)
  - birth_date: Date
- policies_accepted: Boolean
- audit: Object
  - ip_address: String
  - user_agent: String
- indexes:
  - { user_id: 1, status: 1 }
  - { expires_at: 1 }, expireAfterSeconds: 0
  - { package_id: 1, dates.start_date: 1, dates.end_date: 1 }
```

### Collection: availability (Materialized View)
```
- _id: ObjectId
- package_id: String
- date: Date
- available_slots: Number
- reserved_slots: Number
- version: Number (para optimistic locking)
- indexes:
  - { package_id: 1, date: 1 }
  - { available_slots: 1 }
```

## Diagrama de Sequência

### Fluxo: Criar Reserva (REQ-1.1)

```
Frontend    BFF(Python)    Backend(Go)    MongoDB    Valkey
   │              │              │            │          │
   │ POST /res    │              │            │          │
   │─────────────>│              │            │          │
   │              │ 1. Valida    │            │          │
   │              │    sessão    │            │          │
   │              │              │            │          │
   │              │ 2. Check     │            │          │
   │              │    Cache     │───────────>│          │
   │              │    Catálogo  │  GET pkg   │          │
   │              │              │<───────────│          │
   │              │              │            │          │
   │              │ 3. Cria      │            │          │
   │              │    Reserva   │────────────│          │
   │              │              │ 4. Check   │          │
   │              │              │    Disp.   │          │
   │              │              │───────────>│          │
   │              │              │ 5. Lock    │          │
   │              │              │    Atomic  │          │
   │              │              │<──────────│          │
   │              │              │ 6. Insert  │          │
   │              │              │    Reserva │─────────>│
   │              │              │<──────────│          │
   │              │              │ 7. Pub     │          │
   │              │              │    Event   │─────────>│ Stream
   │              │              │            │          │
   │              │<─────────────│            │          │
   │ 8. Resp 201  │              │            │          │
   │<─────────────│              │            │          │
```

### Fluxo: Expirar Reserva (REQ-1.4)

```
Worker(Go)    Valkey    MongoDB
   │            │          │
   │ 1. Consumo │          │
   │    Stream  │────────>│ XREAD
   │            │          │
   │ 2. Msg     │          │
   │    expire  │<────────│
   │            │          │
   │ 3. Update  │          │
   │    status  │────────>│ updateOne
   │            │          │ {status: expired}
   │            │          │
   │ 4. Release │          │
   │    slots   │          │
   │            │          │
   │ 5. Ack     │          │
   │            │────────>│ XACK
   │            │          │
```

## APM / Observability Design

### Métricas (APM-Mx)

| ID | Nome | Tipo | Labels | Descrição |
|----|------|------|--------|-----------|
| APM-M1 | business.reservation.created | Counter | package_id, destination | Reserva criada com sucesso |
| APM-M2 | business.reservation.failed | Counter | reason | Falha ao criar reserva |
| APM-M3 | business.reservation.travelers.captured | Counter | count | Viajantes registrados |
| APM-M4 | business.reservation.summary_viewed | Counter | - | Visualização de resumo |
| APM-M5 | business.reservation.expired | Counter | reason | Reservas expiradas |
| APM-M6 | reservation.api.duration | Histogram | endpoint, method | Latência da API |
| APM-M7 | reservation.db.duration | Histogram | operation | Latência MongoDB |
| APM-M8 | reservation.cache.hit | Counter | cache_name | Cache hit |
| APM-M9 | reservation.cache.miss | Counter | cache_name | Cache miss |

### Eventos (APM-Ex)

| ID | Nome | Atributos | Trigger |
|----|------|-----------|---------|
| APM-E1 | reservation.created | reservation_id, user_id, amount, package_id | Após INSERT bem-sucedido |
| APM-E2 | reservation.create_failed | reason, package_id, error | Após falha de disponibilidade |
| APM-E3 | reservation.travelers_updated | reservation_id, traveler_count | Após PUT /travelers |
| APM-E4 | reservation.validation_failed | field, error_code | Quando validação falha |
| APM-E5 | reservation.expired | reservation_id, elapsed_minutes | Worker de expiração |

### Traces

- **Span name**: `initiate_reservation`
  - **Attributes**: `package_id`, `destination`, `traveler_count`, `total_amount`
  
- **Span name**: `capture_traveler_data`
  - **Attributes**: `reservation_id`, `traveler_count`, `validation_time_ms`
  
- **Span name**: `generate_reservation_summary`
  - **Attributes**: `reservation_id`, `generation_time_ms`
  
- **Span name**: `expire_reservation`
  - **Attributes**: `reservation_id`, `elapsed_minutes`, `reason`

## Decisões de Design

### DD-001: TTL de Reserva via MongoDB vs Worker Manual
**Alternativas consideradas**:
- Opção A: MongoDB TTL Index (automático, simples, mas sem lógica custom)
- Opção B: Worker manual com Valkey Streams (mais controle, notificações)

**Decisão**: Usar ambos - TTL Index como fallback + Worker para eventos
**Justificativa**: TTL Index garante cleanup mesmo se worker falhar; Worker publica eventos para notificações e métricas de negócio

### DD-002: Criptografia de Documentos
**Alternativas consideradas**:
- Opção A: MongoDB Client-Side Field Level Encryption (CSFLE)
- Opção B: Criptografia manual AES-256 no Backend

**Decisão**: Criptografia AES-256 manual + hashing para logs
**Justificativa**: Maior controle sobre chaves; hashing permite correlacionar em logs sem expor dados

### DD-003: Controle de Disponibilidade
**Alternativas consideradas**:
- Opção A: Pessimistic Lock (bloqueia leitura, mais seguro)
- Opção B: Optimistic Locking (version field, melhor performance)

**Decisão**: Optimistic Locking na collection availability
**Justificativa**: Melhor performance para leituras concorrentes; retry em caso de conflito

### DD-004: Estrutura de Preços
**Alternativas consideradas**:
- Opção A: Calcular dinamicamente (baseado em preço atual do pacote)
- Opção B: Congelar preço no momento da reserva

**Decisão**: Congelar preço no momento da reserva (snapshot pricing)
**Justificativa**: Evita surpresas para o usuário; preço mostrado é o preço pago

## Alinhamento com Architecture

- [x] Consistente com ADR-001 (BFF orquestra, Backend Core executa)
- [x] Usa dependências listadas em architecture.md (MongoDB, Valkey)
- [x] Não introduz novas dependências externas
- [x] Segue padrões de código definidos (Repository pattern Go, FastAPI Python)

## Checklist de Validação

- [x] Cada REQ-x.x tem detalhamento técnico
- [x] OBS-x mapeado para APM-Mx ou APM-Ex
- [x] Contratos de API documentados
- [x] Modelo de dados definido
- [x] Diagrama de sequência para fluxo principal
- [x] Não há conflito com constitution.md ou architecture.md

---
**Status**: review
**Próximo passo**: Após aprovação, gerar tasks.md
