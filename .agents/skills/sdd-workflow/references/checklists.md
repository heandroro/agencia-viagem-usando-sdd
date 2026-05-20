# SDD Workflow — Checklists de Cada Etapa

## Checklist — Requirements

```
[ ] Todos os requisitos têm ID REQ-x.x
[ ] Nenhum detalhe técnico presente (tecnologias, bibliotecas, arquitetura)
[ ] Seções OBS-x presentes descrevendo o QUE observar (não o como)
[ ] Critérios de aceite mensuráveis para cada requisito funcional
[ ] Linguagem funcional — sem menção a implementação
[ ] Aprovação humana obtida antes de avançar para design.md
```

---

## Checklist — Design

```
[ ] requirements.md aprovado pelo humano antes de iniciar
[ ] Alinhamento com architecture.md verificado — sem conflitos
[ ] Seção APM / Observability Design presente e preenchida
[ ] Para cada OBS-x em requirements.md, existe APM-Mx ou APM-Ex correspondente
[ ] IDs de telemetria seguem o padrão: APM-Mx (métricas), APM-Ex (eventos)
[ ] Contratos de interface documentados (APIs, eventos, SDKs)
[ ] Nenhuma dependência externa não listada em architecture.md
[ ] Aprovação humana obtida antes de avançar para tasks.md
```

---

## Checklist — Tasks

```
[ ] design.md aprovado pelo humano antes de iniciar
[ ] Cada task referencia pelo menos um [REQ-x.x]
[ ] Tasks T-APM-01 a T-APM-05 presentes e NÃO marcadas como opcionais
[ ] Tasks de APM referenciam IDs [APM-Mx] / [APM-Ex] do design.md
[ ] Nenhuma task parece grande demais para uma sessão (subdividir se necessário)
[ ] Checklist de promoção para architecture.md presente ao final do arquivo
```

---

## Checklist — Execução (por task)

```
[ ] Task confirmada com o humano antes de iniciar
[ ] Implementação completa e testada
[ ] Marcada [x] no tasks.md
[ ] Aprovação humana obtida antes de avançar para a próxima task
```

---

## Checklist Final (após todas as tasks concluídas)

```
[ ] Todas as tasks T-APM-xx concluídas e validadas
[ ] Telemetria visível na plataforma APM (métricas, eventos, traces)
[ ] Alertas configurados e testados
[ ] Dashboard atualizado (se aplicável)
[ ] Decisões arquiteturais relevantes promovidas para architecture.md
[ ] ADR criado se houver decisão arquitetural significativa (.sdd/adr/)
[ ] Spec descartada ou arquivada (spec-first)
[ ] CHANGELOG.md atualizado com a entrega
```

---

## Tasks APM obrigatórias (T-APM-01 a T-APM-05)

Estas tasks **nunca são opcionais** — devem aparecer em todo `tasks.md`:

| Task | Descrição |
|---|---|
| T-APM-01 | Implementar métricas customizadas (APM-Mx do design.md) |
| T-APM-02 | Implementar eventos de negócio (APM-Ex do design.md) |
| T-APM-03 | Implementar traces distribuídos (spans obrigatórios) |
| T-APM-04 | Configurar alertas (condições e severidades do design.md) |
| T-APM-05 | Validar telemetria em staging (métricas visíveis, alertas disparando) |
