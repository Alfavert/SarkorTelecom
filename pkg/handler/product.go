package handler

import (
	"SarkorTelekom/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (h *Handler) CreateProduct(c *gin.Context) {
	var input repository.Products
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.CreateProduct(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
func (h *Handler) updateProduct(c *gin.Context) {
	idStr := c.Request.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	var input repository.UpdateProducts
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if err := h.services.Update(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success"})
}

func (h *Handler) deleteProduct(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Request.URL.Query().Get("id"))
	if err != nil {
		logrus.Error(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	err = h.services.Delete(itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success"})
}
func (h *Handler) getProduct(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Request.URL.Query().Get("id"))

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	item, err := h.services.GetById(itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}
func (h *Handler) getProducts(c *gin.Context) {
	items, err := h.services.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)
}
