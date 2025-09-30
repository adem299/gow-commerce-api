package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID    uint
	User      User
	CartItems []CartItem
}

type CartItem struct {
	gorm.Model
	CartID    uint
	Cart      Cart
	ProductID uint
	Product   Product
	Quantity  int `gorm:"default:1"`
}
