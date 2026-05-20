package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/heandroro/agencia-viagem/backend/internal/domain/reservation"
	"github.com/heandroro/agencia-viagem/backend/pkg/crypto"
)

// CreateReservationUseCase caso de uso para criar reservas
type CreateReservationUseCase struct {
	reservationRepo  reservation.Repository
	availabilityRepo reservation.AvailabilityRepository
	crypto           crypto.Service
}

// NewCreateReservationUseCase cria um novo caso de uso
func NewCreateReservationUseCase(
	reservationRepo reservation.Repository,
	availabilityRepo reservation.AvailabilityRepository,
	crypto crypto.Service,
) *CreateReservationUseCase {
	return &CreateReservationUseCase{
		reservationRepo:  reservationRepo,
		availabilityRepo: availabilityRepo,
		crypto:           crypto,
	}
}

// CreateReservationOutput output do caso de uso
type CreateReservationOutput struct {
	ReservationID string                `json:"reservation_id"`
	Status        reservation.Status    `json:"status"`
	Package       PackageInfo           `json:"package"`
	Dates         reservation.DateRange `json:"dates"`
	Pricing       reservation.Pricing   `json:"pricing"`
	ExpiresAt     time.Time             `json:"expires_at"`
}

// PackageInfo informações do pacote na resposta
type PackageInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Destination string `json:"destination"`
}

// Execute executa o caso de uso
func (uc *CreateReservationUseCase) Execute(ctx context.Context, input reservation.CreateReservationInput) (*CreateReservationOutput, error) {
	// Validar datas
	if input.StartDate.After(input.EndDate) || input.StartDate.Before(time.Now()) {
		return nil, errors.New("invalid_dates")
	}

	// Calcular número de noites
	nights := int(input.EndDate.Sub(input.StartDate).Hours() / 24)
	if nights <= 0 {
		return nil, errors.New("invalid_dates")
	}

	// TODO: Buscar informações do pacote do catálogo (cache ou MongoDB)
	// Por enquanto, usamos valores mockados para o MVP
	packageInfo := struct {
		Name        string
		Destination string
		PricePerDay float64
	}{
		Name:        fmt.Sprintf("Pacote %s", input.PackageID),
		Destination: "Destino", // Seria obtido do catálogo
		PricePerDay: 500.00,    // Seria obtido do catálogo
	}

	// Calcular preços
	packagePrice := packageInfo.PricePerDay
	subtotal := packagePrice * float64(nights) * float64(input.TravelerCount)
	taxes := subtotal * 0.03 // 3% de taxas
	total := subtotal + taxes

	// Verificar disponibilidade e reservar slots
	err := uc.availabilityRepo.CheckAndReserve(ctx, input.PackageID, input.StartDate, input.EndDate, input.TravelerCount)
	if err != nil {
		return nil, errors.New("package_unavailable")
	}

	// Criar reserva
	reservation := &reservation.Reservation{
		UserID:    input.UserID,
		PackageID: input.PackageID,
		Status:    reservation.StatusPending,
		Dates: reservation.DateRange{
			StartDate: input.StartDate,
			EndDate:   input.EndDate,
			Nights:    nights,
		},
		Pricing: reservation.Pricing{
			PackagePrice: packagePrice,
			Subtotal:     subtotal,
			Taxes:        taxes,
			Total:        total,
			Currency:     "BRL",
		},
		ExpiresAt: time.Now().Add(30 * time.Minute),
		Audit: reservation.AuditInfo{
			IPAddress: input.IPAddress,
			UserAgent: input.UserAgent,
		},
	}

	// Salvar reserva
	if err := uc.reservationRepo.Create(ctx, reservation); err != nil {
		// Em caso de erro, liberar os slots reservados
		_ = uc.availabilityRepo.ReleaseSlots(ctx, input.PackageID, input.StartDate, input.EndDate, input.TravelerCount)
		return nil, err
	}

	// Montar output
	output := &CreateReservationOutput{
		ReservationID: reservation.ID.Hex(),
		Status:        reservation.Status,
		Package: PackageInfo{
			ID:          input.PackageID,
			Name:        packageInfo.Name,
			Destination: packageInfo.Destination,
		},
		Dates:     reservation.Dates,
		Pricing:   reservation.Pricing,
		ExpiresAt: reservation.ExpiresAt,
	}

	return output, nil
}
