FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Imagem final
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binário
COPY --from=builder /app/server .

# Expor porta
EXPOSE 8080

# Executar
CMD ["./server"]
