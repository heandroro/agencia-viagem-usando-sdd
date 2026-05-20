---
description: Analisa a estrutura de um repositório GitHub e registra-o como referência de projeto em .sdd/references/<nome>.md para orientar decisões de design, geração de tasks e código. Requer o MCP do GitHub configurado (io.github.github/github-mcp-server).
argument-hint: "owner/repo do repositório de referência (ex: nestjs/nest)"
input:
  - repo: "Repositório GitHub no formato owner/repo (ex: nestjs/nest)"
  - nome: "Nome curto para identificar a referência (ex: nestjs, fastapi, nextjs)"
allowed-tools: [Read, Write, Edit, mcp_github_get_file_contents, mcp_github_search_code]
---

Você é um arquiteto SDD. Analise o repositório de referência e registre-o como
padrão de projeto para orientar specs e tasks.

**Repositório**: ${input:repo}
**Nome da referência**: ${input:nome}

## Aviso de licença

Antes de qualquer ação, leia o arquivo `LICENSE` ou `LICENSE.md` do repositório.
Se a licença **não for** MIT, Apache-2.0, ISC, BSD-2-Clause, BSD-3-Clause ou similar
permissiva, **pare imediatamente**, informe ao humano a licença encontrada e aguarde
confirmação explícita de que tem permissão para usar o repositório como referência.

## Passos

1. Verifique se `.sdd/references/${input:nome}.md` já existe. Se existir,
   informe ao humano e pergunte se deseja atualizar ou cancelar.

2. Leia o arquivo `LICENSE` do repositório via MCP para verificar a licença antes
   de continuar.

3. Com a licença confirmada, analise a estrutura do repositório:
   - Arquivo raiz de configuração (`package.json`, `pyproject.toml`, `go.mod`,
     `pom.xml`, `Cargo.toml`) para identificar stack e dependências principais
   - Estrutura de pastas de nível 1 e 2 via `get_file_contents` no diretório raiz
   - Arquivos de configuração relevantes: `Dockerfile`, `.github/workflows/`,
     arquivos de lint e test
   - Padrões de nomenclatura e organização de módulos

4. Use `search_code` para extrair padrões de código relevantes para o projeto:
   - Como a aplicação é inicializada (ex: `main.ts`, `app.py`, `main.go`)
   - Como erros são tratados e logados
   - Como testes são estruturados e nomeados

5. Crie `.sdd/references/${input:nome}.md` com o formato:

   ```
   # Referência de Projeto: <nome>

   > Repositório: github.com/${input:repo}
   > Licença: <licença>
   > Registrado em: <data de hoje>

   ## Stack e Dependências Principais
   <linguagem, runtime, framework principal, dependências-chave>

   ## Estrutura de Pastas
   <árvore de nível 1-2 com descrição de cada pasta>

   ## Padrões de Código
   <inicialização, tratamento de erro, logging, organização de testes>

   ## Convenções de Nomenclatura
   <arquivos, módulos, funções, variáveis>

   ## Configurações de Referência
   <trechos relevantes de Dockerfile, CI/CD, lint — sem dados sensíveis>
   ```

6. Leia `.sdd/memory-bank/architecture.md` e adicione a entrada na seção de
   referências (criando a seção "Repositórios de Referência" se não existir):
   ```
   - **${input:nome}** — `.sdd/references/${input:nome}.md`
     (baseado em github.com/${input:repo})
   ```

7. Confirme ao humano o arquivo criado e o caminho.

## Restrições

- **Nunca** inclua tokens, senhas ou segredos encontrados no repositório
- Não copie blocos de código diretamente — extraia padrões, convenções e estruturas
- Limite o arquivo de referência a ~200 linhas — foco no que orienta decisões de design
- Se o repositório for privado e o MCP não tiver acesso, informe ao humano
- Se o MCP do GitHub não estiver configurado, informe:
  > "Este comando requer o MCP do GitHub. Configure-o com:
  > `apm install sdd-framework` (já inclui a dependência) ou adicione
  > `io.github.github/github-mcp-server` manualmente ao seu editor."
