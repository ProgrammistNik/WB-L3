package handler

import (
	"github.com/ProgrammistNik/WB-L3/l3.1/internal/handler/dto"
	"github.com/ProgrammistNik/WB-L3/l3.1/internal/handler/tools"
	"github.com/ProgrammistNik/WB-L3/l3.1/internal/model"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type Handler struct {
	service Service
}

func New(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Router() *ginext.Engine {
	router := ginext.New()

	router.POST("/notify", h.Notify)
	router.GET("/notify/:id", h.NotifyGetID)
	router.DELETE("/notify/:id", h.Delete)

	router.Static("/ui", "./web")
	router.GET("/", func(c *ginext.Context) {
		c.Redirect(302, "/ui/index.html")
	})

	return router
}

func (h *Handler) Notify(c *ginext.Context) {
	var notifyRequest dto.NotificationRequest
	err := c.BindJSON(&notifyRequest)
	if err != nil {
		tools.SendError(c, 400, err.Error())
		return
	}

	notif := model.CastToNotification(notifyRequest)

	err = h.service.CreateNotification(c, notif)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to create notification")
		tools.SendError(c, 503, "failed to schedule notification")
		return
	}

	tools.SendSuccess(c, 202, ginext.H{
		"id":     notif.ID,
		"status": notif.Status,
	})
}

func (h *Handler) NotifyGetID(c *ginext.Context) {
	id := c.Param("id")

	notify, err := h.service.GetStatusByID(c, id)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to find status by notification")
		tools.SendError(c, 503, "failed to get id")
		return
	}

	tools.SendSuccess(c, 202, ginext.H{
		"id":     notify.ID,
		"status": notify.Status,
	})
}

func (h *Handler) Delete(c *ginext.Context) {
	id := c.Param("id")
	err := h.service.DeleteNotify(c, id)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to delete notification")
		tools.SendError(c, 404, "failed to delete notification")
		return
	}
	tools.SendSuccess(c, 200, ginext.H{
		"status": "canceled",
	})
}
