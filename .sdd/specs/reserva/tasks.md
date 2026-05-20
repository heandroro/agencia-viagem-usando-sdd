# [NOME DA FEATURE] — Tasks

## Tasks de Implementação

### T-01: [Descrição da task]
**Referência**: [REQ-x.x]

**Passos**:
1. [Passo 1]
2. [Passo 2]
3. [Passo 3]

**Checklist de conclusão**:
- [ ] Código implementado
- [ ] Testes unitários
- [ ] Testes de integração (se aplicável)
- [ ] Documentação atualizada

---

### T-02: [Descrição da task]
**Referência**: [REQ-x.x]

**Passos**:
1. [Passo 1]
2. [Passo 2]

---

## Tasks APM (Obrigatórias)

### T-APM-01: Instrumentar métricas de negócio
**Referência**: [APM-M1, APM-M2...]

- [ ] Métrica APM-M1 implementada
- [ ] Métrica APM-M2 implementada
- [ ] Dashboard configurado

### T-APM-02: Instrumentar logging estruturado
**Referência**: [APM-E1, APM-E2...]

- [ ] Eventos APM-E1 logados
- [ ] Atributos obrigatórios presentes (business.domain, operation, outcome)
- [ ] Validação: nenhum log com PII

### T-APM-03: Instrumentar traces distribuídos
- [ ] Spans criados para operações principais
- [ ] Contexto propagado entre serviços
- [ ] Atributos de trace obrigatórios presentes

### T-APM-04: Configurar alertas
- [ ] Alerta de latência configurado
- [ ] Alerta de erro configurado
- [ ] Runbooks vinculados aos alertas

### T-APM-05: Validar observabilidade
- [ ] Dashboards verificados em staging
- [ ] Logs estruturados validados (Loki/ELK)
- [ ] Métricas aparecem no Prometheus
- [ ] Traces completos no Tempo/Jaeger

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
**Status**: [pending | in_progress | completed]
