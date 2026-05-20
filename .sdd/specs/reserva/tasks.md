# Reserva de Pacotes — Tasks

## Tasks de Implementação

### T-01: Criar modelo de dados MongoDB (reservations + availability)
**Referência**: [REQ-T.3]

**Passos**:
1. Criar collection `reservations` com schema validado
2. Criar collection `availability` com version para optimistic locking
3. Criar indexes: TTL em `expires_at`, compound em `user_id + status`
4. Implementar criptografia AES-256 para documentos em `travelers.document_encrypted`

**Checklist de conclusão**:
- [x] Schema de reservations implementado
- [x] Schema de availability implementado
- [x] Indexes criados
- [x] Testes de criptografia/descriptografia
- [x] Documentação de schema

---

### T-02: Implementar endpoint POST /api/v1/reservations
**Referência**: [REQ-1.1, APM-M1, APM-E1]

**Passos**:
1. BFF (FastAPI): Validação de sessão e payload
2. Backend (Go): Verificar disponibilidade via optimistic locking
3. Backend: Calcular preço (snapshot) e criar reserva em status "pending"
4. Backend: Publicar evento `reservation.created` no Valkey Stream
5. Configurar TTL de 30 minutos (MongoDB TTL index)

**Checklist de conclusão**:
- [x] Endpoint implementado no Backend
- [x] Endpoint implementado no BFF
- [x] Validação de disponibilidade funcionando (optimistic locking)
- [x] Cálculo de preço correto
- [x] Testes unitários (Backend usecase) ✅ PASS
- [x] Código compila sem erros
- [ ] Evento publicado no stream (Valkey) - será feito no worker T-05
- [ ] Testes de integração end-to-end

---

### T-03: Implementar endpoint PUT /api/v1/reservations/{id}/travelers
**Referência**: [REQ-1.2, APM-M3, APM-E3]

**Passos**:
1. BFF: Validar formato de CPF/passaporte (biblioteca de validação)
2. Backend: Atualizar reserva com dados dos viajantes (criptografados)
3. Implementar máscara de documento para resposta
4. Validar que reserva está em status "pending"

**Checklist de conclusão**:
- [ ] Validação de CPF/passaporte implementada
- [ ] Criptografia de documentos funcionando
- [ ] Máscara de documento aplicada na resposta
- [ ] Validação de status "pending" implementada
- [ ] Testes unitários
- [ ] Testes de integração

---

### T-04: Implementar endpoint GET /api/v1/reservations/{id}/summary
**Referência**: [REQ-1.3, APM-M4]

**Passos**:
1. Agregar dados do pacote (do cache ou MongoDB)
2. Montar resposta com breakdown de preços
3. Incluir políticas de cancelamento
4. Adicionar flag `policies_accepted` no modelo

**Checklist de conclusão**:
- [ ] Endpoint implementado no BFF
- [ ] Agregação de dados funcionando
- [ ] Breakdown de preços correto
- [ ] Políticas incluídas na resposta
- [ ] Testes unitários
- [ ] Testes de integração

---

### T-05: Implementar worker de expiração de reservas
**Referência**: [REQ-1.4, APM-M5, APM-E5]

**Passos**:
1. Criar worker Go que consome de Valkey Stream `reservation.expired`
2. Implementar lógica de atualização de status para "expired"
3. Liberar slots de disponibilidade
4. Configurar retry e dead letter queue

**Checklist de conclusão**:
- [ ] Worker consumindo stream corretamente
- [ ] Expiração atualizando status no MongoDB
- [ ] Slots liberados após expiração
- [ ] Retry e DLQ configurados
- [ ] Testes unitários
- [ ] Testes de integração

---

### T-06: Implementar validações e tratamento de erros
**Referência**: [REQ-T.1, REQ-T.2]

**Passos**:
1. Implementar validação de campos obrigatórios
2. Implementar rate limiting no BFF
3. Tratar conflitos de disponibilidade (409 Conflict)
4. Sanitizar logs (nenhum PII em texto claro)

**Checklist de conclusão**:
- [ ] Validações de campos implementadas
- [ ] Rate limiting configurado
- [ ] Tratamento de conflitos implementado
- [ ] Logs sanitizados (sem PII)
- [ ] Testes de edge cases

---

## Tasks APM (Obrigatórias)

### T-APM-01: Instrumentar métricas de negócio
**Referência**: [APM-M1, APM-M2, APM-M3, APM-M4, APM-M5, APM-M6, APM-M7, APM-M8, APM-M9]

