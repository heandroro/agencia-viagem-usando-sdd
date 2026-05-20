# Memory Bank — Guia de Atualização por Arquivo

---

## constitution.md

**Propósito**: Define os princípios imutáveis do projeto — observabilidade obrigatória,
separação funcional × técnica, spec-first, aprovação humana obrigatória entre etapas.

**Quem pode alterar**: Somente o humano, com decisão consciente e explícita.

**IA**: Nunca edite. Use apenas como fonte de verdade para validar specs e tasks.
Se detectar violação de um princípio, reporte — não corrija silenciosamente.

**Seções obrigatórias**:
- Princípios numerados (ex: `## 1. Princípio da Observabilidade Obrigatória`)
- Cada princípio com declaração e bullets de aplicação

---

## architecture.md

**Propósito**: Registra as decisões arquiteturais vigentes — tecnologias adotadas,
padrões de projeto, dependências aprovadas e ADRs.

**Quem pode alterar**: IA, mas somente via `/promover-adr` e somente após conclusão
de uma task que gerou uma decisão arquitetural relevante.

**Quando atualizar**:
- Ao final de uma task que adotou um novo padrão, biblioteca ou abordagem técnica
- Nunca no início ou meio de uma feature — apenas na conclusão

**Checklist de promoção** (verificar antes de atualizar):
- [ ] A decisão foi tomada durante execução de uma task (não especulativa)
- [ ] O ADR foi criado em `.sdd/adr/ADR-NN-titulo.md`
- [ ] A referência ao ADR foi adicionada à seção de ADRs em `architecture.md`
- [ ] O humano foi informado sobre a promoção

**Seções esperadas**: dependências aprovadas, padrões de projeto, lista de ADRs.

---

## product.md

**Propósito**: Descreve o produto em alto nível — visão, público-alvo, restrições
de negócio, segmentos e contexto estratégico. Usado para orientar a formulação de
requisitos funcionais em `requirements.md`.

**Quem pode alterar**: IA via `/atualizar-produto`, **com aprovação humana explícita**
antes de salvar.

**Quando atualizar**:
- Nova visão ou direção estratégica do produto
- Novo segmento de cliente ou caso de uso principal
- Restrição de negócio relevante que impacta futuros specs
- **Não** atualizar como efeito colateral de uma feature — apenas por decisão deliberada

**Fluxo obrigatório**:
1. Ler `product.md` atual
2. Propor o diff ao humano
3. Aguardar aprovação explícita
4. Aplicar apenas o que foi aprovado

**Seções esperadas**: Visão do Produto, Público-Alvo, Restrições de Negócio,
Contexto Estratégico.

---

## apm-standards.md

**Propósito**: Define os padrões de observabilidade e APM CLI do projeto — pilares
adotados (logs, métricas, traces), campos obrigatórios em spans, ferramentas APM,
convenções de naming e alertas.

**Quem pode alterar**: IA via `/atualizar-apm-standards`, **com aprovação humana explícita**
antes de salvar.

**Quando atualizar**:
- Adoção de nova ferramenta ou plataforma APM
- Novo campo obrigatório em spans ou métricas
- Mudança nas convenções de naming de eventos ou traces
- Novo pilar ou padrão de alerta adotado pelo time
- **Não** atualizar por feature individual — apenas quando o padrão muda globalmente

**Fluxo obrigatório**:
1. Ler `apm-standards.md` atual
2. Ler `constitution.md` para verificar conformidade com Princípio da Observabilidade Obrigatória
3. Propor o diff ao humano
4. Aguardar aprovação explícita
5. Aplicar apenas o que foi aprovado

**Seções esperadas**: Modelo de Observabilidade Adotado, Campos Obrigatórios,
Ferramentas APM, Convenções de Naming, Padrões de Alerta.
