"""
BFF - Backend for Frontend
Agência de Viagem
"""

from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(
    title="Agência Viagem BFF",
    description="Backend for Frontend da Agência de Viagem",
    version="0.1.0",
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Ajustar em produção
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy", "service": "bff"}


@app.get("/")
async def root():
    """Root endpoint"""
    return {
        "message": "Agência Viagem BFF",
        "version": "0.1.0",
    }
