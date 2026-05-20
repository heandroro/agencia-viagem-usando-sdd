package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heandroro/agencia-viagem/backend/internal/domain/reservation"
	"github.com/heandroro/agencia-viagem/backend/internal/usecase"
	"github.com/heandroro/agencia-viagem/backend/pkg/crypto"
)

// CreateReservationHandler handler para criar reservas
type CreateReservationHandler struct {
	createReservationUC *usecase.CreateReservationUseCase
}

// NewCreateReservationHandler cria um novo handler
func NewCreateReservationHandler(
	repo reservation.Repository,
	availRepo reservation.AvailabilityRepository,
	crypto crypto.Service,
) *CreateReservationHandler {
	return &CreateReservationHandler{
		createReservationUC: usecase.NewCreateReservationUseCase(repo, availRepo, crypto),
	}
}

// Handle processa a requisição de criação de reserva
func (h *CreateReservationHandler) Handle(c *gin.Context) {
	var input reservation.CreateReservationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Dados da requisição inválidos",
		})
		return
	}

	// Extrair informações do contexto (seriam populadas por middleware de auth)
	input.UserID = c.GetString("user_id")
	if input.UserID == "" {
		input.UserID = "anonymous" // TODO: Remover quando auth estiver implementado
	}
	input.IPAddress = c.ClientIP()
	input.UserAgent = c.Request.UserAgent()

	// Executar use case
	output, err := h.createReservationUC.Execute(c.Request.Context(), input)
	if err != nil {
		switch err.Error() {
		case "package_unavailable":
			c.JSON(http.StatusConflict, gin.H{
				"error":   "package_unavailable",
				"message": "Pacote não disponível para as datas selecionadas",
			})
		case "invalid_dates":
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_dates",
				"message": "Datas de viagem inválidas",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Erro interno ao processar reserva",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, output)
}

// CreateReservationRequest representa a requisição de criação de reserva
type CreateReservationRequest struct {
	PackageID     string    `json:"package_id" binding:"required"`
	StartDate     time.Time `json:"start_date" binding:"required"`
	EndDate       time.Time `json:"end_date" binding:"required,gtfield=StartDate"`
	TravelerCount int       `json:"traveler_count" binding:"required,min=1,max=10"`
}
