package api

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
		logrus.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.CreateProduct(input)
	if err != nil {
		logrus.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
func (h *Handler) updateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Request.URL.Query().Get("id"))
	if err != nil {
		logrus.Error(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input repository.UpdateProducts
	if err := c.BindJSON(&input); err != nil {
		logrus.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Update(id, input); err != nil {
		logrus.Error(c, http.StatusInternalServerError, err.Error())
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
		logrus.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success"})
}
func (h *Handler) getProduct(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Request.URL.Query().Get("id"))

	if err != nil {
		logrus.Error(c, http.StatusBadRequest, "invalid id param")
		return
	}

	item, err := h.services.GetById(itemId)
	if err != nil {
		logrus.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}
func (h *Handler) getProducts(c *gin.Context) {
	items, err := h.services.GetAll()
	if err != nil {
		logrus.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)
}
