package config

import (
	"os"
	"strconv"
	"time"
)

// Config configuração da aplicação
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Valkey   ValkeyConfig
	Security SecurityConfig
	APM      APMConfig
}

// ServerConfig configuração do servidor
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// PostgresConfig configuração do PostgreSQL
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// ValkeyConfig configuração do Valkey
type ValkeyConfig struct {
	Addr     string
	Password string
	DB       int
}

// SecurityConfig configuração de segurança
type SecurityConfig struct {
	EncryptionKey string
	JWTSecret     string
}

// APMConfig configuração de observabilidade
type APMConfig struct {
	Enabled        bool
	ServiceName    string
	ServiceVersion string
	Environment    string
}

// Load carrega configuração das variáveis de ambiente
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDuration("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			Database: getEnv("POSTGRES_DATABASE", "agencia_viagem"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		Valkey: ValkeyConfig{
			Addr:     getEnv("VALKEY_ADDR", "localhost:6379"),
			Password: getEnv("VALKEY_PASSWORD", ""),
			DB:       getInt("VALKEY_DB", 0),
		},
		Security: SecurityConfig{
			EncryptionKey: getEnv("ENCRYPTION_KEY", "default-key-change-in-production-32b!"),
			JWTSecret:     getEnv("JWT_SECRET", "default-jwt-secret-change-in-production"),
		},
		APM: APMConfig{
			Enabled:        getBool("APM_ENABLED", true),
			ServiceName:    getEnv("APM_SERVICE_NAME", "reservation-service"),
			ServiceVersion: getEnv("APM_SERVICE_VERSION", "0.1.0"),
			Environment:    getEnv("APM_ENVIRONMENT", "development"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
