package reservation

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Status representa o status de uma reserva
type Status string

const (
	StatusPending    Status = "pending"
	StatusConfirmed  Status = "confirmed"
	StatusCancelled  Status = "cancelled"
	StatusExpired    Status = "expired"
)

// TravelerType indica se é viajante principal ou acompanhante
type TravelerType string

const (
	TravelerTypePrimary   TravelerType = "primary"
	TravelerTypeCompanion TravelerType = "companion"
)

// DocumentType tipo de documento do viajante
type DocumentType string

const (
	DocumentTypeCPF      DocumentType = "cpf"
	DocumentTypePassport DocumentType = "passport"
)

// Reservation representa uma reserva de pacote de viagem
type Reservation struct {
	ID               bson.ObjectID `bson:"_id,omitempty" json:"reservation_id"`
	UserID           string        `bson:"user_id" json:"user_id"`
	PackageID        string        `bson:"package_id" json:"package_id"`
	Status           Status        `bson:"status" json:"status"`
	CreatedAt        time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time     `bson:"updated_at" json:"updated_at"`
	ExpiresAt        time.Time     `bson:"expires_at" json:"expires_at"`
	Dates            DateRange     `bson:"dates" json:"dates"`
	Pricing          Pricing       `bson:"pricing" json:"pricing"`
	Travelers        []Traveler    `bson:"travelers" json:"travelers"`
	PoliciesAccepted bool          `bson:"policies_accepted" json:"policies_accepted"`
	Audit            AuditInfo     `bson:"audit" json:"audit"`
}

// DateRange representa as datas da viagem
type DateRange struct {
	StartDate time.Time `bson:"start_date" json:"start"`
	EndDate   time.Time `bson:"end_date" json:"end"`
	Nights    int       `bson:"nights" json:"nights"`
}

// Pricing representa a estrutura de preços
type Pricing struct {
	PackagePrice float64 `bson:"package_price" json:"package_price"`
	Subtotal     float64 `bson:"subtotal" json:"subtotal"`
	Taxes        float64 `bson:"taxes" json:"taxes"`
	Total        float64 `bson:"total" json:"total"`
	Currency     string  `bson:"currency" json:"currency"`
}

// Traveler representa um viajante da reserva
type Traveler struct {
	TravelerID        string       `bson:"traveler_id" json:"traveler_id"`
	Type              TravelerType `bson:"type" json:"type"`
	FullName          string       `bson:"full_name" json:"full_name"`
	DocumentType      DocumentType `bson:"document_type" json:"document_type"`
	DocumentEncrypted string       `bson:"document_encrypted" json:"-"` // Nunca exposto na API
	DocumentHash      string       `bson:"document_hash" json:"-"`    // Para correlacionar em logs
	BirthDate         time.Time    `bson:"birth_date" json:"birth_date"`
}

// AuditInfo informações de auditoria
type AuditInfo struct {
	IPAddress string `bson:"ip_address" json:"ip_address"`
	UserAgent string `bson:"user_agent" json:"user_agent"`
}

// CreateReservationInput dados para criar uma reserva
type CreateReservationInput struct {
	PackageID     string    `json:"package_id" validate:"required"`
	StartDate     time.Time `json:"start_date" validate:"required"`
	EndDate       time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
	TravelerCount int       `json:"traveler_count" validate:"required,min=1,max=10"`
	UserID        string    `json:"-"`
	IPAddress     string    `json:"-"`
	UserAgent     string    `json:"-"`
}

// UpdateTravelersInput dados para atualizar viajantes
type UpdateTravelersInput struct {
	ReservationID string         `json:"-" validate:"required"`
	Travelers     []TravelerInput `json:"travelers" validate:"required,min=1,max=10,dive"`
}

// TravelerInput dados de entrada de um viajante
type TravelerInput struct {
	Type           TravelerType `json:"type" validate:"required,oneof=primary companion"`
	FullName       string       `json:"full_name" validate:"required,min=3,max=100"`
	DocumentType   DocumentType `json:"document_type" validate:"required,oneof=cpf passport"`
	DocumentNumber string       `json:"document_number" validate:"required"`
	BirthDate      time.Time    `json:"birth_date" validate:"required"`
}

// ReservationSummary resumo da reserva para o usuário
type ReservationSummary struct {
	ReservationID string           `json:"reservation_id"`
	Status        Status           `json:"status"`
	Package       PackageInfo      `json:"package"`
	Dates         DateRange        `json:"dates"`
	Travelers     []TravelerSummary `json:"travelers"`
	Pricing       Pricing          `json:"pricing"`
	Policies      Policies         `json:"policies"`
	ExpiresAt     time.Time        `json:"expires_at"`
}

// PackageInfo informações do pacote no resumo
type PackageInfo struct {
	Name        string `json:"name"`
	Hotel       string `json:"hotel"`
	Flight      string `json:"flight"`
	Destination string `json:"destination"`
}

// TravelerSummary versão sanitizada do viajante para API
type TravelerSummary struct {
	FullName       string       `json:"full_name"`
	DocumentMasked string       `json:"document_masked"`
	Type           TravelerType `json:"type"`
}

// Policies políticas da reserva
type Policies struct {
	Cancellation string `json:"cancellation"`
	Modification string `json:"modification"`
}

// IsPending verifica se a reserva está pendente
func (r *Reservation) IsPending() bool {
	return r.Status == StatusPending
}

// IsExpired verifica se a reserva expirou
func (r *Reservation) IsExpired() bool {
	return r.Status == StatusExpired || time.Now().After(r.ExpiresAt)
}

// MaskDocument mascara um documento para exibição (ex: ***45678901)
func MaskDocument(doc string) string {
	if len(doc) <= 4 {
		return "****"
	}
	return "***" + doc[len(doc)-4:]
}
