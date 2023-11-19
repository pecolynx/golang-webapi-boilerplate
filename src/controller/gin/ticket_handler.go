package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	liblog "github.com/pecolynx/golang-webapi-boilerplate/lib/log"
	"github.com/pecolynx/golang-webapi-boilerplate/src/controller/gin/helper"
	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
	"github.com/pecolynx/golang-webapi-boilerplate/src/log"
	"github.com/pecolynx/golang-webapi-boilerplate/src/usecase"
)

type TicketHandler interface {
	AddTicket(c *gin.Context)
}

type ticketHandler struct {
	ticketCreatorUsecase usecase.TicketCreatorUsecase
}

func NewTicketHandler(ticketCreatorUsecase usecase.TicketCreatorUsecase) TicketHandler {
	return &ticketHandler{
		ticketCreatorUsecase: ticketCreatorUsecase,
	}
}

func (h *ticketHandler) AddTicket(c *gin.Context) {
	ctx := c.Request.Context()
	logger := liblog.GetLoggerFromContext(ctx, log.AppControllerLoggerContextKey)

	helper.HandleSecuredFunction(c, func(operatorID domain.AppUserID) error {
		ticketCreatorID, err := domain.NewTicketCreatorID(operatorID.Int())
		if err != nil {
			return err
		}

		entity := TicketAddParameter{}
		if err := c.ShouldBindJSON(&entity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			logger.WarnContext(ctx, "failed to BindJSON", slog.Any("err", err))
			return nil
		}

		param, err := ToTicketAddParameter(&entity)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			logger.WarnContext(ctx, "failed to BindJSON", slog.Any("err", err))
			return nil
		}

		ticketID, err := h.ticketCreatorUsecase.AddTicket(ctx, ticketCreatorID, param)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, IDResponse{ID: ticketID.Int()})
		return nil
	}, h.errorHandle)
}

func (h *ticketHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := liblog.GetLoggerFromContext(ctx, log.AppControllerLoggerContextKey)
	// if errors.Is(err, service.ErrAudioNotFound) {
	// 	logger.WarnContext(ctx, "", slog.Any("err", err))
	// 	c.JSON(http.StatusNotFound, gin.H{"message": ""})
	// 	return true
	// }
	logger.ErrorContext(ctx, "", slog.Any("err", err))
	return false
}
