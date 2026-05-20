package reservation

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/heandroro/agencia-viagem/backend/internal/infra/database"
	"github.com/jackc/pgx/v5"
)

// Availability representa a disponibilidade de um pacote em uma data específica
type Availability struct {
	ID            uuid.UUID `json:"id"`
	PackageID     string    `json:"package_id"`
	Date          time.Time `json:"date"`
	TotalSlots    int       `json:"total_slots"`
	ReservedSlots int       `json:"reserved_slots"`
	Version       int       `json:"version"` // Para optimistic locking
	LastUpdated   time.Time `json:"last_updated"`
}

// AvailabilityRepository interface para operações de disponibilidade
type AvailabilityRepository interface {
	CheckAndReserve(ctx context.Context, packageID string, startDate, endDate time.Time, slots int) error
	ReleaseSlots(ctx context.Context, packageID string, startDate, endDate time.Time, slots int) error
	GetAvailability(ctx context.Context, packageID string, startDate, endDate time.Time) ([]Availability, error)
}

// PostgresAvailabilityRepository implementação PostgreSQL do repositório
type PostgresAvailabilityRepository struct {
	db *database.PostgresClient
}

// NewPostgresAvailabilityRepository cria um novo repositório PostgreSQL
func NewPostgresAvailabilityRepository(db *database.PostgresClient) *PostgresAvailabilityRepository {
	return &PostgresAvailabilityRepository{db: db}
}

// CheckAndReserve verifica disponibilidade e reserva slots usando optimistic locking
func (r *PostgresAvailabilityRepository) CheckAndReserve(ctx context.Context, packageID string, startDate, endDate time.Time, slots int) error {
	// Calcular número de dias
	days := int(endDate.Sub(startDate).Hours() / 24)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i := 0; i < days; i++ {
		date := startDate.Add(time.Duration(i) * 24 * time.Hour)

		// Verificar disponibilidade com row-level locking
		var avail Availability
		err := tx.QueryRow(ctx, `
			SELECT id, package_id, date, total_slots, reserved_slots, version, last_updated
			FROM availability
			WHERE package_id = $1 AND date = $2
			FOR UPDATE SKIP LOCKED
		`, packageID, date).Scan(
			&avail.ID, &avail.PackageID, &avail.Date, &avail.TotalSlots,
			&avail.ReservedSlots, &avail.Version, &avail.LastUpdated,
		)

		if err == pgx.ErrNoRows {
			// Criar novo registro de disponibilidade
			_, err = tx.Exec(ctx, `
				INSERT INTO availability (id, package_id, date, total_slots, reserved_slots, version, last_updated)
				VALUES ($1, $2, $3, 100, $4, 1, $5)
			`, uuid.New(), packageID, date, slots, time.Now())
			if err != nil {
				return err
			}
			continue
		}

		if err != nil {
			return err
		}

		// Verificar se há slots disponíveis
		if avail.TotalSlots-avail.ReservedSlots < slots {
			return errors.New("package_unavailable")
		}

		// Atualizar com optimistic locking
		result, err := tx.Exec(ctx, `
			UPDATE availability
			SET reserved_slots = reserved_slots + $1, version = version + 1, last_updated = $2
			WHERE id = $3 AND version = $4
		`, slots, time.Now(), avail.ID, avail.Version)

		if err != nil {
			return err
		}

		if result.RowsAffected() == 0 {
			return errors.New("concurrency_conflict")
		}
	}

	return tx.Commit(ctx)
}

// ReleaseSlots libera slots de disponibilidade
func (r *PostgresAvailabilityRepository) ReleaseSlots(ctx context.Context, packageID string, startDate, endDate time.Time, slots int) error {
	days := int(endDate.Sub(startDate).Hours() / 24)

	for i := 0; i < days; i++ {
		date := startDate.Add(time.Duration(i) * 24 * time.Hour)

		_, err := r.db.Pool().Exec(ctx, `
			INSERT INTO availability (id, package_id, date, total_slots, reserved_slots, version, last_updated)
			VALUES ($1, $2, $3, 100, -$4, 1, $5)
			ON CONFLICT (package_id, date) DO UPDATE
			SET reserved_slots = availability.reserved_slots - $4,
				version = availability.version + 1,
				last_updated = $5
		`, uuid.New(), packageID, date, slots, time.Now())

		if err != nil {
			return err
		}
	}

	return nil
}

// GetAvailability retorna disponibilidade para um período
func (r *PostgresAvailabilityRepository) GetAvailability(ctx context.Context, packageID string, startDate, endDate time.Time) ([]Availability, error) {
	rows, err := r.db.Pool().Query(ctx, `
		SELECT id, package_id, date, total_slots, reserved_slots, version, last_updated
		FROM availability
		WHERE package_id = $1 AND date >= $2 AND date < $3
		ORDER BY date
	`, packageID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Availability
	for rows.Next() {
		var avail Availability
		err = rows.Scan(
			&avail.ID, &avail.PackageID, &avail.Date, &avail.TotalSlots,
			&avail.ReservedSlots, &avail.Version, &avail.LastUpdated,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, avail)
	}

	return results, rows.Err()
}
