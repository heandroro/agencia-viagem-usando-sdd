package usecase

import (
	"context"
	"errors"

	"github.com/heandroro/agencia-viagem/backend/internal/domain/reservation"
)

// GetReservationSummaryUseCase caso de uso para buscar resumo da reserva
type GetReservationSummaryUseCase struct {
	reservationRepo reservation.Repository
}

// NewGetReservationSummaryUseCase cria um novo caso de uso
func NewGetReservationSummaryUseCase(repo reservation.Repository) *GetReservationSummaryUseCase {
	return &GetReservationSummaryUseCase{reservationRepo: repo}
}

// GetReservationSummaryInput input do caso de uso
type GetReservationSummaryInput struct {
	ReservationID string
	UserID        string
}

// SummaryPackageInfo informações do pacote no resumo
type SummaryPackageInfo struct {
	Name        string `json:"name"`
	Destination string `json:"destination"`
}

// SummaryDates datas no resumo
type SummaryDates struct {
	CheckIn  string `json:"check_in"`
	CheckOut string `json:"check_out"`
	Nights   int    `json:"nights"`
}

// SummaryTraveler viajante no resumo (com máscara)
type SummaryTraveler struct {
	FullName       string `json:"full_name"`
	DocumentMasked string `json:"document_masked"`
	Type           string `json:"type"`
}

// SummaryPricing preços no resumo
type SummaryPricing struct {
	PackagePrice float64 `json:"package_price"`
	Subtotal     float64 `json:"subtotal"`
	Taxes        float64 `json:"taxes"`
	Total        float64 `json:"total"`
	Currency     string  `json:"currency"`
}

// SummaryPolicies políticas no resumo
type SummaryPolicies struct {
	Cancellation string `json:"cancellation"`
	Modification string `json:"modification"`
}

// GetReservationSummaryOutput output do caso de uso
type GetReservationSummaryOutput struct {
	ReservationID string            `json:"reservation_id"`
	Status        string            `json:"status"`
	Package       SummaryPackageInfo `json:"package"`
	Dates         SummaryDates      `json:"dates"`
	Travelers     []SummaryTraveler `json:"travelers"`
	Pricing       SummaryPricing    `json:"pricing"`
	Policies      SummaryPolicies   `json:"policies"`
	ExpiresAt     string            `json:"expires_at"`
}

// Execute executa o caso de uso
func (uc *GetReservationSummaryUseCase) Execute(ctx context.Context, input GetReservationSummaryInput) (*GetReservationSummaryOutput, error) {
	rsv, err := uc.reservationRepo.GetByID(ctx, input.ReservationID)
	if err != nil {
		if err.Error() == "reservation_not_found" {
			return nil, errors.New("reservation_not_found")
		}
		return nil, err
	}

	if rsv.UserID != input.UserID {
		return nil, errors.New("unauthorized")
	}

	travelers := make([]SummaryTraveler, len(rsv.Travelers))
	for i, t := range rsv.Travelers {
		travelers[i] = SummaryTraveler{
			FullName:       t.FullName,
			DocumentMasked: reservation.MaskDocument(t.DocumentHash),
			Type:           string(t.Type),
		}
	}

	output := &GetReservationSummaryOutput{
		ReservationID: rsv.ID.String(),
		Status:        string(rsv.Status),
		Package: SummaryPackageInfo{
			Name:        "Pacote " + rsv.PackageID,
			Destination: "Destino",
		},
		Dates: SummaryDates{
			CheckIn:  rsv.Dates.StartDate.Format("2006-01-02"),
			CheckOut: rsv.Dates.EndDate.Format("2006-01-02"),
			Nights:   rsv.Dates.Nights,
		},
		Travelers: travelers,
		Pricing: SummaryPricing{
			PackagePrice: rsv.Pricing.PackagePrice,
			Subtotal:     rsv.Pricing.Subtotal,
			Taxes:        rsv.Pricing.Taxes,
			Total:        rsv.Pricing.Total,
			Currency:     rsv.Pricing.Currency,
		},
		Policies: SummaryPolicies{
			Cancellation: "Cancelamento gratuito até 48h antes do check-in. Após isso, taxa de 10%.",
			Modification: "Modificações permitidas até 72h antes, sujeito a disponibilidade",
		},
		ExpiresAt: rsv.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
