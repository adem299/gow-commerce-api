package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID      uint
	User        User
	Total       float64
	OrderItems  []OrderItem
	TotalAmount uint
	Status      string `gorm:"default:'pending'"`
}

type OrderItem struct {
	gorm.Model
	OrderID     uint
	Order       Order
	ProductID   uint
	Product     Product
	Quantity    uint
	PriceAtTime float64
}
