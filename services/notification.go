package services

import (
	"log"
	"os"
	"strconv"

	"fmt"

	"github.com/adem299/gow-commerce.git/models"
	"gopkg.in/gomail.v2"
)

type NotificationService interface {
	SendOrderConfirmation(user models.User, order models.Order) error
}

type emailService struct {
	host     string
	port     int
	username string
	password string
}

func NewEmailService() NotificationService {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT: %v", err)
	}

	return &emailService{
		host:     os.Getenv("SMTP_HOST"),
		port:     port,
		username: os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASS"),
	}
}

func (s *emailService) SendOrderConfirmation(user models.User, order models.Order) error {
	m := gomail.NewMessage()

	m.SetHeader("From", s.username)
	m.SetHeader("To", user.Email)
	// m.SetHeader("To", s.username)

	m.SetHeader("Subject", fmt.Sprintf("Order Confirmation #%d", order.ID))
	body := fmt.Sprintf(`
		<h1>Thank you, %s!</h1>
		<p>Your order with ID <strong>#%d</strong> has been successfully placed.</p>
		<p>Total Amount: <strong>USD %d</strong></p>
		`, user.Username, order.ID, order.TotalAmount)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.host, s.port, s.username, s.password)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email for Order ID %d: %v", order.ID, err)
	}

	log.Printf("Email sent for Order ID %d", order.ID)

	return nil
}
