package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/heandroro/agencia-viagem/backend/internal/domain/reservation"
	"github.com/heandroro/agencia-viagem/backend/pkg/crypto"
)

// Mock implementations for testing

type mockReservationRepo struct {
	createCalled bool
	lastCreated  *reservation.Reservation
	createError  error
}

func (m *mockReservationRepo) Create(ctx context.Context, r *reservation.Reservation) error {
	m.createCalled = true
	m.lastCreated = r
	return m.createError
}

func (m *mockReservationRepo) GetByID(ctx context.Context, id string) (*reservation.Reservation, error) {
	return nil, nil
}

func (m *mockReservationRepo) GetByUserID(ctx context.Context, userID string, status reservation.Status) ([]reservation.Reservation, error) {
	return nil, nil
}

func (m *mockReservationRepo) Update(ctx context.Context, r *reservation.Reservation) error {
	return nil
}

func (m *mockReservationRepo) UpdateStatus(ctx context.Context, id string, status reservation.Status) error {
	return nil
}

func (m *mockReservationRepo) UpdateTravelers(ctx context.Context, id string, travelers []reservation.Traveler) error {
	return nil
}

func (m *mockReservationRepo) Delete(ctx context.Context, id string) error {
	return nil
}

type mockAvailabilityRepo struct {
	checkCalled    bool
	releaseCalled  bool
	checkAndReserveError error
}

func (m *mockAvailabilityRepo) CheckAndReserve(ctx context.Context, packageID string, startDate, endDate time.Time, slots int) error {
	m.checkCalled = true
	return m.checkAndReserveError
}

func (m *mockAvailabilityRepo) ReleaseSlots(ctx context.Context, packageID string, startDate, endDate time.Time, slots int) error {
	m.releaseCalled = true
	return nil
}

func (m *mockAvailabilityRepo) GetAvailability(ctx context.Context, packageID string, startDate, endDate time.Time) ([]reservation.Availability, error) {
	return nil, nil
}

func TestCreateReservationUseCase_Execute_Success(t *testing.T) {
	// Setup
	reservationRepo := &mockReservationRepo{}
	availabilityRepo := &mockAvailabilityRepo{}
	cryptoService, _ := crypto.NewAES256Service("test-key-that-is-32-bytes-long!")

	uc := NewCreateReservationUseCase(reservationRepo, availabilityRepo, cryptoService)

	input := reservation.CreateReservationInput{
		PackageID:     "pkg_123",
		StartDate:     time.Now().Add(24 * time.Hour),
		EndDate:       time.Now().Add(48 * time.Hour),
		TravelerCount: 2,
		UserID:        "user_456",
	}

	// Execute
	output, err := uc.Execute(context.Background(), input)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if output == nil {
		t.Fatal("Expected output, got nil")
	}

	if output.Status != reservation.StatusPending {
		t.Errorf("Expected status pending, got: %v", output.Status)
	}

	if !availabilityRepo.checkCalled {
		t.Error("Expected CheckAndReserve to be called")
	}

	if !reservationRepo.createCalled {
		t.Error("Expected Create to be called")
	}

	// Verificar cálculo de preços
	expectedNights := 1
	if output.Dates.Nights != expectedNights {
		t.Errorf("Expected %d nights, got: %d", expectedNights, output.Dates.Nights)
	}

	if output.Pricing.Total <= 0 {
		t.Error("Expected positive total price")
	}
}

func TestCreateReservationUseCase_Execute_PackageUnavailable(t *testing.T) {
	// Setup
	reservationRepo := &mockReservationRepo{}
	availabilityRepo := &mockAvailabilityRepo{
		checkAndReserveError: errors.New("package_unavailable"),
	}
	cryptoService, _ := crypto.NewAES256Service("test-key-that-is-32-bytes-long!")

	uc := NewCreateReservationUseCase(reservationRepo, availabilityRepo, cryptoService)

	input := reservation.CreateReservationInput{
		PackageID:     "pkg_123",
		StartDate:     time.Now().Add(24 * time.Hour),
		EndDate:       time.Now().Add(48 * time.Hour),
		TravelerCount: 2,
		UserID:        "user_456",
	}

	// Execute
	output, err := uc.Execute(context.Background(), input)

	// Assert
	if err == nil {
		t.Error("Expected error for unavailable package")
	}

	if err.Error() != "package_unavailable" {
		t.Errorf("Expected 'package_unavailable' error, got: %v", err.Error())
	}

	if output != nil {
		t.Error("Expected nil output on error")
	}

	if reservationRepo.createCalled {
		t.Error("Expected Create NOT to be called when package unavailable")
	}
}

func TestCreateReservationUseCase_Execute_InvalidDates(t *testing.T) {
	// Setup
	reservationRepo := &mockReservationRepo{}
	availabilityRepo := &mockAvailabilityRepo{}
	cryptoService, _ := crypto.NewAES256Service("test-key-that-is-32-bytes-long!")

	uc := NewCreateReservationUseCase(reservationRepo, availabilityRepo, cryptoService)

	tests := []struct {
		name      string
		startDate time.Time
		endDate   time.Time
	}{
		{
			name:      "Start in past",
			startDate: time.Now().Add(-24 * time.Hour),
			endDate:   time.Now().Add(24 * time.Hour),
		},
		{
			name:      "End before start",
			startDate: time.Now().Add(48 * time.Hour),
			endDate:   time.Now().Add(24 * time.Hour),
		},
		{
			name:      "Same day",
			startDate: time.Now().Add(24 * time.Hour),
			endDate:   time.Now().Add(24 * time.Hour),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := reservation.CreateReservationInput{
				PackageID:     "pkg_123",
				StartDate:     tt.startDate,
				EndDate:       tt.endDate,
				TravelerCount: 2,
				UserID:        "user_456",
			}

			output, err := uc.Execute(context.Background(), input)

			if err == nil || err.Error() != "invalid_dates" {
				t.Errorf("Expected 'invalid_dates' error, got: %v", err)
			}

			if output != nil {
				t.Error("Expected nil output on invalid dates")
			}
		})
	}
}

func TestCreateReservationUseCase_Execute_RollsBackOnCreateError(t *testing.T) {
	// Setup
	reservationRepo := &mockReservationRepo{
		createError: errors.New("database_error"),
	}
	availabilityRepo := &mockAvailabilityRepo{}
	cryptoService, _ := crypto.NewAES256Service("test-key-that-is-32-bytes-long!")

	uc := NewCreateReservationUseCase(reservationRepo, availabilityRepo, cryptoService)

	input := reservation.CreateReservationInput{
		PackageID:     "pkg_123",
		StartDate:     time.Now().Add(24 * time.Hour),
		EndDate:       time.Now().Add(48 * time.Hour),
		TravelerCount: 2,
		UserID:        "user_456",
	}

	// Execute
	_, err := uc.Execute(context.Background(), input)

	// Assert
	if err == nil {
		t.Error("Expected error when create fails")
	}

	if !availabilityRepo.releaseCalled {
		t.Error("Expected ReleaseSlots to be called to rollback reservation")
	}
}
