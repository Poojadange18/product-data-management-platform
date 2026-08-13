package handler

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"product-data-management-platform/model"
	"product-data-management-platform/service"
	"strconv"
	"strings"
)

type ProductHandler struct{ Service *service.ProductService }

func NewProductHandler(s *service.ProductService) *ProductHandler { return &ProductHandler{Service: s} }
func productID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"error": "invalid product ID"})
		return 0, false
	}
	return id, true
}
func serverError(c *gin.Context) { c.JSON(500, gin.H{"error": "an internal server error occurred"}) }
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var p model.Product
	if c.ShouldBindJSON(&p) != nil {
		c.JSON(400, gin.H{"error": "invalid request data"})
		return
	}
	if err := h.Service.CreateProduct(&p); err != nil {
		var dbErr *pq.Error
		if errors.As(err, &dbErr) && dbErr.Code == "23505" {
			c.JSON(409, gin.H{"error": "product ID already exists"})
			return
		}
		if strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "negative") || strings.Contains(err.Error(), "positive") {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		serverError(c)
		return
	}
	c.JSON(201, gin.H{"data": p})
}
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var q model.ProductQuery
	if c.ShouldBindQuery(&q) != nil {
		c.JSON(400, gin.H{"error": "invalid query parameters"})
		return
	}
	if q.Stock != "" && q.Stock != "in-stock" && q.Stock != "low-stock" && q.Stock != "out-of-stock" {
		c.JSON(400, gin.H{"error": "invalid stock filter"})
		return
	}
	result, err := h.Service.GetProducts(q)
	if err != nil {
		serverError(c)
		return
	}
	// Preserve the original unfiltered endpoint contract; use the paginated
	// response whenever a query option is supplied.
	if c.Request.URL.RawQuery == "" {
		c.JSON(200, result.Data)
		return
	}
	c.JSON(200, result)
}
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, ok := productID(c)
	if !ok {
		return
	}
	p, err := h.Service.GetProductByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, gin.H{"error": "product not found"})
		return
	}
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			serverError(c)
		}
		return
	}
	c.JSON(200, gin.H{"data": p})
}
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, ok := productID(c)
	if !ok {
		return
	}
	var p model.Product
	if c.ShouldBindJSON(&p) != nil {
		c.JSON(400, gin.H{"error": "invalid request data"})
		return
	}
	if err := h.Service.UpdateProduct(id, &p); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "product not found"})
			return
		}
		if strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "negative") {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		serverError(c)
		return
	}
	c.JSON(200, gin.H{"data": p})
}
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, ok := productID(c)
	if !ok {
		return
	}
	if err := h.Service.DeleteProduct(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "product not found"})
		} else {
			serverError(c)
		}
		return
	}
	c.JSON(200, gin.H{"message": "product deleted successfully"})
}
