# Backend - Agência de Viagem

Serviço backend em Go para o sistema de reservas de agência de viagem.

## Estrutura

```
backend/
├── cmd/server/          # Entry point
├── internal/
│   ├── domain/
│   │   └── reservation/ # Domínio de reservas
│   │       ├── model.go        # Entidades
│   │       ├── repository.go   # Persistência
│   │       └── availability.go # Disponibilidade
│   ├── infra/
│   │   ├── database/    # MongoDB
│   │   └── config/      # Configurações
│   └── api/
│       └── handlers/    # HTTP handlers
├── pkg/
│   └── crypto/          # Criptografia AES-256
└── go.mod
```

## Tecnologias

- **Go 1.23**
- **Gin** - Web framework
- **MongoDB** - Database
- **Valkey** - Cache e streams
- **AES-256-GCM** - Criptografia de dados sensíveis

## Executar Localmente

```bash
# Setup infraestrutura
make setup

# Executar
make run-backend

# Testes
make test-backend
```

## Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `SERVER_PORT` | Porta do servidor | `8080` |
| `MONGODB_URI` | URI do MongoDB | `mongodb://localhost:27017` |
| `MONGODB_DATABASE` | Nome do database | `agencia_viagem` |
| `VALKEY_ADDR` | Endereço do Valkey | `localhost:6379` |
| `ENCRYPTION_KEY` | Chave AES-256 (mínimo 32 chars) | - |
| `JWT_SECRET` | Secret para JWT | - |

## API Endpoints

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/api/v1/reservations` | Criar reserva |
| PUT | `/api/v1/reservations/:id/travelers` | Atualizar viajantes |
| GET | `/api/v1/reservations/:id/summary` | Resumo da reserva |
| GET | `/health` | Health check |

## Segurança

- Documentos (CPF/passaporte) são criptografados com AES-256-GCM
- Campos sensíveis nunca são expostos na API (usam máscara)
- Hash dos documentos é usado para correlacionar em logs
