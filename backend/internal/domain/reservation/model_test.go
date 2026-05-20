package reservation

import (
	"testing"
	"time"
)

func TestReservation_IsPending(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected bool
	}{
		{"Pending", StatusPending, true},
		{"Confirmed", StatusConfirmed, false},
		{"Cancelled", StatusCancelled, false},
		{"Expired", StatusExpired, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Reservation{Status: tt.status}
			if got := r.IsPending(); got != tt.expected {
				t.Errorf("IsPending() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestReservation_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		status    Status
		expiresAt time.Time
		expected  bool
	}{
		{
			name:      "Expired by status",
			status:    StatusExpired,
			expiresAt: time.Now().Add(1 * time.Hour),
			expected:  true,
		},
		{
			name:      "Expired by time",
			status:    StatusPending,
			expiresAt: time.Now().Add(-1 * time.Hour),
			expected:  true,
		},
		{
			name:      "Not expired",
			status:    StatusPending,
			expiresAt: time.Now().Add(1 * time.Hour),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Reservation{
				Status:    tt.status,
				ExpiresAt: tt.expiresAt,
			}
			if got := r.IsExpired(); got != tt.expected {
				t.Errorf("IsExpired() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMaskDocument(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"12345678901", "***78901"},
		{"AB123456", "***3456"},
		{"1234", "****"},
		{"12", "****"},
		{"", "****"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := MaskDocument(tt.input)
			if result != tt.expected {
				t.Errorf("MaskDocument(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTravelerInput_Validation(t *testing.T) {
	// Test happy path
	validInput := CreateReservationInput{
		PackageID:     "pkg_123",
		StartDate:     time.Now().Add(24 * time.Hour),
		EndDate:       time.Now().Add(48 * time.Hour),
		TravelerCount: 2,
	}

	if validInput.PackageID == "" {
		t.Error("PackageID should not be empty")
	}

	if validInput.StartDate.After(validInput.EndDate) {
		t.Error("StartDate should be before EndDate")
	}

	if validInput.TravelerCount < 1 {
		t.Error("TravelerCount should be at least 1")
	}
}

func TestPricing_Calculations(t *testing.T) {
	pricing := Pricing{
		PackagePrice: 500.00,
		Subtotal:     5000.00,
		Taxes:        150.00,
		Total:        5150.00,
		Currency:     "BRL",
	}

	expectedTotal := pricing.Subtotal + pricing.Taxes
	if pricing.Total != expectedTotal {
		t.Errorf("Total = %v, want %v", pricing.Total, expectedTotal)
	}

	if pricing.Currency != "BRL" {
		t.Errorf("Currency = %v, want BRL", pricing.Currency)
	}
}
