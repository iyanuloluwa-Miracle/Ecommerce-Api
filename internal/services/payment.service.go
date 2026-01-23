package services

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/models"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/repositories"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/webhook"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckoutSession struct {
	ID  string
	URL string
}

// ProcessCheckout creates a Stripe checkout session for the cart
func ProcessCheckout(userID primitive.ObjectID, cart models.Cart) (*CheckoutSession, error) {
	// Calculate cart total
	var lineItems []*stripe.CheckoutSessionLineItemParams
	var totalAmount float64

	for _, item := range cart.Items {
		product, err := repositories.GetProductByID(item.ProductID)
		if err != nil {
			return nil, errors.New("product not found: " + item.ProductID.Hex())
		}

		// Check stock
		if product.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		itemTotal := product.Price * float64(item.Quantity)
		totalAmount += itemTotal

		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("usd"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:        stripe.String(product.Name),
					Description: stripe.String(product.Description),
				},
				UnitAmount: stripe.Int64(int64(product.Price * 100)), // Stripe uses cents
			},
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	// Create Stripe checkout session
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String(os.Getenv("FRONTEND_URL") + "/checkout/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:          stripe.String(os.Getenv("FRONTEND_URL") + "/checkout/cancel"),
		Metadata: map[string]string{
			"user_id": userID.Hex(),
		},
	}

	s, err := session.New(params)
	if err != nil {
		return nil, err
	}

	// Create pending order
	order := models.Order{
		UserID:        userID,
		Items:         cart.Items,
		TotalAmount:   totalAmount,
		PaymentStatus: "pending",
		StripeSession: s.ID,
		CreatedAt:     time.Now(),
	}

	if err := repositories.CreateOrder(order); err != nil {
		return nil, err
	}

	return &CheckoutSession{
		ID:  s.ID,
		URL: s.URL,
	}, nil
}

// HandleStripeWebhook processes Stripe webhook events
func HandleStripeWebhook(payload []byte, signature string) error {
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSecret == "" {
		return errors.New("STRIPE_WEBHOOK_SECRET not set")
	}

	event, err := webhook.ConstructEvent(payload, signature, webhookSecret)
	if err != nil {
		return err
	}

	switch event.Type {
	case "checkout.session.completed":
		var sessionData stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &sessionData); err != nil {
			return err
		}

		// Update order status
		if err := repositories.UpdateOrderPaymentStatus(sessionData.ID, "paid"); err != nil {
			return err
		}

		// Get the order to update stock
		order, err := repositories.GetOrderBySessionID(sessionData.ID)
		if err != nil {
			return err
		}

		// Update product stock
		for _, item := range order.Items {
			if err := repositories.DecrementProductStock(item.ProductID, item.Quantity); err != nil {
				return err
			}
		}

		// Clear the user's cart
		if err := repositories.ClearCart(order.UserID); err != nil {
			return err
		}

	case "checkout.session.expired":
		var sessionData stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &sessionData); err != nil {
			return err
		}

		// Update order status to expired
		if err := repositories.UpdateOrderPaymentStatus(sessionData.ID, "expired"); err != nil {
			return err
		}
	}

	return nil
}

// GetUserOrders returns all orders for a user
func GetUserOrders(userID primitive.ObjectID) ([]models.Order, error) {
	return repositories.GetOrdersByUserID(userID)
}
