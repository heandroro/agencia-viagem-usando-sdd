---
description: Inicializa o memory bank do projeto (.sdd/memory-bank/) com os quatro arquivos obrigatórios preenchidos de acordo com o contexto do projeto fornecido pelo humano. Use em projetos novos ou quando o memory bank ainda não existe.
argument-hint: "descrição do projeto em 1-3 frases (ex: API de pagamentos para e-commerce B2B)"
input:
  - projeto: "Descrição do projeto em 1-3 frases"
allowed-tools: [Read, Write]
---

Você é um arquiteto SDD. Inicialize o memory bank do projeto com base na
descrição fornecida.

**Projeto**: ${input:projeto}

## Contexto

O memory bank é composto por quatro arquivos obrigatórios em `.sdd/memory-bank/`:

| Arquivo | Propósito |
|---|---|
| `constitution.md` | Princípios imutáveis — **nunca personalizado por projeto** |
| `architecture.md` | Decisões arquiteturais, estilo, dependências, C4 |
| `product.md` | Visão, personas, KPIs, restrições de negócio |
| `apm-standards.md` | Padrões de observabilidade adotados |

## Passos

1. Verifique se `.sdd/memory-bank/` já existe com arquivos presentes.
   Se existir, **pare imediatamente** e informe ao humano — não sobrescreva
   um memory bank existente.

2. **Faça perguntas ao humano** antes de criar qualquer arquivo.
   Colete as seguintes informações:

   **Para `architecture.md`**:
   - Estilo arquitetural (microsserviços, monolito, event-driven, etc.)
   - Linguagens e plataformas principais
   - Dependências externas já conhecidas (banco de dados, brokers, gateways)

   **Para `product.md`**:
   - Personas principais e suas necessidades
   - 3-5 objetivos de negócio
   - KPIs de negócio que devem ser observáveis via APM
   - Restrições conhecidas (compliance, regulatório, contrato)

   **Para `apm-standards.md`**:
   - Plataforma APM a usar (ex: Azure Monitor, Datadog, OpenTelemetry)
   - Campos adicionais obrigatórios além dos padrão SDD
   - Ferramentas de alerta e dashboards planejadas

3. Com as respostas em mãos, crie os quatro arquivos:

   - `constitution.md` — copie o conteúdo padrão SDD **sem nenhuma alteração**;
     este arquivo é idêntico em todos os projetos
   - `architecture.md` — preencha com as informações coletadas; use `[PREENCHER]`
     apenas onde o humano não forneceu informação
   - `product.md` — preencha com as informações coletadas
   - `apm-standards.md` — preencha adaptando os padrões à plataforma informada,
     mantendo os campos obrigatórios padrão SDD

4. Apresente um resumo do que foi criado e oriente o próximo passo:
   > "Memory bank inicializado. Você pode agora usar `/criar-spec` para iniciar
   > o primeiro spec do projeto."

## Restrições

- **Nunca sobrescreva** um memory bank existente — apenas `/atualizar-produto` e
  `/atualizar-apm-standards` fazem atualizações pontuais após a inicialização
- `constitution.md` deve ser copiado **sem modificações** — não personalize
- Não invente informações arquiteturais — use `[PREENCHER]` onde o humano não respondeu
- Não inclua tokens ou secrets em `apm-standards.md` — use `${VAR}`
