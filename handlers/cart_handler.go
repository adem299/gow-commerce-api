package handlers

import (
	"net/http"

	"github.com/adem299/gow-commerce.git/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) AddToCart(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(models.User)

	var input struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Quantity < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity must be at least 1"})
		return
	}

	var product models.Product
	if err := h.DB.First(&product, input.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var cart models.Cart
	h.DB.FirstOrCreate(&cart, models.Cart{UserID: currentUser.ID})

	var cartItem models.CartItem
	err := h.DB.Where("cart_id = ? AND product_id = ?", cart.ID, input.ProductID).First(&cartItem).Error

	if err == nil {
		cartItem.Quantity += input.Quantity
		h.DB.Save(&cartItem)
	} else {
		newCartItem := models.CartItem{
			CartID:    cart.ID,
			ProductID: input.ProductID,
			Quantity:  input.Quantity,
		}

		h.DB.Create(&newCartItem)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart"})
}

func (h *Handler) GetCart(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := user.(models.User)

	type CartItemResponse struct {
		ID          uint    `json:"id"`
		ProductName string  `json:"product_name"`
		Price       float32 `json:"price"`
		Quantity    uint    `json:"quantity"`
		Subtotal    uint    `json:"subtotal"`
	}

	var cart models.Cart
	err := h.DB.Preload("CartItems.Product").Where("user_id = ?", currentUser.ID).First(&cart).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"cart": []CartItemResponse{}, "total": 0})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
		return
	}

	var cartItemsResponse []CartItemResponse
	var total float64 = 0
	for _, item := range cart.CartItems {
		subtotal := float64(item.Product.Price) * float64(item.Quantity)
		cartItemsResponse = append(cartItemsResponse, CartItemResponse{
			ID:          item.ID,
			ProductName: item.Product.Name,
			Price:       float32(item.Product.Price),
			Quantity:    uint(item.Quantity),
			Subtotal:    uint(subtotal),
		})
		total += float64(subtotal)
	}

	c.JSON(http.StatusOK, gin.H{"cart_items": cartItemsResponse, "total": total})
}
