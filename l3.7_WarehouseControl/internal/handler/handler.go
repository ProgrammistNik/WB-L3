package handler

import (
	"strconv"

	"github.com/ProgrammistNik/WB-L3/l3.7_WarehouseControl/internal/handler/dto"
	"github.com/ProgrammistNik/WB-L3/l3.7_WarehouseControl/internal/middlerware"
	"github.com/wb-go/wbf/ginext"
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
	r := ginext.New("debug")

	r.Use(ginext.Logger(), ginext.Recovery())
	r.StaticFile("/", "./web/index.html")

	api := r.Group("/")
	api.POST("/login", h.Login)

	items := api.Group("/items")
	items.GET("", middlerware.Auth("admin", "manager", "viewer"), h.List)
	items.POST("", middlerware.Auth("admin", "manager"), h.Create)
	items.PUT("/:id", middlerware.Auth("admin", "manager"), h.Update)
	items.DELETE("/:id", middlerware.Auth("admin"), h.Delete)
	items.GET("/:id/history", middlerware.Auth("admin", "manager", "viewer"), h.History)

	return r
}

func (h *Handler) Login(c *ginext.Context) {
	var req dto.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, ginext.H{"error": err.Error()})
		return
	}

	role := req.Username
	token, err := middlerware.GenerateToken(req.Username, role)
	if err != nil {
		c.JSON(500, ginext.H{"error": "cannot generate token"})
		return
	}

	c.JSON(200, ginext.H{"token": token})
}

func (h *Handler) List(c *ginext.Context) {
	res, err := h.service.ListItems()
	if err != nil {
		c.JSON(500, ginext.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) Create(c *ginext.Context) {
	var req dto.CreateItemRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, ginext.H{"error": err.Error()})
		return
	}

	username := c.GetString("user")

	err := h.service.CreateItem(username, req.Name, req.Quantity)
	if err != nil {
		c.JSON(500, ginext.H{"error": err.Error()})
		return
	}

	c.Status(200)
}

func (h *Handler) Update(c *ginext.Context) {
	var req dto.UpdateItemRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, ginext.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, ginext.H{"error": "invalid id"})
		return
	}

	username := c.GetString("user")

	if err := h.service.UpdateItem(username, id, req.Name, req.Quantity); err != nil {
		c.JSON(500, ginext.H{"error": err.Error()})
		return
	}

	c.Status(200)
}

func (h *Handler) Delete(c *ginext.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, ginext.H{"error": "invalid id"})
		return
	}

	username := c.GetString("user")

	if err := h.service.DeleteItem(username, id); err != nil {
		c.JSON(500, ginext.H{"error": err.Error()})
		return
	}

	c.Status(200)
}

func (h *Handler) History(c *ginext.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, ginext.H{"error": "invalid id"})
		return
	}

	res, err := h.service.GetHistory(id)
	if err != nil {
		c.JSON(500, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(200, res)
}
