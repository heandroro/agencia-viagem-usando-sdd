# APM Standards — Padrões de Observabilidade

## Plataforma APM

**Stack**: OpenTelemetry + Grafana (Tempo + Loki + Prometheus)

Alternativas aceitas: Datadog, New Relic, Dynatrace (manter estrutura de atributos)

## Campos Obrigatórios (Todos os Spans/Logs)

### Contexto de Negócio
| Atributo | Descrição | Exemplo |
|----------|-----------|---------|
| `business.domain` | Domínio de negócio | `catalog`, `reservation`, `payment` |
| `business.operation` | Operação de negócio | `search_packages`, `create_reservation`, `process_payment` |
| `business.outcome` | Resultado da operação | `success`, `failure`, `timeout` |

### Contexto Técnico
| Atributo | Descrição | Exemplo |
|----------|-----------|---------|
| `service.name` | Nome do serviço | `bff-api`, `catalog-service` |
| `service.version` | Versão do serviço | `1.2.3` |
| `deployment.environment` | Ambiente | `prod`, `staging`, `dev` |
| `trace.trace_id` | ID de correlação | `abc123...` |

### Contexto de Requisição
| Atributo | Descrição | Exemplo |
|----------|-----------|---------|
| `http.method` | Método HTTP | `GET`, `POST` |
| `http.route` | Rota | `/api/v1/packages` |
| `http.status_code` | Código de resposta | `200`, `500` |
| `user.id` | ID do usuário (hash) | `usr_sha256_abc...` |
| `session.id` | ID da sessão | `sess_xyz789...` |

## Métricas de Negócio (Business Metrics)

### Catalog
```
business.catalog.search.performed  (counter + destination, dates)
business.catalog.search.results    (histogram + count)
business.catalog.cache.hit         (counter)
business.catalog.cache.miss        (counter)
```

### Reservation
```
business.reservation.created       (counter)
business.reservation.confirmed     (counter)
business.reservation.cancelled     (counter + reason)
business.reservation.checkout.duration  (histogram)
business.reservation.cart.abandoned (counter + step)
```

### Payment
```
business.payment.initiated         (counter)
business.payment.approved          (counter)
business.payment.declined          (counter + reason)
business.payment.refunded          (counter)
```

## Logs Estruturados

```json
{
  "timestamp": "2026-05-20T10:30:00Z",
  "level": "INFO",
  "message": "Reserva criada com sucesso",
  "trace_id": "abc123",
  "span_id": "def456",
  "service": "reservation-service",
  "business": {
    "domain": "reservation",
    "operation": "create_reservation",
    "reservation_id": "res_789",
    "amount": 2500.00,
    "currency": "BRL"
  },
  "context": {
    "user_id": "usr_sha256_...",
    "session_id": "sess_xyz..."
  }
}
```

## Alertas (SLOs)

| SLO | Condição | Severidade |
|-----|----------|------------|
| Latência de busca | P99 > 2s por 5min | warning |
| Latência de busca | P99 > 5s por 5min | critical |
| Taxa de erro | Error rate > 1% por 5min | warning |
| Taxa de erro | Error rate > 5% por 5min | critical |
| Disponibilidade | < 99.9% por 1min | critical |
| Cache hit ratio | < 70% por 10min | warning |
| Pagamentos falhos | > 5% por 5min | critical |

## Dashboards Obrigatórios

1. **Business Overview**: KPIs de negócio (conversão, tempo de checkout, revenue)
2. **Technical Overview**: Latência, throughput, error rate por serviço
3. **Cache Performance**: Hit/miss ratio, evictions, memory usage
4. **Payment Health**: Taxa de aprovação, erros por gateway, tempo de processamento

## Rastreabilidade de Erros

Todos os erros devem incluir:
- `error.type`: Classificação (validation, external_service, infrastructure)
- `error.message`: Mensagem sanitizada (sem PII)
- `error.stack`: Stack trace (somente em dev/staging)
- `error.retryable`: Se pode ser tentado novamente

## PII e Segurança

- **NUNCA** logar: números de cartão, CVV, senhas, documentos completos
- **Hash** user_id em logs: `usr_sha256_${hash}`
- **Mascarar** emails: `leandro.yam***@email.com`
- Usar campos dedicados para dados sensíveis que precisam ser rastreados
