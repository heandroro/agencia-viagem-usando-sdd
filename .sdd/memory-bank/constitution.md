# Constitution — Princípios Fundamentais SDD

## 1. Observabilidade Obrigatória

Todo código em produção deve ser instrumentado. Tasks T-APM-xx nunca são opcionais.
Métricas, logs e traces são parte da definição de pronto (Definition of Done).

## 2. Separação Funcional × Técnica

- `requirements.md` = O QUÊ (o que o sistema deve fazer)
- `design.md` = O COMO (como o sistema será implementado)

Nunca misture detalhes técnicos em requisitos funcionais.

## 3. Contrato de Interface

Interfaces documentadas em `design.md` antes de qualquer código.
Alterações em contratos estabelecidos requerem aprovação explícita.

## 4. Rastreabilidade

- Cada task referencia `[REQ-x.x]` do requisito correspondente
- Cada task APM referencia `[APM-Mx]` (métrica) ou `[APM-Ex]` (evento)

## 5. Tamanho de Task

Uma task = uma sessão de trabalho. Subdividir se parecer grande demais.

## 6. Revisão Humana

Aprovação explícita entre cada etapa:
```
requirements.md → [aprovação] → design.md → [aprovação] → tasks.md → execução
```
Nunca avance sozinho.

## 7. Consistência com Memory Bank

`design.md` sempre alinhado com `architecture.md`.
Não introduza dependências externas não listadas em `architecture.md`.

## 8. Commit Semântico

Todos os commits seguem [Conventional Commits](https://www.conventionalcommits.org/pt-br/):
- `feat:` nova funcionalidade
- `fix:` correção de bug
- `docs:` documentação
- `refactor:` refatoração
- `test:` testes
- `chore:` manutenção
