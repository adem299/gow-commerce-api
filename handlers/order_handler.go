package handlers

import (
	"net/http"
	"time"

	"github.com/adem299/gow-commerce.git/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) CreateOrder(c *gin.Context) {
	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	var cart models.Cart
	if err := h.DB.Preload("CartItems.Product").Where("user_id = ?", currentUser.ID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	if len(cart.CartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
	}

	var newOrder models.Order
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var totalAmount uint = 0
		for _, item := range cart.CartItems {
			totalAmount += uint(item.Product.Price) * uint(item.Quantity)
		}

		newOrder = models.Order{
			UserID:      currentUser.ID,
			TotalAmount: totalAmount,
			Status:      "pending",
		}
		if err := tx.Create(&newOrder).Error; err != nil {
			return err
		}

		for _, item := range cart.CartItems {
			orderItem := models.OrderItem{
				OrderID:     newOrder.ID,
				ProductID:   item.ProductID,
				Quantity:    uint(item.Quantity),
				PriceAtTime: item.Product.Price,
			}

			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Pesanan berhasil dibuat",
		"order_id":   newOrder.ID,
		"total":      newOrder.TotalAmount,
		"status":     newOrder.Status,
		"created_at": newOrder.CreatedAt,
	})
}

func (h *Handler) GetOrders(c *gin.Context) {
	user, _ := c.Get("currentUser")
	currentUser := user.(models.User)

	type OrderItemResponse struct {
		ProductName string  `json:"product_name"`
		Quantity    uint    `json:"quantity"`
		PriceAtTime float64 `json:"price_at_time"`
	}

	type OrderResponse struct {
		ID          uint                `json:"id"`
		TotalAmount uint                `json:"total_amount"`
		Status      string              `json:"status"`
		Items       []OrderItemResponse `json:"items"`
		CreatedAt   time.Time           `json:"created_at"`
	}

	var orders []models.Order
	h.DB.Preload("OrderItems.Product").Where("user_id = ?", currentUser.ID).Order("created_at desc").Find(&orders)

	var ordersResponse []OrderResponse
	for _, order := range orders {
		var itemsResponse []OrderItemResponse
		for _, item := range order.OrderItems {
			itemsResponse = append(itemsResponse, OrderItemResponse{
				ProductName: item.Product.Name,
				Quantity:    item.Quantity,
				PriceAtTime: item.PriceAtTime,
			})
		}

		ordersResponse = append(ordersResponse, OrderResponse{
			ID:          order.ID,
			CreatedAt:   order.CreatedAt,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			Items:       itemsResponse,
		})
	}

	c.JSON(http.StatusOK, gin.H{"orders": ordersResponse})
}

// admin (only)
func (h *Handler) GetAllOrders(c *gin.Context) {
	var orders []models.Order
	h.DB.Preload("OrderItems.Product").Preload("User").Order("created_at desc").Find(&orders)

	type OrderItemResponse struct {
		ProductName string  `json:"product_name"`
		Quantity    uint    `json:"quantity"`
		PriceAtTime float64 `json:"price_at_time"`
	}

	type OrderResponse struct {
		ID          uint                `json:"id"`
		TotalAmount uint                `json:"total_amount"`
		Status      string              `json:"status"`
		Items       []OrderItemResponse `json:"items"`
		CreatedAt   time.Time           `json:"created_at"`
		User        models.User         `json:"user"`
	}

	var ordersResponse []OrderResponse
	for _, order := range orders {
		var itemsResponse []OrderItemResponse
		for _, item := range order.OrderItems {
			itemsResponse = append(itemsResponse, OrderItemResponse{
				ProductName: item.Product.Name,
				Quantity:    item.Quantity,
				PriceAtTime: item.PriceAtTime,
			})
		}
		ordersResponse = append(ordersResponse, OrderResponse{
			ID:          order.ID,
			CreatedAt:   order.CreatedAt,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			Items:       itemsResponse,
			User:        order.User,
		})
	}

	c.JSON(http.StatusOK, gin.H{"orders": ordersResponse})
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID := c.Param("id")
	var order models.Order
	if err := h.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.Status = input.Status
	h.DB.Save(&order)

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}
