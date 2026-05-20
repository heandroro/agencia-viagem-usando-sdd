package reservation

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/heandroro/agencia-viagem/backend/internal/infra/database"
	"github.com/jackc/pgx/v5"
)

// Repository interface para operações de reserva
type Repository interface {
	Create(ctx context.Context, reservation *Reservation) error
	GetByID(ctx context.Context, id string) (*Reservation, error)
	GetByUserID(ctx context.Context, userID string, status Status) ([]Reservation, error)
	Update(ctx context.Context, reservation *Reservation) error
	UpdateStatus(ctx context.Context, id string, status Status) error
	UpdateTravelers(ctx context.Context, id string, travelers []Traveler) error
	Delete(ctx context.Context, id string) error
}

// PostgresRepository implementação PostgreSQL do repositório
type PostgresRepository struct {
	db *database.PostgresClient
}

// NewPostgresRepository cria um novo repositório PostgreSQL
func NewPostgresRepository(db *database.PostgresClient) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create cria uma nova reserva
func (r *PostgresRepository) Create(ctx context.Context, reservation *Reservation) error {
	reservation.ID = uuid.New()
	reservation.CreatedAt = time.Now()
	reservation.UpdatedAt = time.Now()
	reservation.Status = StatusPending

	sql := `
		INSERT INTO reservations (
			id, user_id, package_id, status, start_date, end_date, nights,
			package_price, subtotal, taxes, total, currency, expires_at,
			ip_address, user_agent, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err := r.db.Pool().Exec(ctx, sql,
		reservation.ID, reservation.UserID, reservation.PackageID, reservation.Status,
		reservation.Dates.StartDate, reservation.Dates.EndDate, reservation.Dates.Nights,
		reservation.Pricing.PackagePrice, reservation.Pricing.Subtotal,
		reservation.Pricing.Taxes, reservation.Pricing.Total, reservation.Pricing.Currency,
		reservation.ExpiresAt,
		reservation.Audit.IPAddress,
		reservation.Audit.UserAgent,
		reservation.CreatedAt, reservation.UpdatedAt,
	)

	return err
}

// GetByID busca uma reserva pelo ID
func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Reservation, error) {
	reservationID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid_reservation_id")
	}

	sql := `
		SELECT id, user_id, package_id, status, start_date, end_date, nights,
			package_price, subtotal, taxes, total, currency, expires_at,
			ip_address, user_agent, created_at, updated_at
		FROM reservations WHERE id = $1
	`

	var rsv Reservation
	var ipAddress, userAgent *string

	err = r.db.Pool().QueryRow(ctx, sql, reservationID).Scan(
		&rsv.ID, &rsv.UserID, &rsv.PackageID, &rsv.Status,
		&rsv.Dates.StartDate, &rsv.Dates.EndDate, &rsv.Dates.Nights,
		&rsv.Pricing.PackagePrice, &rsv.Pricing.Subtotal,
		&rsv.Pricing.Taxes, &rsv.Pricing.Total, &rsv.Pricing.Currency,
		&rsv.ExpiresAt,
		&ipAddress, &userAgent,
		&rsv.CreatedAt, &rsv.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New("reservation_not_found")
	}
	if err != nil {
		return nil, err
	}

	if ipAddress != nil {
		rsv.Audit.IPAddress = *ipAddress
	}
	if userAgent != nil {
		rsv.Audit.UserAgent = *userAgent
	}

	// Carregar viajantes
	travelers, err := r.getTravelers(ctx, rsv.ID)
	if err != nil {
		return nil, err
	}
	rsv.Travelers = travelers

	return &rsv, nil
}

// GetByUserID busca reservas de um usuário
func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string, status Status) ([]Reservation, error) {
	sql := `
		SELECT id, user_id, package_id, status, start_date, end_date, nights,
			package_price, subtotal, taxes, total, currency, expires_at,
			ip_address, user_agent, created_at, updated_at
		FROM reservations WHERE user_id = $1
	`
	args := []interface{}{userID}

	if status != "" {
		sql += " AND status = $2"
		args = append(args, status)
	}

	sql += " ORDER BY created_at DESC"

	rows, err := r.db.Pool().Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []Reservation
	for rows.Next() {
		var rsv Reservation
		var ipAddress, userAgent *string

		err = rows.Scan(
			&rsv.ID, &rsv.UserID, &rsv.PackageID, &rsv.Status,
			&rsv.Dates.StartDate, &rsv.Dates.EndDate, &rsv.Dates.Nights,
			&rsv.Pricing.PackagePrice, &rsv.Pricing.Subtotal,
			&rsv.Pricing.Taxes, &rsv.Pricing.Total, &rsv.Pricing.Currency,
			&rsv.ExpiresAt,
			&ipAddress, &userAgent,
			&rsv.CreatedAt, &rsv.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if ipAddress != nil {
			rsv.Audit.IPAddress = *ipAddress
		}
		if userAgent != nil {
			rsv.Audit.UserAgent = *userAgent
		}
		reservations = append(reservations, rsv)
	}

	return reservations, rows.Err()
}

// Update atualiza uma reserva completa
func (r *PostgresRepository) Update(ctx context.Context, reservation *Reservation) error {
	reservation.UpdatedAt = time.Now()

	sql := `
		UPDATE reservations SET
			status = $1, start_date = $2, end_date = $3, nights = $4,
			package_price = $5, subtotal = $6, taxes = $7, total = $8,
			updated_at = $9
		WHERE id = $10
	`

	_, err := r.db.Pool().Exec(ctx, sql,
		reservation.Status, reservation.Dates.StartDate, reservation.Dates.EndDate, reservation.Dates.Nights,
		reservation.Pricing.PackagePrice, reservation.Pricing.Subtotal,
		reservation.Pricing.Taxes, reservation.Pricing.Total,
		reservation.UpdatedAt, reservation.ID,
	)

	return err
}

// UpdateStatus atualiza apenas o status da reserva
func (r *PostgresRepository) UpdateStatus(ctx context.Context, id string, status Status) error {
	reservationID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid_reservation_id")
	}

	sql := `UPDATE reservations SET status = $1, updated_at = $2 WHERE id = $3`
	_, err = r.db.Pool().Exec(ctx, sql, status, time.Now(), reservationID)
	return err
}

// UpdateTravelers atualiza os viajantes da reserva
func (r *PostgresRepository) UpdateTravelers(ctx context.Context, id string, travelers []Traveler) error {
	reservationID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid_reservation_id")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Deletar viajantes existentes
	_, err = tx.Exec(ctx, "DELETE FROM travelers WHERE reservation_id = $1", reservationID)
	if err != nil {
		return err
	}

	// Inserir novos viajantes
	for _, t := range travelers {
		_, err = tx.Exec(ctx, `
			INSERT INTO travelers (id, reservation_id, "type", full_name, document_type, document_number_encrypted, document_hash, birth_date)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, uuid.New(), reservationID, string(t.Type), t.FullName, string(t.DocumentType), t.DocumentEncrypted, t.DocumentHash, t.BirthDate)
		if err != nil {
			return err
		}
	}

	// Atualizar timestamp da reserva
	_, err = tx.Exec(ctx, "UPDATE reservations SET updated_at = $1 WHERE id = $2", time.Now(), reservationID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Delete remove uma reserva
func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	reservationID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid_reservation_id")
	}

	_, err = r.db.Pool().Exec(ctx, "DELETE FROM reservations WHERE id = $1", reservationID)
	return err
}

// getTravelers carrega os viajantes de uma reserva
func (r *PostgresRepository) getTravelers(ctx context.Context, reservationID uuid.UUID) ([]Traveler, error) {
	rows, err := r.db.Pool().Query(ctx, `
		SELECT type, full_name, document_type, document_number_encrypted, document_hash, birth_date
		FROM travelers WHERE reservation_id = $1
	`, reservationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var travelers []Traveler
	for rows.Next() {
		var t Traveler
		err = rows.Scan(&t.Type, &t.FullName, &t.DocumentType, &t.DocumentEncrypted, &t.DocumentHash, &t.BirthDate)
		if err != nil {
			return nil, err
		}
		travelers = append(travelers, t)
	}

	return travelers, rows.Err()
}
