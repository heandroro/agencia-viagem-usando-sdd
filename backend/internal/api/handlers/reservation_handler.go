package handlers

import (
	"log"
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

// UpdateTravelersHandler handler para atualizar viajantes
type UpdateTravelersHandler struct {
	updateTravelersUC *usecase.UpdateTravelersUseCase
}

// NewUpdateTravelersHandler cria um novo handler
func NewUpdateTravelersHandler(
	repo reservation.Repository,
	cryptoService crypto.Service,
) *UpdateTravelersHandler {
	return &UpdateTravelersHandler{
		updateTravelersUC: usecase.NewUpdateTravelersUseCase(repo, cryptoService),
	}
}

// Handle processa a requisição de atualização de viajantes
func (h *UpdateTravelersHandler) Handle(c *gin.Context) {
	reservationID := c.Param("id")
	if reservationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "ID da reserva é obrigatório",
		})
		return
	}

	var req UpdateTravelersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Dados da requisição inválidos: " + err.Error(),
		})
		return
	}

	// Converter viajantes
	travelers := make([]usecase.TravelerInput, len(req.Travelers))
	for i, t := range req.Travelers {
		birthDate, err := time.Parse("2006-01-02", t.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_date_format",
				"message": "Formato de data inválido. Use YYYY-MM-DD",
			})
			return
		}

		travelers[i] = usecase.TravelerInput{
			Type:           reservation.TravelerType(t.Type),
			FullName:       t.FullName,
			DocumentType:   reservation.DocumentType(t.DocumentType),
			DocumentNumber: t.DocumentNumber,
			BirthDate:      birthDate,
		}
	}

	// Criar input para use case
	input := usecase.UpdateTravelersInput{
		ReservationID: reservationID,
		UserID:        c.GetString("user_id"),
		Travelers:     travelers,
	}

	if input.UserID == "" {
		input.UserID = "anonymous" // TODO: Remover quando auth estiver implementado
	}

	// Executar use case
	output, err := h.updateTravelersUC.Execute(c.Request.Context(), input)
	if err != nil {
		switch err.Error() {
		case "reservation_not_found":
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "reservation_not_found",
				"message": "Reserva não encontrada",
			})
		case "unauthorized":
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "unauthorized",
				"message": "Não autorizado a modificar esta reserva",
			})
		case "reservation_not_pending":
			c.JSON(http.StatusConflict, gin.H{
				"error":   "reservation_not_pending",
				"message": "Reserva não está em status pendente",
			})
		default:
			log.Printf("[ERROR] UpdateTravelers: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Erro interno ao atualizar viajantes",
			})
		}
		return
	}

	// Converter para resposta
	response := UpdateTravelersResponse{
		ReservationID: output.ReservationID,
		Status:        string(output.Status),
		Travelers:     make([]TravelerResponse, len(output.Travelers)),
	}
	for i, t := range output.Travelers {
		response.Travelers[i] = TravelerResponse{
			Type:           t.Type,
			FullName:       t.FullName,
			DocumentType:   t.DocumentType,
			DocumentNumber: t.DocumentNumber,
			BirthDate:      t.BirthDate,
		}
	}

	c.JSON(http.StatusOK, response)
}

// CreateReservationRequest representa a requisição de criação de reserva
type CreateReservationRequest struct {
	PackageID     string `json:"package_id" binding:"required"`
	StartDate     string `json:"start_date" binding:"required"`
	EndDate       string `json:"end_date" binding:"required"`
	TravelerCount int    `json:"traveler_count" binding:"required,min=1,max=10"`
}

// TravelerRequest representa um viajante na requisição
type TravelerRequest struct {
	Type           string `json:"type" binding:"required,oneof=primary companion"`
	FullName       string `json:"full_name" binding:"required"`
	DocumentType   string `json:"document_type" binding:"required,oneof=cpf passport"`
	DocumentNumber string `json:"document_number" binding:"required"`
	BirthDate      string `json:"birth_date" binding:"required"`
}

// UpdateTravelersRequest representa a requisição de atualização de viajantes
type UpdateTravelersRequest struct {
	Travelers []TravelerRequest `json:"travelers" binding:"required,dive,required"`
}

// UpdateTravelersResponse representa a resposta de atualização de viajantes
type UpdateTravelersResponse struct {
	ReservationID string             `json:"reservation_id"`
	Status        string             `json:"status"`
	Travelers     []TravelerResponse `json:"travelers"`
}

// TravelerResponse representa um viajante na resposta (com máscara)
type TravelerResponse struct {
	Type           string `json:"type"`
	FullName       string `json:"full_name"`
	DocumentType   string `json:"document_type"`
	DocumentNumber string `json:"document_number"` // Mascarado
	BirthDate      string `json:"birth_date"`
}

// GetReservationSummaryHandler handler para buscar resumo da reserva
type GetReservationSummaryHandler struct {
	getSummaryUC *usecase.GetReservationSummaryUseCase
}

// NewGetReservationSummaryHandler cria um novo handler
func NewGetReservationSummaryHandler(repo reservation.Repository) *GetReservationSummaryHandler {
	return &GetReservationSummaryHandler{
		getSummaryUC: usecase.NewGetReservationSummaryUseCase(repo),
	}
}

// Handle processa a requisição de busca do resumo
func (h *GetReservationSummaryHandler) Handle(c *gin.Context) {
	reservationID := c.Param("id")
	if reservationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "ID da reserva é obrigatório",
		})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		userID = "anonymous" // TODO: Remover quando auth estiver implementado
	}

	input := usecase.GetReservationSummaryInput{
		ReservationID: reservationID,
		UserID:        userID,
	}

	output, err := h.getSummaryUC.Execute(c.Request.Context(), input)
	if err != nil {
		switch err.Error() {
		case "reservation_not_found":
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "reservation_not_found",
				"message": "Reserva não encontrada",
			})
		case "unauthorized":
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "unauthorized",
				"message": "Não autorizado a acessar esta reserva",
			})
		default:
			log.Printf("[ERROR] GetReservationSummary: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": "Erro interno ao buscar resumo",
			})
		}
		return
	}

	c.JSON(http.StatusOK, output)
}
