package services

import (
	"time"

	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/models"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateOrder creates a new order
func CreateOrder(userID primitive.ObjectID, items []models.CartItem, total float64) error {
	order := models.Order{
		UserID:        userID,
		Items:         items,
		TotalAmount:   total,
		PaymentStatus: "pending",
		CreatedAt:     time.Now(),
	}
	return repositories.CreateOrder(order)
}
