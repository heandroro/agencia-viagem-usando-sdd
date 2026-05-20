-- Migration: Criar tabela de reservas
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS reservations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) NOT NULL,
    package_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    nights INT NOT NULL,
    package_price DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    taxes DECIMAL(10,2) NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'BRL',
    expires_at TIMESTAMP NOT NULL,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index para busca por usuário
CREATE INDEX idx_reservations_user_status ON reservations(user_id, status);

-- Index para busca por pacote
CREATE INDEX idx_reservations_package ON reservations(package_id, start_date, end_date);

-- Index para expiração (TTL manual via worker)
CREATE INDEX idx_reservations_expires ON reservations(expires_at) WHERE status = 'pending';

-- Tabela de viajantes
CREATE TABLE IF NOT EXISTS travelers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reservation_id UUID NOT NULL REFERENCES reservations(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('primary', 'companion')),
    full_name VARCHAR(100) NOT NULL,
    document_type VARCHAR(20) NOT NULL CHECK (document_type IN ('cpf', 'passport')),
    document_number_encrypted TEXT NOT NULL,
    document_hash VARCHAR(64) NOT NULL,
    birth_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_travelers_reservation ON travelers(reservation_id);
