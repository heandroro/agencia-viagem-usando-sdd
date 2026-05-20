-- Migration: Criar tabela de disponibilidade
CREATE TABLE IF NOT EXISTS availability (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    package_id VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    total_slots INT NOT NULL DEFAULT 10,
    reserved_slots INT NOT NULL DEFAULT 0,
    version INT NOT NULL DEFAULT 1,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(package_id, date)
);

-- Index para busca de disponibilidade
CREATE INDEX idx_availability_package_date ON availability(package_id, date);

-- Função para optimistic locking (row-level locking)
CREATE OR REPLACE FUNCTION check_availability(
    p_package_id VARCHAR,
    p_date DATE,
    p_slots INT
) RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM availability
        WHERE package_id = p_package_id
        AND date = p_date
        AND (total_slots - reserved_slots) >= p_slots
        FOR UPDATE SKIP LOCKED
    );
END;
$$ LANGUAGE plpgsql;
