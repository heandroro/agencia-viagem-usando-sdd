package reservation

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

// MongoRepository implementação MongoDB do repositório
type MongoRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository cria um novo repositório MongoDB
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	repo := &MongoRepository{
		collection: db.Collection("reservations"),
	}
	
	// Criar índices
	ctx := context.Background()
	
	// Índice TTL para expiração automática
	repo.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	
	// Índice composto para busca por usuário e status
	repo.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "status", Value: 1},
		},
	})
	
	// Índice para busca por pacote e datas
	repo.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "package_id", Value: 1},
			{Key: "dates.start_date", Value: 1},
			{Key: "dates.end_date", Value: 1},
		},
	})
	
	return repo
}

// Create cria uma nova reserva
func (r *MongoRepository) Create(ctx context.Context, reservation *Reservation) error {
	reservation.ID = bson.NewObjectID()
	reservation.CreatedAt = time.Now()
	reservation.UpdatedAt = time.Now()
	reservation.Status = StatusPending
	
	_, err := r.collection.InsertOne(ctx, reservation)
	return err
}

// GetByID busca uma reserva pelo ID
func (r *MongoRepository) GetByID(ctx context.Context, id string) (*Reservation, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid_reservation_id")
	}
	
	var reservation Reservation
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&reservation)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("reservation_not_found")
	}
	if err != nil {
		return nil, err
	}
	
	return &reservation, nil
}

// GetByUserID busca reservas de um usuário
func (r *MongoRepository) GetByUserID(ctx context.Context, userID string, status Status) ([]Reservation, error) {
	filter := bson.M{"user_id": userID}
	if status != "" {
		filter["status"] = status
	}
	
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{
		{Key: "created_at", Value: -1},
	}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var reservations []Reservation
	if err := cursor.All(ctx, &reservations); err != nil {
		return nil, err
	}
	
	return reservations, nil
}

// Update atualiza uma reserva completa
func (r *MongoRepository) Update(ctx context.Context, reservation *Reservation) error {
	reservation.UpdatedAt = time.Now()
	
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": reservation.ID}, reservation)
	return err
}

// UpdateStatus atualiza apenas o status da reserva
func (r *MongoRepository) UpdateStatus(ctx context.Context, id string, status Status) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid_reservation_id")
	}
	
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	})
	return err
}

// UpdateTravelers atualiza apenas os viajantes da reserva
func (r *MongoRepository) UpdateTravelers(ctx context.Context, id string, travelers []Traveler) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid_reservation_id")
	}
	
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{
		"$set": bson.M{
			"travelers":  travelers,
			"updated_at": time.Now(),
		},
	})
	return err
}

// Delete remove uma reserva (soft delete recomendado para produção)
func (r *MongoRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid_reservation_id")
	}
	
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
