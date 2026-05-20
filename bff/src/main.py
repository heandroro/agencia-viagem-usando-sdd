"""
BFF - Backend for Frontend
Agência de Viagem
"""

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from src.api.routes import reservations
from src.core.config.settings import settings

app = FastAPI(
    title=settings.APP_NAME,
    description="Backend for Frontend da Agência de Viagem",
    version=settings.APP_VERSION,
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include routers
app.include_router(reservations.router)


@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy", "service": "bff", "version": settings.APP_VERSION}


@app.get("/")
async def root():
    """Root endpoint"""
    return {
        "message": "Agência Viagem BFF",
        "version": settings.APP_VERSION,
        "environment": settings.ENVIRONMENT,
    }
