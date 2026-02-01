package handler

import (
	"strconv"

	"github.com/ProgrammistNik/WB-L3/tree/main/l3.3_CommentTree/internal/handler/dto"
	"github.com/ProgrammistNik/WB-L3/tree/main/l3.3_CommentTree/internal/handler/tools"
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

	router.POST("/comments", h.CreateComment)
	router.GET("/comments", h.GetAllComments)
	router.DELETE("/comments/:id", h.Delete)

	router.Static("/static", "./web")
	router.GET("/", func(c *ginext.Context) {
		c.File("./web/index.html")
	})

	return router
}

func (h *Handler) CreateComment(c *ginext.Context) {
	var commentIn dto.CommentRequest
	err := c.BindJSON(&commentIn)
	if err != nil {
		tools.SendError(c, 400, err.Error())
		return
	}

	comment, err := h.service.CreateComment(commentIn)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to create cooment")
		tools.SendError(c, 500, "failed to create cooment")
		return
	}

	tools.SendSuccess(c, 202, ginext.H{
		"comment": comment,
	})

}

func (h *Handler) GetAllComments(c *ginext.Context) {
	parent := c.Query("parent")
	search := c.Query("search")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page := 1
	limit := 20
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	if search != "" {
		comments, err := h.service.SearchComments(search, page, limit)
		if err != nil {
			tools.SendError(c, 404, "Comments not found")
			return
		}
		tools.SendSuccess(c, 200, ginext.H{"comments": comments})
		return
	}

	comments, err := h.service.GetAllComments(parent)
	if err != nil {
		tools.SendError(c, 404, "Comments not found")
		return
	}

	tools.SendSuccess(c, 200, ginext.H{"comments": comments})
}

func (h *Handler) Delete(c *ginext.Context) {
	id := c.Param("id")
	err := h.service.DeleteComment(id)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to delete comment")
		tools.SendError(c, 404, "failed to delete comment")
		return
	}

	tools.SendSuccess(c, 200, ginext.H{
		"status": "cancel",
	})
}
