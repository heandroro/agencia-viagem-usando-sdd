# Product — Contexto de Produto

## Visão

Plataforma de agência de viagem digital que simplifica a busca, comparação e reserva de pacotes de viagem (voo + hotel), oferecendo experiência fluida e transparente de checkout.

## Personas

### Persona 1: Viajante (Cliente)
- **Necessidades**: Buscar pacotes por destino/datas, comparar preços, reservar com segurança, acompanhar reservas
- **Dores**: Processo de reserva complexo, falta de transparência de preços, medo de problemas na viagem
- **Objetivos**: Encontrar o melhor pacote pelo melhor preço, reservar com confiança

### Persona 2: Administrador (Operador da Agência)
- **Necessidades**: CRUD de pacotes no catálogo, visualizar reservas, gerenciar disponibilidade
- **Dores**: Atualização manual de preços, falta de visibilidade de conversão
- **Objetivos**: Manter catálogo atualizado, maximizar taxa de conversão

## Objetivos de Negócio

1. **Taxa de conversão**: 5% de buscas → reservas confirmadas no MVP
2. **Tempo de checkout**: Média de < 5 minutos do início da busca à confirmação
3. **Satisfação do cliente**: NPS > 40 após primeira viagem
4. **Disponibilidade do catálogo**: Cache hit > 80% para buscas populares

## KPIs de Negócio (Observáveis via APM)

| KPI | Métrica Técnica | Target |
|-----|-----------------|--------|
| Conversão | `business.reservation.completed / business.search.performed` | > 5% |
| Tempo de checkout | Histograma `checkout.duration` (P50, P95, P99) | P95 < 5min |
| Abandono de carrinho | `business.cart.abandoned` vs `business.cart.initiated` | < 60% |
| Disponibilidade catálogo | `cache.hit_ratio.catalog` | > 80% |
| Erros de pagamento | `business.payment.failed` | < 2% |
| Latência de busca | `search.response_time` P99 | < 2s |

## Jornada do Usuário

```
Busca → Resultados → Detalhes → Reserva → Pagamento → Confirmação
  │         │           │          │          │            │
  │         │           │          │          │            └── Email de confirmação
  │         │           │          │          └── Integração Stripe
  │         │           │          └── Form dados viajantes
  │         │           └── Fotos, preço, disponibilidade
  │         └── Lista de pacotes com filtros
  └── Origem, destino, datas, viajantes
```

## Restrições de Negócio

### Compliance
- **LGPD**: Consentimento explícito para dados pessoais, direito ao esquecimento
- **PCI-DSS**: Não armazenar dados de cartão (usar Stripe com tokenização)
- **CDC**: Transparência em preços, condições de cancelamento claras

### Regras de Domínio
- Reserva só é confirmada após pagamento aprovado
- Preços de pacotes podem mudar a cada 1 hora (TTL do cache)
- Cancelamento permitido até 48h antes do voo com taxa de 10%
- Voucher emitido automaticamente após confirmação

## Roadmap (MVP → V1 → V2)

### MVP (Atual)
- Busca simples de pacotes
- Reserva com pagamento
- "Minhas reservas"
- Admin básico (CRUD catálogo)

### V1 (Futuro)
- Busca avançada (filtros, ordenação)
- Sistema de reviews
- Notificações push/email
- Dashboard de analytics para admin

### V2 (Futuro)
- Integração com GDS (Amadeus/Sabre) para voos reais
- Recomendações ML
- Programa de fidelidade
- App mobile

## Glossário de Domínio

| Termo | Definição |
|-------|-----------|
| Pacote | Combinação de voo + hotel com preço único |
| Reserva | Solicitação de compra de um pacote por um cliente |
| Viajante | Pessoa que realizará a viagem (pode ser diferente do comprador) |
| Voucher | Documento de confirmação enviado após reserva confirmada |
| Catálogo | Conjunto de pacotes disponíveis para venda |
| Disponibilidade | Indica se o pacote tem vagas para as datas solicitadas |
