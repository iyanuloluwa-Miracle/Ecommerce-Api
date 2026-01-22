package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	Items         []CartItem         `bson:"items" json:"items"`
	TotalAmount   float64            `bson:"total_amount" json:"total_amount"`
	PaymentStatus string             `bson:"payment_status" json:"payment_status"`
	StripeSession string             `bson:"stripe_session" json:"stripe_session,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}
