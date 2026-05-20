package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/heandroro/agencia-viagem/backend/internal/domain/reservation"
	"github.com/heandroro/agencia-viagem/backend/pkg/crypto"
)

// UpdateTravelersUseCase caso de uso para atualizar viajantes
type UpdateTravelersUseCase struct {
	reservationRepo reservation.Repository
	crypto          crypto.Service
}

// UpdateTravelersInput dados de entrada para atualização de viajantes
type UpdateTravelersInput struct {
	ReservationID string
	UserID        string
	Travelers     []TravelerInput
}

// TravelerInput dados de um viajante
type TravelerInput struct {
	Type           reservation.TravelerType
	FullName       string
	DocumentType   reservation.DocumentType
	DocumentNumber string
	BirthDate      time.Time
}

// UpdateTravelersOutput dados de saída da atualização
type UpdateTravelersOutput struct {
	ReservationID string
	Status        reservation.Status
	Travelers     []TravelerOutput
}

// TravelerOutput viajante na resposta (com documento mascarado)
type TravelerOutput struct {
	Type           string `json:"type"`
	FullName       string `json:"full_name"`
	DocumentType   string `json:"document_type"`
	DocumentNumber string `json:"document_number"` // Mascarado
	BirthDate      string `json:"birth_date"`
}

// NewUpdateTravelersUseCase cria um novo caso de uso
func NewUpdateTravelersUseCase(
	reservationRepo reservation.Repository,
	crypto crypto.Service,
) *UpdateTravelersUseCase {
	return &UpdateTravelersUseCase{
		reservationRepo: reservationRepo,
		crypto:          crypto,
	}
}

// Execute executa o caso de uso
func (uc *UpdateTravelersUseCase) Execute(ctx context.Context, input UpdateTravelersInput) (*UpdateTravelersOutput, error) {
	// Buscar reserva
	rsv, err := uc.reservationRepo.GetByID(ctx, input.ReservationID)
	if err != nil {
		if err.Error() == "reservation_not_found" {
			return nil, errors.New("reservation_not_found")
		}
		return nil, err
	}

	// Validar que reserva pertence ao usuário
	if rsv.UserID != input.UserID {
		return nil, errors.New("unauthorized")
	}

	// Validar que reserva está em status "pending"
	if rsv.Status != reservation.StatusPending {
		return nil, errors.New("reservation_not_pending")
	}

	// Converter input para domain.Traveler com criptografia
	travelers := make([]reservation.Traveler, len(input.Travelers))
	for i, t := range input.Travelers {
		// Criptografar documento
		encryptedDoc, err := uc.crypto.Encrypt(t.DocumentNumber)
		if err != nil {
			return nil, err
		}

		// Calcular hash do documento
		docHash := uc.crypto.Hash(t.DocumentNumber)

		travelers[i] = reservation.Traveler{
			Type:              t.Type,
			FullName:          t.FullName,
			DocumentType:      t.DocumentType,
			DocumentEncrypted: encryptedDoc,
			DocumentHash:      docHash,
			BirthDate:         t.BirthDate,
		}
	}

	// Atualizar viajantes no repositório
	if err := uc.reservationRepo.UpdateTravelers(ctx, input.ReservationID, travelers); err != nil {
		return nil, err
	}

	// Montar resposta (com documentos mascarados)
	output := &UpdateTravelersOutput{
		ReservationID: rsv.ID.String(),
		Status:        rsv.Status,
		Travelers:     make([]TravelerOutput, len(travelers)),
	}

	for i, t := range travelers {
		output.Travelers[i] = TravelerOutput{
			Type:           string(t.Type),
			FullName:       t.FullName,
			DocumentType:   string(t.DocumentType),
			DocumentNumber: reservation.MaskDocument(t.DocumentEncrypted), // Mascarado
			BirthDate:      t.BirthDate.Format("2006-01-02"),
		}
	}

	return output, nil
}
