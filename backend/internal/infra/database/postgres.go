package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresConfig configuração do PostgreSQL
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
	Timeout  time.Duration
}

// PostgresClient wrapper para pool PostgreSQL
type PostgresClient struct {
	pool *pgxpool.Pool
}

// NewPostgresClient cria uma nova conexão PostgreSQL
func NewPostgresClient(cfg PostgresConfig) (*PostgresClient, error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode, int(cfg.Timeout.Seconds()),
	)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Testar conexão
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresClient{pool: pool}, nil
}

// Pool retorna o pool de conexões
func (c *PostgresClient) Pool() *pgxpool.Pool {
	return c.pool
}

// Close fecha o pool
func (c *PostgresClient) Close() {
	c.pool.Close()
}

// Health verifica se a conexão está saudável
func (c *PostgresClient) Health(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

// Begin inicia uma transação
func (c *PostgresClient) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.pool.Begin(ctx)
}

// Exec executa uma query sem retorno
func (c *PostgresClient) Exec(ctx context.Context, sql string, arguments ...interface{}) error {
	_, err := c.pool.Exec(ctx, sql, arguments...)
	return err
}

// QueryRow executa uma query que retorna uma única linha
func (c *PostgresClient) QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row {
	return c.pool.QueryRow(ctx, sql, arguments...)
}

// Query executa uma query que retorna múltiplas linhas
func (c *PostgresClient) Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error) {
	return c.pool.Query(ctx, sql, arguments...)
}
