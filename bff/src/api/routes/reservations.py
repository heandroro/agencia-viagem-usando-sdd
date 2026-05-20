"""Rotas de reservas"""
from typing import Annotated
from fastapi import APIRouter, HTTPException, Header, Depends
from pydantic import BaseModel, Field

from src.infra.http_client.backend_client import get_backend_client, BackendClient


router = APIRouter(prefix="/api/v1/reservations", tags=["reservations"])


# Request/Response Models
class CreateReservationRequest(BaseModel):
    """Requisição de criação de reserva"""
    package_id: str = Field(..., description="ID do pacote")
    start_date: str = Field(..., description="Data de início (ISO 8601)")
    end_date: str = Field(..., description="Data de término (ISO 8601)")
    traveler_count: int = Field(..., ge=1, le=10, description="Número de viajantes")


class TravelerInput(BaseModel):
    """Dados de um viajante"""
    type: str = Field(..., pattern="^(primary|companion)$")
    full_name: str = Field(..., min_length=3, max_length=100)
    document_type: str = Field(..., pattern="^(cpf|passport)$")
    document_number: str = Field(..., min_length=5)
    birth_date: str = Field(..., description="Data de nascimento (ISO 8601)")


class UpdateTravelersRequest(BaseModel):
    """Requisição de atualização de viajantes"""
    travelers: list[TravelerInput] = Field(..., min_length=1, max_length=10)


class ReservationResponse(BaseModel):
    """Resposta de reserva criada"""
    reservation_id: str
    status: str
    package: dict
    dates: dict
    pricing: dict
    expires_at: str


# Dependencies
async def get_user_id(x_user_id: Annotated[str, Header()] = "anonymous") -> str:
    """Extrai user ID do header (TODO: implementar JWT)"""
    return x_user_id


@router.post("", response_model=ReservationResponse)
async def create_reservation(
    request: CreateReservationRequest,
    user_id: str = Depends(get_user_id),
    backend: BackendClient = Depends(get_backend_client)
) -> ReservationResponse:
    """
    Cria uma nova reserva de pacote de viagem.
    
    - Valida disponibilidade do pacote
    - Calcula preço (snapshot)
    - Cria reserva com TTL de 30 minutos
    """
    try:
        result = await backend.create_reservation(
            package_id=request.package_id,
            start_date=request.start_date,
            end_date=request.end_date,
            traveler_count=request.traveler_count,
            user_id=user_id
        )
        return ReservationResponse(**result)
    except Exception as e:
        if "package_unavailable" in str(e):
            raise HTTPException(
                status_code=409,
                detail={
                    "error": "package_unavailable",
                    "message": "Pacote não disponível para as datas selecionadas"
                }
            )
        raise HTTPException(
            status_code=500,
            detail={"error": "internal_error", "message": "Erro ao criar reserva"}
        )


@router.put("/{reservation_id}/travelers")
async def update_travelers(
    reservation_id: str,
    request: UpdateTravelersRequest,
    user_id: str = Depends(get_user_id),
    backend: BackendClient = Depends(get_backend_client)
) -> dict:
    """
    Atualiza dados dos viajantes de uma reserva.
    
    - Valida formato de CPF/passaporte
    - Criptografa documentos
    - Requer reserva em status "pending"
    """
    try:
        travelers_data = [t.model_dump() for t in request.travelers]
        result = await backend.update_travelers(
            reservation_id=reservation_id,
            travelers=travelers_data,
            user_id=user_id
        )
        return result
    except Exception as e:
        if "validation" in str(e).lower():
            raise HTTPException(
                status_code=400,
                detail={"error": "validation_error", "message": "Dados inválidos"}
            )
        raise HTTPException(
            status_code=500,
            detail={"error": "internal_error", "message": "Erro ao atualizar viajantes"}
        )


@router.get("/{reservation_id}/summary")
async def get_reservation_summary(
    reservation_id: str,
    user_id: str = Depends(get_user_id),
    backend: BackendClient = Depends(get_backend_client)
) -> dict:
    """
    Retorna resumo completo da reserva.
    
    - Agrega dados do pacote
    - Inclui breakdown de preços
    - Mostra políticas de cancelamento
    """
    try:
        result = await backend.get_reservation_summary(
            reservation_id=reservation_id,
            user_id=user_id
        )
        return result
    except Exception as e:
        if "not_found" in str(e).lower():
            raise HTTPException(
                status_code=404,
                detail={"error": "not_found", "message": "Reserva não encontrada"}
            )
        raise HTTPException(
            status_code=500,
            detail={"error": "internal_error", "message": "Erro ao obter resumo"}
        )
