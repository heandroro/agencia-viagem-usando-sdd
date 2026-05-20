"""Configurações do BFF"""
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    """Configurações da aplicação"""
    
    # App
    APP_NAME: str = "Agência Viagem BFF"
    APP_VERSION: str = "0.1.0"
    ENVIRONMENT: str = "development"
    PORT: int = 8000
    
    # Backend
    BACKEND_URL: str = "http://localhost:8080"
    BACKEND_TIMEOUT: float = 10.0
    
    # Valkey/Redis
    VALKEY_URL: str = "redis://localhost:6379/0"
    
    # JWT
    JWT_SECRET: str = "default-jwt-secret-change-in-production"
    JWT_ALGORITHM: str = "HS256"
    JWT_EXPIRATION_HOURS: int = 24
    
    # CORS
    CORS_ORIGINS: list[str] = ["*"]
    
    class Config:
        env_file = ".env"


settings = Settings()
