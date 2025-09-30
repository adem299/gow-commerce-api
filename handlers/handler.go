package handlers

import (
	"github.com/adem299/gow-commerce.git/services"
	"gorm.io/gorm"
)

type Handler struct {
	DB           *gorm.DB
	NotifService services.NotificationService
}

func NewHandler(db *gorm.DB, notifService services.NotificationService) *Handler {
	return &Handler{
		DB:           db,
		NotifService: notifService,
	}
}
