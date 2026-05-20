# Reserva de Pacotes — Requirements

## Contexto

O usuário (Viajante) precisa reservar um pacote de viagem (voo + hotel) após encontrá-lo na busca. O fluxo de reserva deve garantir disponibilidade, capturar dados dos viajantes, calcular o preço total e preparar para o pagamento.

## Requisitos Funcionais

### REQ-1.1 Iniciar reserva de pacote
**Descrição**: O sistema deve permitir que o usuário inicie uma reserva a partir de um pacote selecionado, verificando disponibilidade para as datas escolhidas.

**Critérios de Aceite**:
- [ ] Usuário pode selecionar um pacote e datas de viagem
- [ ] Sistema valida disponibilidade do pacote para as datas
- [ ] Sistema calcula preço total (pacote × número de viajantes × dias)
- [ ] Reserva é criada em status "pendente" com TTL de 30 minutos

**OBS-1.1** — Observabilidade
- **Métrica**: `business.reservation.created` (counter + package_id, destination)
- **Log**: Evento `reservation.created` com reservation_id, user_id, amount
- **Trace**: Span `initiate_reservation` com atributos: package_id, dates, traveler_count

### REQ-1.2 Capturar dados dos viajantes
**Descrição**: O sistema deve coletar e validar dados de todos os viajantes (nome completo, documento, data de nascimento) antes de permitir avançar.

**Critérios de Aceite**:
- [ ] Formulário captura dados de cada viajante (comprador + acompanhantes)
- [ ] Validação de campos obrigatórios
- [ ] Validação de formato de documento (CPF/passaporte)
- [ ] Dados são persistidos vinculados à reserva

**OBS-1.2** — Observabilidade
- **Métrica**: `business.reservation.travelers.captured` (counter + count)
- **Log**: Evento `reservation.travelers_updated` com reservation_id, traveler_count
- **Trace**: Span `capture_traveler_data` com validação de campos

### REQ-1.3 Calcular e exibir resumo da reserva
**Descrição**: O sistema deve exibir um resumo detalhado da reserva antes da confirmação, incluindo detalhes do pacote, viajantes, preços e condições.

**Critérios de Aceite**:
- [ ] Exibe detalhes do pacote (hotel, voo, datas)
- [ ] Lista todos os viajantes com dados cadastrados
- [ ] Mostra breakdown de preços (pacote, taxas, total)
- [ ] Apresenta política de cancelamento
- [ ] Requer aceite dos termos para prosseguir

**OBS-1.3** — Observabilidade
- **Métrica**: `business.reservation.summary_viewed` (counter)
- **Log**: Evento `reservation.summary_viewed` com reservation_id, total_amount
- **Trace**: Span `generate_reservation_summary` com tempo de geração

### REQ-1.4 Gerenciar expiração de reserva pendente
**Descrição**: Reservas em status "pendente" devem expirar automaticamente após 30 minutos sem conclusão do pagamento.

**Critérios de Aceite**:
- [ ] Reserva pendente expira após 30 minutos
- [ ] Sistema libera disponibilidade do pacote ao expirar
- [ ] Usuário é notificado sobre expiração
- [ ] Reserva expirada não pode ser paga

**OBS-1.4** — Observabilidade
- **Métrica**: `business.reservation.expired` (counter + reason)
- **Log**: Evento `reservation.expired` com reservation_id, time_elapsed
- **Trace**: Span `expire_reservation` com motivo da expiração

## Requisitos Técnicos

### REQ-T.1 Performance
**Descrição**: A criação de reserva e cálculos devem ser rápidos para não frustrar o usuário.

**Critérios**:
- Latência P99 < 500ms para criação de reserva
- Latência P99 < 200ms para cálculo de preço
- Throughput > 100 req/s para criação de reserva

### REQ-T.2 Segurança
**Descrição**: Proteção de dados pessoais dos viajantes.

**Critérios**:
- Dados pessoais (CPF, passaporte) criptografados em repouso
- Nenhum log contém documentos em texto claro (hash ou máscara)
- Acesso à reserva apenas pelo usuário dono ou admin

### REQ-T.3 Consistência
**Descrição**: Prevenir overbooking através de controle de concorrência.

**Critérios**:
- Uso de pessimistic locking ou atomic operations para verificar disponibilidade
- Reserva só é criada se houver disponibilidade confirmada

## Dependências

- [x] Depende de: Catálogo de pacotes (pacote deve existir para reservar)
- [ ] Impacta em: Fluxo de pagamento (reserva precisa existir para pagar)
- [ ] Impacta em: "Minhas reservas" (listagem de reservas do usuário)

## Checklist de Validação

- [x] Todos os requisitos têm IDs únicos
- [x] Cada requisito funcional tem OBS-x correspondente
- [x] Não há detalhes de implementação técnica (banco, endpoints)
- [x] Critérios de aceite são verificáveis
- [x] Dependências estão mapeadas

---
**Status**: review
**Próximo passo**: Após aprovação, gerar design.md
