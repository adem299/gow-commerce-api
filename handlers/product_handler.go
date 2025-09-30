package handlers

import (
	"net/http"

	"github.com/adem299/gow-commerce.git/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateProduct(c *gin.Context) {
	var input struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{Name: input.Name, Description: input.Description, Price: input.Price}
	h.DB.Create(&product)

	c.JSON(http.StatusCreated, gin.H{"data": product})
}

func (h *Handler) GetProducts(c *gin.Context) {
	var products []models.Product
	// database.DB.Find(&products)
	h.DB.Find(&products)

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func (h *Handler) GetProductByID(c *gin.Context) {
	var product models.Product

	if err := h.DB.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	var product models.Product

	if err := h.DB.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	h.DB.Save(&product)

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	var product models.Product

	if err := h.DB.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	h.DB.Delete(&product)
	c.JSON(http.StatusOK, gin.H{"data deleted": true})
}
