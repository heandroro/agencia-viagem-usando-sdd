# Architecture — Decisões Arquiteturais

## Estilo Arquitetural

**Arquitetura de 3 Camadas com API Gateway/BFF**
- Frontend SPA → BFF → Backend Core → Database/Cache

## Linguagens e Plataformas

| Camada | Tecnologia | Justificativa |
|--------|------------|---------------|
| Frontend | React + TypeScript | Componentização, tipagem forte, ecossistema maduro |
| BFF | Python (FastAPI) | Produtividade, async nativo, excelente para orquestração |
| Backend Core | Go (Gin/Echo) | Performance, concorrência, confiabilidade em transações |
| Database | PostgreSQL | Transações ACID, consistência forte, JSONB para flexibilidade |
| Cache | Valkey | Compatibilidade Redis, performance, persistência opcional |
| Message Queue | Valkey Streams | Eventos assíncronos sem adicionar nova dependência |

## Decisões Arquiteturais (ADRs)

### ADR-001: Separação BFF + Backend Core
**Status**: Aceito
**Contexto**: Frontend React precisa de dados agregados de múltiplos domínios
**Decisão**: BFF (Python) orquestra chamadas ao Backend Core (Go)
**Consequências**: 
- Frontend desacoplado da complexidade interna
- BFF pode fazer caching e transformação de dados
- Latência extra de uma hop de rede

### ADR-002: Valkey para Cache e Eventos
**Status**: Aceito
**Contexto**: Necessidade de cache distribuído e eventos assíncronos
**Decisão**: Valkey para ambos (cache e streams de eventos)
**Consequências**:
- Stack simplificada (uma ferramenta para dois propósitos)
- Cache hit/miss para catálogo e buscas
- Eventos de confirmação de reserva via Streams

### ADR-003: PostgreSQL como Database Principal
**Status**: Aceito
**Contexto**: Necessidade de transações ACID fortes para reservas e consistência de dados
**Decisão**: PostgreSQL para catálogo, reservas, usuários e pagamentos
**Consequências**:
- Transações ACID robustas para reservas (evita overbooking)
- Consistência forte e integridade referencial
- JSONB disponível para campos flexíveis quando necessário
- Melhor suporte a locking e concorrência

## Diagrama de Componentes (C4 - Nível 2)

```
┌─────────────────────────────────────────────────────────────┐
│                         Browser                              │
│                      (React SPA)                             │
└──────────────────────────┬──────────────────────────────────┘
                           │ HTTPS
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                      BFF (Python/FastAPI)                    │
│  • Autenticação/Autorização                                  │
│  • Agregação de dados para frontend                          │
│  • Rate limiting                                             │
└──────────────────────────┬──────────────────────────────────┘
                           │ gRPC/HTTP
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                   Backend Core (Go)                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Catálogo   │  │  Reservas   │  │     Pagamentos      │  │
│  │  Service    │  │  Service    │  │      Service        │  │
│  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘  │
└─────────┼────────────────┼───────────────────┼─────────────┘
          │                │                   │
          ▼                ▼                   ▼
     ┌─────────┐      ┌─────────┐       ┌─────────────┐
     │ MongoDB │      │ MongoDB │       │   Stripe    │
     │(Pacotes)│      │(Reservas)│      │    API      │
     └─────────┘      └─────────┘       └─────────────┘
          │                │
          └────────────────┘
                    │
                    ▼
            ┌──────────────┐
            │   Valkey     │
            │ Cache/Events │
            └──────────────┘
```

## Dependências Externas

| Dependência | Uso | Integração |
|-------------|-----|------------|
| Stripe | Processamento de pagamentos | API REST |
| Valkey | Cache e message queue | Redis protocol |
| MongoDB Atlas (ou self-hosted) | Persistência | MongoDB driver |

## Padrões de Código

- **Go**: Repository pattern, handlers REST, services de domínio
- **Python**: Dependency injection, async/await, Pydantic para schemas
- **React**: Componentes funcionais, hooks customizados, context API para estado global

## Segurança

- JWT para autenticação stateless
- HTTPS obrigatório em produção
- Rate limiting no BFF
- Sanitização de inputs (zod/validator)
- Nunca logar PII ou dados de cartão
