package reservation

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Availability representa a disponibilidade de um pacote em uma data específica
type Availability struct {
	ID              bson.ObjectID `bson:"_id,omitempty"`
	PackageID       string        `bson:"package_id"`
	Date            time.Time     `bson:"date"`
	AvailableSlots  int           `bson:"available_slots"`
	ReservedSlots   int           `bson:"reserved_slots"`
	Version         int           `bson:"version"` // Para optimistic locking
	LastUpdated     time.Time     `bson:"last_updated"`
}

// AvailabilityRepository interface para operações de disponibilidade
type AvailabilityRepository interface {
	CheckAndReserve(ctx context.Context, packageID string, startDate, endDate time.Time, slots int) error
	ReleaseSlots(ctx context.Context, packageID string, startDate, endDate time.Time, slots int) error
	GetAvailability(ctx context.Context, packageID string, startDate, endDate time.Time) ([]Availability, error)
}

// MongoAvailabilityRepository implementação MongoDB do repositório
type MongoAvailabilityRepository struct {
	collection *mongo.Collection
}

// NewMongoAvailabilityRepository cria um novo repositório
func NewMongoAvailabilityRepository(db *mongo.Database) *MongoAvailabilityRepository {
	return &MongoAvailabilityRepository{
		collection: db.Collection("availability"),
	}
}

// CheckAndReserve verifica disponibilidade e reserva slots usando optimistic locking
func (r *MongoAvailabilityRepository) CheckAndReserve(ctx context.Context, packageID string, startDate, endDate time.Time, slots int) error {
	// Calcular número de dias
	days := int(endDate.Sub(startDate).Hours() / 24)
	
	for i := 0; i < days; i++ {
		date := startDate.Add(time.Duration(i) * 24 * time.Hour)
		
		// Tentativa de reserva com optimistic locking
		for retries := 0; retries < 3; retries++ {
			// Buscar registro atual
			var avail Availability
			err := r.collection.FindOne(ctx, bson.M{
				"package_id": packageID,
				"date":       date,
			}).Decode(&avail)
			
			if err == mongo.ErrNoDocuments {
				// Criar novo registro de disponibilidade
				_, err = r.collection.InsertOne(ctx, Availability{
					PackageID:      packageID,
					Date:           date,
					AvailableSlots: 100, // Default
					ReservedSlots:  slots,
					Version:        1,
					LastUpdated:    time.Now(),
				})
				if err != nil {
					return err
				}
				break
			}
			
			if err != nil {
				return err
			}
			
			// Verificar se há slots disponíveis
			if avail.AvailableSlots-avail.ReservedSlots < slots {
				return errors.New("package_unavailable")
			}
			
			// Tentar atualizar com optimistic locking
			result, err := r.collection.UpdateOne(ctx, bson.M{
				"_id":     avail.ID,
				"version": avail.Version,
			}, bson.M{
				"$inc": bson.M{
					"reserved_slots": slots,
					"version":        1,
				},
				"$set": bson.M{
					"last_updated": time.Now(),
				},
			})
			
			if err != nil {
				return err
			}
			
			if result.ModifiedCount > 0 {
				break // Sucesso
			}
			
			// Conflito de versão, tentar novamente
			if retries == 2 {
				return errors.New("concurrency_conflict")
			}
			
			time.Sleep(time.Millisecond * 10)
		}
	}
	
	return nil
}

// ReleaseSlots libera slots de disponibilidade
func (r *MongoAvailabilityRepository) ReleaseSlots(ctx context.Context, packageID string, startDate, endDate time.Time, slots int) error {
	days := int(endDate.Sub(startDate).Hours() / 24)
	
	for i := 0; i < days; i++ {
		date := startDate.Add(time.Duration(i) * 24 * time.Hour)
		
		_, err := r.collection.UpdateOne(ctx, bson.M{
			"package_id": packageID,
			"date":       date,
		}, bson.M{
			"$inc": bson.M{
				"reserved_slots": -slots,
				"version":        1,
			},
			"$set": bson.M{
				"last_updated": time.Now(),
			},
		}, options.Update().SetUpsert(true))
		
		if err != nil {
			return err
		}
	}
	
	return nil
}

// GetAvailability retorna disponibilidade para um período
func (r *MongoAvailabilityRepository) GetAvailability(ctx context.Context, packageID string, startDate, endDate time.Time) ([]Availability, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"package_id": packageID,
		"date": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []Availability
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	
	return results, nil
}
