package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/services"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Checkout processes the checkout and creates an order with Stripe payment
// @Summary Create checkout session
// @Description Process checkout and create a Stripe payment session for the user's cart
// @Tags Checkout
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=object{checkout_url=string,session_id=string}} "Checkout session created"
// @Failure 400 {object} utils.Response "Bad request - cart is empty"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /checkout [post]
func Checkout(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid user ID")
		return
	}

	// Get the user's cart
	cart, err := services.GetCart(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch cart")
		return
	}

	if len(cart.Items) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "cart is empty")
		return
	}

	// Calculate total and create payment intent
	checkoutSession, err := services.ProcessCheckout(userID, cart)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "checkout session created", gin.H{
		"checkout_url": checkoutSession.URL,
		"session_id":   checkoutSession.ID,
	})
}

// GetOrders returns all orders for the authenticated user
// @Summary Get user orders
// @Description Retrieve all orders for the authenticated user
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=object{orders=[]models.Order}} "Orders retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /orders [get]
func GetOrders(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid user ID")
		return
	}

	orders, err := services.GetUserOrders(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch orders")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "orders retrieved successfully", gin.H{"orders": orders})
}

// WebhookHandler handles Stripe webhook events
// @Summary Handle Stripe webhook
// @Description Process Stripe webhook events for payment status updates
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param Stripe-Signature header string true "Stripe webhook signature"
// @Success 200 {object} utils.Response{data=object{received=bool}} "Webhook processed"
// @Failure 400 {object} utils.Response "Bad request - invalid signature or payload"
// @Router /webhook/stripe [post]
func WebhookHandler(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "failed to read request body")
		return
	}

	signature := c.GetHeader("Stripe-Signature")

	if err := services.HandleStripeWebhook(payload, signature); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "webhook processed", gin.H{"received": true})
}
