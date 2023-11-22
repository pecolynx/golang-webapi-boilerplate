package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/pecolynx/golang-webapi-boilerplate/src/controller/gin/helper"
	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
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
	helper.HandleSecuredFunction(c, func(ctx context.Context, logger *slog.Logger, operatorID domain.AppUserID) error {
		standardUserID, err := domain.NewStandardUserID(operatorID.Int())
		if err != nil {
			return err
		}

		entity := TicketAddParameter{}
		if err := c.ShouldBindJSON(&entity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			logger.WarnContext(ctx, "c.ShouldBindJSON", slog.Any("err", err))
			return nil
		}

		param, err := ToTicketAddParameter(&entity)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			logger.WarnContext(ctx, "toTicketAddParameter", slog.Any("err", err))
			return nil
		}

		ticketID, err := h.ticketCreatorUsecase.AddTicket(ctx, standardUserID, param)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, IDResponse{ID: ticketID.Int()})
		return nil
	}, h.errorHandle)
}

func (h *ticketHandler) FindMyTickets(c *gin.Context) {
	helper.HandleSecuredFunction(c, func(ctx context.Context, logger *slog.Logger, operatorID domain.AppUserID) error {
		standardUserID, err := domain.NewStandardUserID(operatorID.Int())
		if err != nil {
			return err
		}

		paramEntity := TicketFindParameter{}
		if err := c.ShouldBindJSON(&paramEntity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			logger.WarnContext(ctx, "c.ShouldBindJSON", slog.Any("err", err))
			return nil
		}

		param, err := ToTicketSearchCondition(&paramEntity)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			logger.WarnContext(ctx, "toTicketFindParameter", slog.Any("err", err))
			return nil
		}

		result, err := h.ticketCreatorUsecase.FindTickets(ctx, standardUserID, param)
		if err != nil {
			return err
		}

		response, err := ToTicketFindResponse(ctx, result)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *ticketHandler) errorHandle(ctx context.Context, logger *slog.Logger, c *gin.Context, err error) bool {
	// if errors.Is(err, service.ErrAudioNotFound) {
	// 	logger.WarnContext(ctx, "", slog.Any("err", err))
	// 	c.JSON(http.StatusNotFound, gin.H{"message": ""})
	// 	return true
	// }
	logger.ErrorContext(ctx, "", slog.Any("err", err))
	return false
}

func NewInitTicketRouterFunc(ticketCreatorUsecase usecase.TicketCreatorUsecase) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		ticket := parentRouterGroup.Group("ticket")
		ticketHandler := NewTicketHandler(ticketCreatorUsecase)
		for _, m := range middleware {
			ticket.Use(m)
		}
		ticket.POST("", ticketHandler.AddTicket)
		return nil
	}
}
