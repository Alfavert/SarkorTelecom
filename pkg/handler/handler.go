package handler

import (
	"SarkorTelekom/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	api := router.Group("/api")
	{
		api.POST("/products", h.CreateProduct)
		api.GET("/product", h.getProduct)
		api.PUT("/product", h.updateProduct)
		api.DELETE("/product", h.deleteProduct)
		api.GET("/products", h.getProducts)
	}
	return router
}
