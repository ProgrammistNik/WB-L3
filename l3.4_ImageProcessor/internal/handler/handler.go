package handler

import (
	"net/http"

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

	router.POST("/upload", h.UploadImage)
	router.GET("/image/:id", h.GetImage)
	router.DELETE("/image/:id", h.DeleteImage)

	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("web/*.html")

	router.GET("/", func(c *ginext.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	return router
}

func (h *Handler) UploadImage(c *ginext.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "image file not loaded"})
		return
	}
	defer file.Close()

	imageID, err := h.service.UploadImage(c.Request.Context(), file, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginext.H{
		"id":     imageID,
		"status": "pending",
	})
}

func (h *Handler) GetImage(c *ginext.Context) {
	id := c.Param("id")
	image, err := h.service.GetImage(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": err.Error()})
		return
	}

	if image.Status != "completed" {
		c.JSON(http.StatusAccepted, ginext.H{
			"id":     image.ID,
			"status": image.Status,
		})
		return
	}

	c.File(image.OriginalPath)
}

func (h *Handler) DeleteImage(c *ginext.Context) {
	id := c.Param("id")
	if err := h.service.DeleteImage(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ginext.H{"status": "deleted"})
}
