FROM python:3.12-slim

WORKDIR /app

# Instalar dependências do sistema
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Copiar arquivos de dependências
COPY pyproject.toml .

# Instalar dependências
RUN pip install --no-cache-dir -e ".[dev]"

# Copiar código
COPY src/ ./src/

# Expor porta
EXPOSE 8000

# Executar
CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "8000"]