- [ ] APM-M1 `business.reservation.created` (counter + package_id, destination)
- [ ] APM-M2 `business.reservation.failed` (counter + reason)
- [ ] APM-M3 `business.reservation.travelers.captured` (counter + count)
- [ ] APM-M4 `business.reservation.summary_viewed` (counter)
- [ ] APM-M5 `business.reservation.expired` (counter + reason)
- [ ] APM-M6 `reservation.api.duration` (histogram + endpoint, method)
- [ ] APM-M7 `reservation.db.duration` (histogram + operation)
- [ ] APM-M8 `reservation.cache.hit` (counter + cache_name)
- [ ] APM-M9 `reservation.cache.miss` (counter + cache_name)
- [ ] Dashboard "Reservations Overview" criado no Grafana

### T-APM-02: Instrumentar logging estruturado
**Referência**: [APM-E1, APM-E2, APM-E3, APM-E4, APM-E5]

- [ ] APM-E1 `reservation.created` logado em INFO
- [ ] APM-E2 `reservation.create_failed` logado em WARN
- [ ] APM-E3 `reservation.travelers_updated` logado em INFO
- [ ] APM-E4 `reservation.validation_failed` logado em WARN
- [ ] APM-E5 `reservation.expired` logado em INFO
- [ ] Atributos obrigatórios em todos os logs: `business.domain`, `business.operation`, `business.outcome`
- [ ] Validação: nenhum log com documento em texto claro (somente hash ou máscara)
- [ ] Correlation ID propagado entre todos os logs

### T-APM-03: Instrumentar traces distribuídos
**Referência**: [Traces definidos em design.md]

- [ ] Span `initiate_reservation` com atributos: package_id, destination, traveler_count, total_amount
- [ ] Span `capture_traveler_data` com atributos: reservation_id, traveler_count, validation_time_ms
- [ ] Span `generate_reservation_summary` com atributos: reservation_id, generation_time_ms
- [ ] Span `expire_reservation` com atributos: reservation_id, elapsed_minutes, reason
- [ ] Span `check_availability` com latência de consulta
- [ ] Span `encrypt_traveler_data` com tempo de criptografia
- [ ] Contexto de trace propagado: Frontend → BFF → Backend → MongoDB/Valkey
- [ ] Atributos obrigatórios: `service.name`, `service.version`, `deployment.environment`, `trace.trace_id`

### T-APM-04: Configurar alertas
**Referência**: [apm-standards.md - SLOs]

- [ ] Alerta: Latência P99 > 500ms por 5min (warning)
- [ ] Alerta: Latência P99 > 2s por 5min (critical)
- [ ] Alerta: Error rate > 1% por 5min (warning)
- [ ] Alerta: Error rate > 5% por 5min (critical)
- [ ] Alerta: `business.reservation.failed` > 10% em 1h (critical)
- [ ] Runbook "Reserva - Alta Latência" vinculado
- [ ] Runbook "Reserva - Erros de Disponibilidade" vinculado

### T-APM-05: Validar observabilidade
**Referência**: [apm-standards.md]

- [ ] Dashboard "Reservations Business Overview" verificado em staging
- [ ] Dashboard "Reservations Technical Overview" verificado em staging
- [ ] Logs estruturados aparecendo no Loki
- [ ] Métricas `business.reservation.*` aparecendo no Prometheus
- [ ] Traces completos (Frontend → BFF → Backend → DB) no Tempo/Jaeger
- [ ] Alertas testados (firing test)
- [ ] Validação: Nenhum log contém documento ou email em texto claro

---

## Checklist Final

### Antes de merge
- [ ] Todas as tasks implementadas
- [ ] T-APM-01 a T-APM-05 concluídas
- [ ] CI passando (lint, testes, build)
- [ ] Code review aprovado

### Promoção para architecture.md (se necessário)
- [ ] Alguma decisão arquitetural nova surgiu durante implementação?
- [ ] Se sim, usar `/promover-adr` para registrar

---
**Status**: pending

## Ordem de Execução Sugerida

1. **T-01**: Modelo de dados (base para todas as outras)
2. **T-APM-01 a T-APM-05**: Instrumentação APM (fazer junto com as tasks de feature)
3. **T-02**: Criar reserva (core do fluxo)
4. **T-03**: Capturar viajantes
5. **T-04**: Resumo da reserva
6. **T-05**: Worker de expiração
7. **T-06**: Validações e edge cases

Cada task deve ser executada **uma por vez**, com code review e aprovação antes de prosseguir.
