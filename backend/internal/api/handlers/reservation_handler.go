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
	var req CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Dados da requisição inválidos: " + err.Error(),
		})
		return
	}

	// Converter datas de string para time.Time
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_date_format",
			"message": "Formato de data inválido. Use YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_date_format",
			"message": "Formato de data inválido. Use YYYY-MM-DD",
		})
		return
	}

	// Criar input para use case
	input := reservation.CreateReservationInput{
		PackageID:     req.PackageID,
		StartDate:     startDate,
		EndDate:       endDate,
		TravelerCount: req.TravelerCount,
		UserID:        c.GetString("user_id"),
		IPAddress:     c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
	}

	if input.UserID == "" {
		input.UserID = "anonymous" // TODO: Remover quando auth estiver implementado
	}

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
	PackageID     string `json:"package_id" binding:"required"`
	StartDate     string `json:"start_date" binding:"required"`
	EndDate       string `json:"end_date" binding:"required"`
	TravelerCount int    `json:"traveler_count" binding:"required,min=1,max=10"`
}
