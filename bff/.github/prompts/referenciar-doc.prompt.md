---
description: Busca documentação de uma tecnologia, biblioteca, banco de dados ou plataforma a partir de uma URL e registra um resumo em .sdd/references/<tech>.md para uso em specs futuros. Registra a referência em architecture.md.
argument-hint: "URL da documentação (ex: https://docs.kafka.apache.org/quickstart)"
input:
  - url: "URL da documentação de referência"
  - tech: "Nome curto da tecnologia (ex: kafka, postgres, redis, stripe)"
allowed-tools: [Read, Fetch, Write, Edit]
---

Você é um arquiteto SDD. Registre a documentação de referência informada.

**URL**: ${input:url}
**Tecnologia**: ${input:tech}

## Passos

1. Verifique se `.sdd/references/${input:tech}.md` já existe. Se existir, informe
   ao humano e pergunte se deseja sobrescrever ou apenas complementar.

2. Acesse `${input:url}` e extraia apenas o que é relevante para o trabalho de
   specs SDD. Foque em:
   - Conceitos principais e terminologia
   - Tipos de dados, schemas e contratos de interface
   - Padrões de uso recomendados e anti-padrões conhecidos
   - Limites, cotas e restrições operacionais relevantes
   - Exemplos de configuração (sem valores sensíveis — use `${VAR}`)

3. Crie `.sdd/references/${input:tech}.md` com o seguinte formato:

   ```
   # Referência: <nome da tecnologia>

   > Fonte: <url>
   > Registrado em: <data de hoje>

   ## Conceitos principais
   ...

   ## Contratos e tipos
   ...

   ## Padrões de uso
   ...

   ## Restrições e limites
   ...
   ```

4. Leia `.sdd/memory-bank/architecture.md` e adicione (ou atualize) a entrada
   da tecnologia na seção de dependências, com referência ao arquivo:
   ```
   - **${input:tech}** — `.sdd/references/${input:tech}.md`
   ```

5. Confirme ao humano o que foi salvo e o caminho do arquivo.

## Restrições

- Nunca inclua tokens, senhas ou valores de configuração sensíveis — use `${VAR}`
- Não copie a documentação inteira — extraia apenas o que orienta decisões de design
- Se a URL retornar erro ou conteúdo insuficiente, informe ao humano e sugira
  uma URL alternativa
- O arquivo de referência é um resumo para orientar specs, não uma cópia da documentação oficial
