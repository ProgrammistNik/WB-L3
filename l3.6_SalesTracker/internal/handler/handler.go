package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.6_SalesTracker/internal/handler/dto"
	"github.com/ProgrammistNik/WB-L3/l3.6_SalesTracker/internal/model"
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
	router := ginext.New("release")
	router.Use(ginext.Logger(), ginext.Recovery())

	router.POST("/items", h.CreateItem)
	router.GET("/analytics", h.GetAnalytics)
	router.GET("/items", h.GetItems)
	router.PUT("/items/:id", h.UpdateItem)
	router.DELETE("/items/:id", h.DeleteItem)

	router.Static("/static", "./web")
	router.GET("/", func(c *ginext.Context) {
		c.File("./web/index.html")
	})

	return router
}

func (h *Handler) CreateItem(c *ginext.Context) {
	var req dto.ItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid date format, must be YYYY-MM-DD"})
		return
	}

	createdItem, err := h.service.CreateItem(c, model.Item{
		Type:     req.Type,
		Category: req.Category,
		Amount:   req.Amount,
		Date:     date,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdItem)
}

func (h *Handler) GetAnalytics(c *ginext.Context) {
	var req dto.AnalyticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	filter := model.ItemsFilter{}

	if req.From != nil {
		fromDate, err := time.Parse("2006-01-02", *req.From)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid 'from' date format"})
			return
		}
		filter.From = &fromDate
	}

	if req.To != nil {
		toDate, err := time.Parse("2006-01-02", *req.To)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid 'to' date format"})
			return
		}
		filter.To = &toDate
	}

	filter.Category = req.Category

	analytics, err := h.service.GetAnalytics(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

func (h *Handler) GetItems(c *ginext.Context) {
	var req dto.GetItemsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	filter := model.ItemsFilter{}

	if req.From != nil {
		fromDate, err := time.Parse("2006-01-02", *req.From)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid 'from' date format"})
			return
		}
		filter.From = &fromDate
	}

	if req.To != nil {
		toDate, err := time.Parse("2006-01-02", *req.To)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid 'to' date format"})
			return
		}
		filter.To = &toDate
	}

	filter.Category = req.Category
	filter.Limit = req.Limit
	filter.Offset = req.Offset

	items, err := h.service.GetItems(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) UpdateItem(c *ginext.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid item id"})
		return
	}

	var req dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	item := model.Item{ID: id}

	if req.Type != nil {
		item.Type = *req.Type
	}
	if req.Category != nil {
		item.Category = *req.Category
	}
	if req.Amount != nil {
		item.Amount = *req.Amount
	}
	if req.Date != nil {
		date, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid date format"})
			return
		}
		item.Date = date
	}

	updatedItem, err := h.service.UpdateItem(c, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedItem)
}

func (h *Handler) DeleteItem(c *ginext.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid item id"})
		return
	}

	if err := h.service.DeleteItem(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
