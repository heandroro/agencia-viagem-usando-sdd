package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// Config configuração do MongoDB
type Config struct {
	URI      string
	Database string
	Timeout  time.Duration
}

// Client wrapper para cliente MongoDB
type Client struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewClient cria uma nova conexão MongoDB
func NewClient(cfg Config) (*Client, error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.URI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	// Ping para verificar conexão
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	return &Client{
		client:   client,
		database: client.Database(cfg.Database),
	}, nil
}

// Database retorna o banco de dados
func (c *Client) Database() *mongo.Database {
	return c.database
}

// Close fecha a conexão
func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

// Health verifica se a conexão está saudável
func (c *Client) Health(ctx context.Context) error {
	return c.client.Ping(ctx, readpref.Primary())
}
