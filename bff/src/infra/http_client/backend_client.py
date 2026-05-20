"""Cliente HTTP para backend Go"""
import httpx
from typing import Optional, Any
from src.core.config.settings import settings


class BackendClient:
    """Cliente para comunicação com backend Go"""
    
    def __init__(self, base_url: Optional[str] = None, timeout: Optional[float] = None):
        self.base_url = base_url or settings.BACKEND_URL
        self.timeout = timeout or settings.BACKEND_TIMEOUT
        self.client = httpx.AsyncClient(timeout=self.timeout)
    
    async def create_reservation(
        self, 
        package_id: str, 
        start_date: str, 
        end_date: str, 
        traveler_count: int,
        user_id: str
    ) -> dict[str, Any]:
        """Cria uma nova reserva no backend"""
        payload = {
            "package_id": package_id,
            "start_date": start_date,
            "end_date": end_date,
            "traveler_count": traveler_count,
            "user_id": user_id
        }
        
        response = await self.client.post(
            f"{self.base_url}/api/v1/reservations",
            json=payload,
            headers={
                "X-User-ID": user_id,
                "Content-Type": "application/json"
            }
        )
        response.raise_for_status()
        return response.json()
    
    async def update_travelers(
        self,
        reservation_id: str,
        travelers: list[dict],
        user_id: str
    ) -> dict[str, Any]:
        """Atualiza viajantes de uma reserva"""
        payload = {"travelers": travelers}
        
        response = await self.client.put(
            f"{self.base_url}/api/v1/reservations/{reservation_id}/travelers",
            json=payload,
            headers={"X-User-ID": user_id}
        )
        response.raise_for_status()
        return response.json()
    
    async def get_reservation_summary(
        self,
        reservation_id: str,
        user_id: str
    ) -> dict[str, Any]:
        """Obtém resumo da reserva"""
        response = await self.client.get(
            f"{self.base_url}/api/v1/reservations/{reservation_id}/summary",
            headers={"X-User-ID": user_id}
        )
        response.raise_for_status()
        return response.json()
    
    async def health_check(self) -> dict[str, Any]:
        """Verifica saúde do backend"""
        response = await self.client.get(f"{self.base_url}/health")
        response.raise_for_status()
        return response.json()
    
    async def close(self) -> None:
        """Fecha conexões"""
        await self.client.aclose()


# Singleton instance
_backend_client: Optional[BackendClient] = None


def get_backend_client() -> BackendClient:
    """Obtém instância do cliente backend (singleton)"""
    global _backend_client
    if _backend_client is None:
        _backend_client = BackendClient()
    return _backend_client
