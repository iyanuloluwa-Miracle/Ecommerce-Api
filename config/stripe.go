package config

import (
	"log"
	"os"

	"github.com/stripe/stripe-go/v76"
)

func InitStripe() {
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		log.Println("Warning: STRIPE_SECRET_KEY not set, payment features will not work")
		return
	}
	stripe.Key = stripeKey
	log.Println("Stripe initialized successfully")
}
